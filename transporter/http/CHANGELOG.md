# Changelog

## [1.0.0](https://github.com/Liphium/neoroute/compare/transporter/http/v0.7.0...transporter/http/v1.0.0) (2026-08-29)


### ⚠ BREAKING CHANGES

* Upgrade everything to Go 1.27 generic methods

### Features

* Use writers instead of reading the complete byte stream ([e578fcf](https://github.com/Liphium/neoroute/commit/e578fcfad3b99c748aeff3c3e39f70dbd1ec5fc8))


### Code Refactoring

* Upgrade everything to Go 1.27 generic methods ([67f0a91](https://github.com/Liphium/neoroute/commit/67f0a91e478ef21047696c4c210d9ea4c2ae4d9a))

## [0.7.0](https://github.com/Liphium/neoroute/compare/transporter/http/v0.6.0...transporter/http/v0.7.0) (2026-07-21)


### Features

* **http:** provide empty struct to NewSession function to satisfy definition ([4ed6356](https://github.com/Liphium/neoroute/commit/4ed63569862dd4404d5233b8a80e8bd903e478b5))

## [0.6.0](https://github.com/Liphium/neoroute/compare/transporter/http/v0.5.0...transporter/http/v0.6.0) (2026-07-11)


### Features

* add panic protection to websocket and http neoroute transporter ([42aad5a](https://github.com/Liphium/neoroute/commit/42aad5a32297dc8998b75aeea1ef6608aa4dcb94))
* add schema generation for transporters ([604cdb9](https://github.com/Liphium/neoroute/commit/604cdb9660594678998b88eb30ea28688484f97a))
* support http transporter and full generation ([cb3e08e](https://github.com/Liphium/neoroute/commit/cb3e08ee12d71465f70e766dd31852c441136567))

## [0.5.0](https://github.com/Liphium/neoroute/compare/transporter/http/v0.4.1...transporter/http/v0.5.0) (2026-06-28)


### Features

* **http:** use replacer to always use the newest version of neoroute with new releases ([d00b204](https://github.com/Liphium/neoroute/commit/d00b204da765b1c55fdd7380e96131b97244563e))

## [0.4.1](https://github.com/Liphium/neoroute/compare/transporter/http-v0.4.0...transporter/http/v0.4.1) (2026-06-27)


### Bug Fixes

* trigger release manually to test it ([3f0c5ce](https://github.com/Liphium/neoroute/commit/3f0c5ce22a42fc95a757d204c8cb66d287a394ad))

## [0.4.0](https://github.com/Liphium/neoroute/compare/transporter/http-v0.3.0...transporter/http-v0.4.0) (2026-06-27)


### ⚠ BREAKING CHANGES

* move transporter to extra modules, to avoid unused imports

### Bug Fixes

* **transporter:** make transporter importable ([4db8ba3](https://github.com/Liphium/neoroute/commit/4db8ba3bb97c690e57a5c98e7aaa212fc9b0dd42))


### Code Refactoring

* move transporter to extra modules, to avoid unused imports ([652ccd7](https://github.com/Liphium/neoroute/commit/652ccd7c425245255240e5a2918352bfc8f75d2f))

## [0.3.0](https://github.com/Liphium/neoroute/compare/v0.2.0...v0.3.0) (2026-06-27)


### ⚠ BREAKING CHANGES

* move transporter to extra modules, to avoid unused imports

### Bug Fixes

* **transporter:** make transporter importable ([4db8ba3](https://github.com/Liphium/neoroute/commit/4db8ba3bb97c690e57a5c98e7aaa212fc9b0dd42))


### Code Refactoring

* move transporter to extra modules, to avoid unused imports ([652ccd7](https://github.com/Liphium/neoroute/commit/652ccd7c425245255240e5a2918352bfc8f75d2f))

## [0.2.0](https://github.com/Liphium/neoroute/compare/v0.1.1...v0.2.0) (2026-06-27)


### ⚠ BREAKING CHANGES

* move transporter to extra modules, to avoid unused imports

### Bug Fixes

* **transporter:** make transporter importable ([4db8ba3](https://github.com/Liphium/neoroute/commit/4db8ba3bb97c690e57a5c98e7aaa212fc9b0dd42))


### Code Refactoring

* move transporter to extra modules, to avoid unused imports ([652ccd7](https://github.com/Liphium/neoroute/commit/652ccd7c425245255240e5a2918352bfc8f75d2f))
