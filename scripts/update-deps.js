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
	[".", "client"],
	[
		"client/transporter/http",
		"client/transporter/websocket",
		"transporter/http",
		"transporter/websocket",
	],
	["pkg/neodebug"],
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
	const branch =
		process.env.UPDATER_REF ||
		run("git rev-parse --abbrev-ref HEAD").trim();

	// Version of every released module, from the release-please manifest.
	const manifest = JSON.parse(
		fs.readFileSync(path.join(process.cwd(), MANIFEST_FILE), "utf8"),
	);

	for (let step = 0; step < DEPENDENCY_TREE.length; step++) {
		log(`=== Step ${step}: [${DEPENDENCY_TREE[step].join(", ")}] ===`);

		// Modules from earlier steps are already at their latest versions.
		const previousModules = [];
		for (let i = 0; i < step; i++) {
			for (const mod of DEPENDENCY_TREE[i]) {
				const version = manifest[mod];
				if (version === undefined) {
					log(`WARNING: no version in manifest for "${mod}", skipping`);
					continue;
				}
				previousModules.push({ module: mod, version });
			}
		}

		if (previousModules.length === 0) continue;

		for (const mod of DEPENDENCY_TREE[step]) {
			log(`--- Processing "${mod}" ---`);

			let goMod;
			try {
				goMod = fs.readFileSync(path.join(mod, "go.mod"), "utf8");
			} catch (err) {
				console.error(`[update-deps] ERROR: cannot read ${mod}/go.mod: ${err.message}`);
				process.exit(1);
			}

			for (const { module: dep, version } of previousModules) {
				const importPath = `${MODULE_PREFIX}/${dep}`;
				// Root module is imported without a path prefix.
				const rootImport = dep === ".";
				const found = rootImport
					? new RegExp(
							`^\\s*${MODULE_PREFIX.replace(/\./g, "\\.")}\\s`,
							"m",
						).test(goMod)
					: goMod.includes(importPath);

				if (!found) {
					log(`"${importPath}" not imported in ${mod}, skipping`);
					continue;
				}

				const target = rootImport
					? `${MODULE_PREFIX}@v${version}`
					: `${importPath}@v${version}`;
				log(`go get ${target} (in ${mod})`);
				try {
					const output = run(`go get ${target}`, { cwd: mod });
					if (output.trim()) console.log(output.trim());
				} catch (err) {
					console.error(`[update-deps] ERROR: go get ${target} in ${mod} failed:`);
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
		.flatMap((m) => [`${m}/go.mod`, `${m}/go.sum`])
		.filter((f) => fs.existsSync(f));
	run(`git add ${files.join(" ")}`);
	run('git commit -m "chore: upgrade dependencies"');
	run(`git pull --rebase origin ${branch}`);
	run(`git push origin HEAD:${branch}`);
	log("Pushed dependency updates.");
}

main();
