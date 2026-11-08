# Changelog

## [2.0.0](https://github.com/Liphium/neoroute/compare/pkg/neodebug/v1.1.0...pkg/neodebug/v2.0.0) (2026-11-08)


### ⚠ BREAKING CHANGES

* Neoroute v2 ([#49](https://github.com/Liphium/neoroute/issues/49))

### Features

* **neodebug:** Add playground for trying out neodebug ([fd7b428](https://github.com/Liphium/neoroute/commit/fd7b42885b8d9dee864a4bdf67087e900f12646e))
* **neodebug:** Fuzzy search for finding routes more easily ([fd7b428](https://github.com/Liphium/neoroute/commit/fd7b42885b8d9dee864a4bdf67087e900f12646e))
* **neodebug:** Input of time is now available ([fd7b428](https://github.com/Liphium/neoroute/commit/fd7b42885b8d9dee864a4bdf67087e900f12646e))
* **neodebug:** Presets: pre-defined requests ([fd7b428](https://github.com/Liphium/neoroute/commit/fd7b42885b8d9dee864a4bdf67087e900f12646e))
* **neogen:** Support to generate structs as tuples in MessagePack ([fd7b428](https://github.com/Liphium/neoroute/commit/fd7b42885b8d9dee864a4bdf67087e900f12646e))
* **neogen:** time.Time is now a specific type (not generated as struct) ([fd7b428](https://github.com/Liphium/neoroute/commit/fd7b42885b8d9dee864a4bdf67087e900f12646e))
* Neoroute v2 ([#49](https://github.com/Liphium/neoroute/issues/49)) ([fd7b428](https://github.com/Liphium/neoroute/commit/fd7b42885b8d9dee864a4bdf67087e900f12646e))


### Bug Fixes

* **neogen:** Neogen no longer generates unexported struct fields ([fd7b428](https://github.com/Liphium/neoroute/commit/fd7b42885b8d9dee864a4bdf67087e900f12646e))
* Various fixes and performance improvements ([fd7b428](https://github.com/Liphium/neoroute/commit/fd7b42885b8d9dee864a4bdf67087e900f12646e))

## [1.1.0](https://github.com/Liphium/neoroute/compare/pkg/neodebug/v1.0.0...pkg/neodebug/v1.1.0) (2026-09-30)


### Features

* **neodebug:** Add hotkey (ctrl+d) for toggling mouse capture ([7c8c3ff](https://github.com/Liphium/neoroute/commit/7c8c3ff387b18741f8172c8f5468f914f0a00c4f))
* **neodebug:** Add request history for resending requests ([31cb561](https://github.com/Liphium/neoroute/commit/31cb561f1a65b34deaf06e1eeb0c26c77d775ef7))
* **neodebug:** New features for better usability ([#40](https://github.com/Liphium/neoroute/issues/40)) ([31cb561](https://github.com/Liphium/neoroute/commit/31cb561f1a65b34deaf06e1eeb0c26c77d775ef7))
* **neodebug:** Replace individual bytes with file uploads ([31cb561](https://github.com/Liphium/neoroute/commit/31cb561f1a65b34deaf06e1eeb0c26c77d775ef7))

## 1.0.0 (2026-08-29)


### ⚠ BREAKING CHANGES

* Upgrade everything to Go 1.27 generic methods

### Features

* Add debugger for neoroute transporters ([#28](https://github.com/Liphium/neoroute/issues/28)) ([47b32a2](https://github.com/Liphium/neoroute/commit/47b32a2540cfbf0747d568bfa4d840faed390428))
* **neodebug:** Add debugger for neoroute servers ([47b32a2](https://github.com/Liphium/neoroute/commit/47b32a2540cfbf0747d568bfa4d840faed390428))


### Bug Fixes

* Make sure middlewares work properly across groups ([47b32a2](https://github.com/Liphium/neoroute/commit/47b32a2540cfbf0747d568bfa4d840faed390428))
* trigger release manually to test it ([3f0c5ce](https://github.com/Liphium/neoroute/commit/3f0c5ce22a42fc95a757d204c8cb66d287a394ad))


### Code Refactoring

* Upgrade everything to Go 1.27 generic methods ([67f0a91](https://github.com/Liphium/neoroute/commit/67f0a91e478ef21047696c4c210d9ea4c2ae4d9a))
