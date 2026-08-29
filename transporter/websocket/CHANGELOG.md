# Changelog

## [2.0.0](https://github.com/Liphium/neoroute/compare/transporter/websocket/v1.0.0...transporter/websocket/v2.0.0) (2026-08-29)


### ⚠ BREAKING CHANGES

* Upgrade everything to Go 1.27 generic methods
* **websocket:** remove transporter from EnterNetworkFunc as it's no longer needed to register adapter
* **websocket:** add support for adapt und disconnect on session and remove old adapt function
* move transporter to extra modules, to avoid unused imports

### Features

* add generator support for all default transporters ([15a2992](https://github.com/Liphium/neoroute/commit/15a29929f5397f064657eaf9b2b1c4d0d601c02d))
* add panic protection to websocket and http neoroute transporter ([42aad5a](https://github.com/Liphium/neoroute/commit/42aad5a32297dc8998b75aeea1ef6608aa4dcb94))
* add schema generation for transporters ([604cdb9](https://github.com/Liphium/neoroute/commit/604cdb9660594678998b88eb30ea28688484f97a))
* support http transporter and full generation ([cb3e08e](https://github.com/Liphium/neoroute/commit/cb3e08ee12d71465f70e766dd31852c441136567))
* Use writers instead of reading the complete byte stream ([e578fcf](https://github.com/Liphium/neoroute/commit/e578fcfad3b99c748aeff3c3e39f70dbd1ec5fc8))
* **websocket:** Add option to disable origin checking ([f09dbd9](https://github.com/Liphium/neoroute/commit/f09dbd9fe22ab4f05cbe643288a93b7c8a381a01))
* **websocket:** Add option to disable origin checking ([63cff11](https://github.com/Liphium/neoroute/commit/63cff1139254fbc6d0db0dcf6f4b811a138e5438))
* **websocket:** add support for adapt und disconnect on session and remove old adapt function ([262d9dc](https://github.com/Liphium/neoroute/commit/262d9dc47d280d7e17ba89667b0a7f8b3f1ff408))
* **websocket:** Let the AcceptOptions bet set ([bd4b6fc](https://github.com/Liphium/neoroute/commit/bd4b6fc3df2799b004ab7645ee389bb56f1eeded))
* **websocket:** remove transporter from EnterNetworkFunc as it's no longer needed to register adapter ([0c975e4](https://github.com/Liphium/neoroute/commit/0c975e4b4f3b9efe94f11bed3e89fd9dfefee0df))
* **websocket:** remove unneeded config parameters ([0a8a9db](https://github.com/Liphium/neoroute/commit/0a8a9dbf7f703c9eda599cf9a5f89d94c4a20e9a))
* **websocket:** use replacer to always use the newest version of neoroute with new releases ([c2e7554](https://github.com/Liphium/neoroute/commit/c2e7554607fc17d638d4b3bea05d43099c7b8020))


### Bug Fixes

* **transporter:** make transporter importable ([4db8ba3](https://github.com/Liphium/neoroute/commit/4db8ba3bb97c690e57a5c98e7aaa212fc9b0dd42))
* trigger release manually to test it ([3f0c5ce](https://github.com/Liphium/neoroute/commit/3f0c5ce22a42fc95a757d204c8cb66d287a394ad))


### Code Refactoring

* move transporter to extra modules, to avoid unused imports ([652ccd7](https://github.com/Liphium/neoroute/commit/652ccd7c425245255240e5a2918352bfc8f75d2f))
* Upgrade everything to Go 1.27 generic methods ([67f0a91](https://github.com/Liphium/neoroute/commit/67f0a91e478ef21047696c4c210d9ea4c2ae4d9a))

## [1.0.0](https://github.com/Liphium/neoroute/compare/transporter/websocket/v0.8.0...transporter/websocket/v1.0.0) (2026-08-29)


### ⚠ BREAKING CHANGES

* Upgrade everything to Go 1.27 generic methods

### Features

* Use writers instead of reading the complete byte stream ([e578fcf](https://github.com/Liphium/neoroute/commit/e578fcfad3b99c748aeff3c3e39f70dbd1ec5fc8))


### Code Refactoring

* Upgrade everything to Go 1.27 generic methods ([67f0a91](https://github.com/Liphium/neoroute/commit/67f0a91e478ef21047696c4c210d9ea4c2ae4d9a))

## [0.8.0](https://github.com/Liphium/neoroute/compare/transporter/websocket/v0.7.0...transporter/websocket/v0.8.0) (2026-08-06)


### Features

* **websocket:** Add option to disable origin checking ([f09dbd9](https://github.com/Liphium/neoroute/commit/f09dbd9fe22ab4f05cbe643288a93b7c8a381a01))
* **websocket:** Add option to disable origin checking ([63cff11](https://github.com/Liphium/neoroute/commit/63cff1139254fbc6d0db0dcf6f4b811a138e5438))
* **websocket:** Let the AcceptOptions bet set ([bd4b6fc](https://github.com/Liphium/neoroute/commit/bd4b6fc3df2799b004ab7645ee389bb56f1eeded))

## [0.7.0](https://github.com/Liphium/neoroute/compare/transporter/websocket/v0.6.0...transporter/websocket/v0.7.0) (2026-07-21)


### ⚠ BREAKING CHANGES

* **websocket:** remove transporter from EnterNetworkFunc as it's no longer needed to register adapter
* **websocket:** add support for adapt und disconnect on session and remove old adapt function

### Features

* **websocket:** add support for adapt und disconnect on session and remove old adapt function ([262d9dc](https://github.com/Liphium/neoroute/commit/262d9dc47d280d7e17ba89667b0a7f8b3f1ff408))
* **websocket:** remove transporter from EnterNetworkFunc as it's no longer needed to register adapter ([0c975e4](https://github.com/Liphium/neoroute/commit/0c975e4b4f3b9efe94f11bed3e89fd9dfefee0df))

## [0.6.0](https://github.com/Liphium/neoroute/compare/transporter/websocket/v0.5.0...transporter/websocket/v0.6.0) (2026-07-11)


### Features

* add panic protection to websocket and http neoroute transporter ([42aad5a](https://github.com/Liphium/neoroute/commit/42aad5a32297dc8998b75aeea1ef6608aa4dcb94))
* add schema generation for transporters ([604cdb9](https://github.com/Liphium/neoroute/commit/604cdb9660594678998b88eb30ea28688484f97a))
* support http transporter and full generation ([cb3e08e](https://github.com/Liphium/neoroute/commit/cb3e08ee12d71465f70e766dd31852c441136567))

## [0.5.0](https://github.com/Liphium/neoroute/compare/transporter/websocket/v0.4.1...transporter/websocket/v0.5.0) (2026-06-28)


### Features

* **websocket:** use replacer to always use the newest version of neoroute with new releases ([c2e7554](https://github.com/Liphium/neoroute/commit/c2e7554607fc17d638d4b3bea05d43099c7b8020))

## [0.4.1](https://github.com/Liphium/neoroute/compare/transporter/websocket-v0.4.0...transporter/websocket/v0.4.1) (2026-06-27)


### Bug Fixes

* trigger release manually to test it ([3f0c5ce](https://github.com/Liphium/neoroute/commit/3f0c5ce22a42fc95a757d204c8cb66d287a394ad))

## [0.4.0](https://github.com/Liphium/neoroute/compare/transporter/websocket-v0.3.0...transporter/websocket-v0.4.0) (2026-06-27)


### ⚠ BREAKING CHANGES

* move transporter to extra modules, to avoid unused imports

### Features

* **websocket:** remove unneeded config parameters ([0a8a9db](https://github.com/Liphium/neoroute/commit/0a8a9dbf7f703c9eda599cf9a5f89d94c4a20e9a))


### Bug Fixes

* **transporter:** make transporter importable ([4db8ba3](https://github.com/Liphium/neoroute/commit/4db8ba3bb97c690e57a5c98e7aaa212fc9b0dd42))


### Code Refactoring

* move transporter to extra modules, to avoid unused imports ([652ccd7](https://github.com/Liphium/neoroute/commit/652ccd7c425245255240e5a2918352bfc8f75d2f))

## [0.3.0](https://github.com/Liphium/neoroute/compare/v0.2.0...v0.3.0) (2026-06-27)


### ⚠ BREAKING CHANGES

* move transporter to extra modules, to avoid unused imports

### Features

* **websocket:** remove unneeded config parameters ([0a8a9db](https://github.com/Liphium/neoroute/commit/0a8a9dbf7f703c9eda599cf9a5f89d94c4a20e9a))


### Bug Fixes

* **transporter:** make transporter importable ([4db8ba3](https://github.com/Liphium/neoroute/commit/4db8ba3bb97c690e57a5c98e7aaa212fc9b0dd42))


### Code Refactoring

* move transporter to extra modules, to avoid unused imports ([652ccd7](https://github.com/Liphium/neoroute/commit/652ccd7c425245255240e5a2918352bfc8f75d2f))

## [0.2.0](https://github.com/Liphium/neoroute/compare/v0.1.2...v0.2.0) (2026-06-27)


### ⚠ BREAKING CHANGES

* move transporter to extra modules, to avoid unused imports

### Features

* **websocket:** remove unneeded config parameters ([0a8a9db](https://github.com/Liphium/neoroute/commit/0a8a9dbf7f703c9eda599cf9a5f89d94c4a20e9a))


### Bug Fixes

* **transporter:** make transporter importable ([4db8ba3](https://github.com/Liphium/neoroute/commit/4db8ba3bb97c690e57a5c98e7aaa212fc9b0dd42))


### Code Refactoring

* move transporter to extra modules, to avoid unused imports ([652ccd7](https://github.com/Liphium/neoroute/commit/652ccd7c425245255240e5a2918352bfc8f75d2f))
