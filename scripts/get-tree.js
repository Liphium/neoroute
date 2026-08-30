#!/usr/bin/env node
"use strict";

const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");

// ---------------------------------------------------------------------------
// Dependency tree: each inner array is one "step". Steps must be executed in
// order — everything in a later step depends on everything in earlier steps.
// ---------------------------------------------------------------------------
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
	console.log(`[get-tree] ${msg}`);
}

function run(cmd, opts = {}) {
	log(`$ ${cmd}`);
	return execSync(cmd, {
		stdio: ["ignore", "pipe", "pipe"],
		encoding: "utf8",
		...opts,
	});
}

function main() {
	const step = parseInt(process.argv[2], 10);
	const branch = process.argv[3]; // ${{ github.head_ref || github.ref_name }}
	if (
		Number.isNaN(step) ||
		step < 0 ||
		step >= DEPENDENCY_TREE.length ||
		!branch
	) {
		console.error(
			`Usage: node scripts/get-tree.js <step> <branch> (step: 0..${
				DEPENDENCY_TREE.length - 1
			}, branch: \${{ github.head_ref || github.ref_name }})`,
		);
		process.exit(1);
	}

	log(`Executing step ${step}: [${DEPENDENCY_TREE[step].join(", ")}]`);

	// Load versions from the release-please manifest.
	const manifestPath = path.join(process.cwd(), MANIFEST_FILE);
	const manifest = JSON.parse(fs.readFileSync(manifestPath, "utf8"));

	// Collect all modules from previous steps with their versions.
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

	if (previousModules.length === 0) {
		log("No modules from previous steps, nothing to wire up.");
		return;
	}

	// For each module in the current step, check go.mod for imports of
	// previous modules and run `go get` for the ones found.
	for (const mod of DEPENDENCY_TREE[step]) {
		log(`--- Processing "${mod}" ---`);

		let goMod;
		try {
			goMod = fs.readFileSync(path.join(mod, "go.mod"), "utf8");
		} catch (err) {
			console.error(
				`[get-tree] ERROR: cannot read ${mod}/go.mod: ${err.message}`,
			);
			process.exit(1);
		}

		for (const { module: dep, version } of previousModules) {
			const importPath = `${MODULE_PREFIX}/${dep}`;
			// Skip when the dep is the root module (imported without prefix).
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
			log(`Updating: cd ${mod} && go get ${target}`);
			try {
				const output = execSync(`go get ${target}`, {
					cwd: mod,
					stdio: ["ignore", "pipe", "pipe"],
					encoding: "utf8",
				});
				if (output.trim()) console.log(output.trim());
				log(`Done: ${target} in ${mod}`);
			} catch (err) {
				console.error(`[get-tree] ERROR: go get ${target} in ${mod} failed:`);
				console.error(err.stderr || err.message);
				process.exit(1);
			}
		}
	}

	log(`Step ${step} finished.`);

	finishStep(step, branch);
}

// ---------------------------------------------------------------------------
// Commit dependency upgrades, then release every module in the current step.
// ---------------------------------------------------------------------------
function finishStep(step, branch) {
	log("=== Finishing step: git + releases ===");

	// Commit updates if there are any.
	const status = run("git status --porcelain").trim();
	if (status) {
		log("Changes detected, committing...");
		run(
			`git add ${DEPENDENCY_TREE[step]
				.flatMap((m) => [`${m}/go.mod`, `${m}/go.sum`])
				.join(" ")}`,
		);
		run('git commit -m "chore: upgrade dependencies"');
		run(`git pull --rebase origin ${branch}`);
		run(`git push origin HEAD:${branch}`);
	} else {
		log("No changes, skipping commit.");
	}

	// Release all modules in the current step.
	for (const mod of DEPENDENCY_TREE[step]) {
		log(`Releasing "${mod}"...`);
		try {
			const output = run(`bash scripts/release.sh ${mod}`);
			if (output.trim()) console.log(output.trim());
			log(`Release done: ${mod}`);
		} catch (err) {
			console.error(`[get-tree] ERROR: release.sh ${mod} failed:`);
			console.error(err.stderr || err.message);
			process.exit(1);
		}
	}

	log("Step fully finished.");
}

main();
