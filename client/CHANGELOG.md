# Changelog

## [0.6.0](https://github.com/Liphium/neoroute/compare/client/v0.5.0...client/v0.6.0) (2026-07-11)


### ⚠ BREAKING CHANGES

* **client:** rename the RouteNoop function to RoutePing

### Features

* **client:** add run error handler function with default logger to avoid nil pointer ([7ee50b3](https://github.com/Liphium/neoroute/commit/7ee50b310c1ffa580948bd12d4547302d5d3c87b))


### Code Refactoring

* **client:** rename the RouteNoop function to RoutePing ([a75c651](https://github.com/Liphium/neoroute/commit/a75c651a350bb0a35dc47bcbeceafb5a859fa28b))

## [0.5.0](https://github.com/Liphium/neoroute/compare/client-v0.4.0...client/v0.5.0) (2026-06-30)


### ⚠ BREAKING CHANGES

* **client:** return only one error instead of an user error string and  an error and use a custom error for user errors instead

### Features

* **client:** return only one error instead of an user error string and  an error and use a custom error for user errors instead ([5e76825](https://github.com/Liphium/neoroute/commit/5e76825a63e5a21df34bd93289f5ec71cbf81404))

## [0.4.0](https://github.com/Liphium/neoroute/compare/client-v0.3.0...client-v0.4.0) (2026-06-27)


### ⚠ BREAKING CHANGES

* move transporter to extra modules, to avoid unused imports

### Features

* add client ([1865a72](https://github.com/Liphium/neoroute/commit/1865a72bfe17567ae289ba6fa2ce6a50cb8fe8ee))
* add http transporter for client and make Receiver easier to use with send ([50dfebb](https://github.com/Liphium/neoroute/commit/50dfebb047baee3bf55b935c3d45359526209ceb))
* remove error return from client receiver ([4862728](https://github.com/Liphium/neoroute/commit/4862728bc9df031adf6589e3aad8ba16f4c0bcce))


### Bug Fixes

* **http:** return error if handshake fails or body cant be read ([5337467](https://github.com/Liphium/neoroute/commit/5337467f9474ed53384cef4ea038b292acf26e39))


### Code Refactoring

* move transporter to extra modules, to avoid unused imports ([652ccd7](https://github.com/Liphium/neoroute/commit/652ccd7c425245255240e5a2918352bfc8f75d2f))

## [0.3.0](https://github.com/Liphium/neoroute/compare/v0.2.0...v0.3.0) (2026-06-27)


### ⚠ BREAKING CHANGES

* move transporter to extra modules, to avoid unused imports

### Code Refactoring

* move transporter to extra modules, to avoid unused imports ([652ccd7](https://github.com/Liphium/neoroute/commit/652ccd7c425245255240e5a2918352bfc8f75d2f))

## [0.2.0](https://github.com/Liphium/neoroute/compare/v0.1.1...v0.2.0) (2026-06-27)


### ⚠ BREAKING CHANGES

* move transporter to extra modules, to avoid unused imports

### Bug Fixes

* **http:** return error if handshake fails or body cant be read ([5337467](https://github.com/Liphium/neoroute/commit/5337467f9474ed53384cef4ea038b292acf26e39))


### Code Refactoring

* move transporter to extra modules, to avoid unused imports ([652ccd7](https://github.com/Liphium/neoroute/commit/652ccd7c425245255240e5a2918352bfc8f75d2f))
