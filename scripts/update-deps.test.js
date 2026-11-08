"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const test = require("node:test");
const vm = require("node:vm");

const source = fs.readFileSync(`${__dirname}/update-deps.js`, "utf8");

function simulate({ race = false, dirty = false, pushFailure = false } = {}) {
  const commands = [];
  let base = "initial";
  let pushes = 0;
  let changed = false;
  const invoke = (cmd) => {
    commands.push(cmd);
    if (cmd === "git status --porcelain") return dirty ? " M user-file" : "";
    if (cmd === "git rev-parse FETCH_HEAD") return base;
    if (cmd.startsWith("git status --porcelain --")) return changed ? " M go.mod" : "";
    if (cmd.startsWith("go get")) changed = true;
    if (cmd.startsWith("git push")) {
      pushes++;
      if (pushFailure || (race && pushes === 1)) {
        if (race) base = "concurrent";
        throw new Error("push rejected");
      }
    }
    return "";
  };
  const fakeFs = {
    existsSync: () => true,
    readFileSync(file) {
      if (file.endsWith(".release-please-manifest.json")) {
        commands.push(`read manifest ${base}`);
        return JSON.stringify({ ".": base === "initial" ? "2.0.0-rc.1" : "2.0.0" });
      }
      return "require github.com/Liphium/neoroute/v2 v2.0.0\n";
    },
  };
  let error;
  try {
    vm.runInNewContext(source, {
      require(name) {
        if (name === "fs") return fakeFs;
        if (name === "child_process") return {
          execSync: invoke,
          execFileSync: (file, args) => invoke([file, ...args].join(" ")),
        };
        return require(name);
      },
      process: { env: { UPDATER_REF: "main" }, cwd: () => "/repo" },
      console: { log() {}, error() {} },
    });
  } catch (err) { error = err; }
  return { commands, pushes, error };
}

test("syncs before reading manifest and preserves requested downgrades", () => {
  const { commands, error } = simulate();
  assert.ifError(error);
  assert(commands.indexOf("git merge --ff-only initial") < commands.indexOf("read manifest initial"));
  assert(commands.includes("go get github.com/Liphium/neoroute/v2@v2.0.0-rc.1"));
});

test("regenerates using the new manifest after a concurrent push", () => {
  const { commands, pushes, error } = simulate({ race: true });
  assert.ifError(error);
  assert.equal(pushes, 2);
  assert(commands.includes("git reset --keep concurrent"));
  assert(commands.includes("read manifest concurrent"));
  assert(commands.includes("go get github.com/Liphium/neoroute/v2@v2.0.0"));
  assert(!commands.some((cmd) => cmd.includes("rebase")));
});

test("does not retry unrelated push failures", () => {
  const { commands, pushes, error } = simulate({ pushFailure: true });
  assert.match(error.message, /push rejected/);
  assert.equal(pushes, 1);
  assert(!commands.some((cmd) => cmd.startsWith("git reset")));
});

test("rejects a dirty working tree before fetching or changing files", () => {
  const { commands, error } = simulate({ dirty: true });
  assert.match(error.message, /clean working tree/);
  assert(!commands.some((cmd) => cmd.startsWith("git fetch") || cmd.startsWith("go get")));
});
