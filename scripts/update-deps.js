#!/usr/bin/env node
"use strict";

// Updates every module's go.mod to the latest released versions of the other
// neoroute modules it imports, walking the dependency tree bottom-up. Commits
// "chore: upgrade dependencies" to the current branch (main) when anything
// changed. Does NOT tag, release, or edit changelogs — release-please owns
// those now.

const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");

// Each inner array is one "step". Steps run in order — everything in a later
// step depends on everything in earlier steps.
const DEPENDENCY_TREE = [
  [
    { path: ".", version: "v2" },
    { path: "client/sdk", version: "v2" },
  ],
  [
    { path: "client/transporter/http", version: "v2" },
    { path: "client/transporter/websocket", version: "v2" },
    { path: "transporter/http", version: "v2" },
    { path: "transporter/websocket", version: "v2" },
  ],
  [{ path: "pkg/neodebug", version: "v2" }],
];

const MODULE_PREFIX = "github.com/Liphium/neoroute";
const MANIFEST_FILE = ".release-please-manifest.json";

function log(msg) {
  console.log(`[update-deps] ${msg}`);
}

function run(cmd, opts = {}) {
  return execSync(cmd, { encoding: "utf8", ...opts });
}

function main() {
  const branch = process.env.UPDATER_REF || run("git rev-parse --abbrev-ref HEAD").trim();

  // Version of every released module, from the release-please manifest.
  const manifest = JSON.parse(fs.readFileSync(path.join(process.cwd(), MANIFEST_FILE), "utf8"));

  for (let step = 0; step < DEPENDENCY_TREE.length; step++) {
    log(`=== Step ${step}: [${DEPENDENCY_TREE[step].map((dep) => dep.path).join(", ")}] ===`);

    // Modules from earlier steps are already at their latest versions.
    const previousModules = [];
    for (let i = 0; i < step; i++) {
      for (const mod of DEPENDENCY_TREE[i]) {
        const version = manifest[mod.path];
        if (version === undefined) {
          log(`WARNING: no version in manifest for "${mod.path}", skipping`);
          continue;
        }
        // During a major-version rollout, some modules still have older releases.
        const releasedMajor = `v${version.split(".")[0]}`;
        if (releasedMajor !== mod.version) {
          log(`"${mod.path}" latest release is v${version}, waiting for a ${mod.version} release, skipping`);
          continue;
        }
        previousModules.push({ module: mod, version });
      }
    }

    if (previousModules.length === 0) continue;

    for (const mod of DEPENDENCY_TREE[step]) {
      log(`--- Processing "${mod.path}" ---`);

      let goMod;
      try {
        goMod = fs.readFileSync(path.join(mod.path, "go.mod"), "utf8");
      } catch (err) {
        console.error(`[update-deps] ERROR: cannot read ${mod.path}/go.mod: ${err.message}`);
        process.exit(1);
      }

      for (const { module: dep, version } of previousModules) {
        const importPath = [MODULE_PREFIX, dep.path === "." ? null : dep.path, dep.version]
          .filter(Boolean) // removes falsy values from an array, including null, undefined, "", 0, and false.
          .join("/");
        const found = new RegExp(`^\\s*(?:require\\s+)?${importPath.replace(/\./g, "\\.")}\\s`, "m").test(goMod);

        if (!found) {
          log(`"${importPath}" not imported in ${mod.path}, skipping`);
          continue;
        }

        const target = `${importPath}@v${version}`;
        log(`go get ${target} (in ${mod.path})`);
        try {
          const output = run(`go get ${target}`, { cwd: mod.path });
          if (output.trim()) console.log(output.trim());
        } catch (err) {
          console.error(`[update-deps] ERROR: go get ${target} in ${mod.path} failed:`);
          console.error(err.stderr || err.message);
          process.exit(1);
        }
      }
    }
  }

  // Commit updates across all steps, if there are any.
  const status = run("git status --porcelain").trim();
  if (!status) {
    log("No changes, skipping commit.");
    return;
  }

  log("Changes detected, committing...");
  const files = DEPENDENCY_TREE.flat()
    .flatMap((m) => [`${m.path}/go.mod`, `${m.path}/go.sum`])
    .filter((f) => fs.existsSync(f));
  run(`git add ${files.join(" ")}`);
  run('git commit -m "chore: upgrade dependencies"');
  run(`git pull --rebase origin ${branch}`);
  run(`git push origin HEAD:${branch}`);
  log("Pushed dependency updates.");
}

main();
