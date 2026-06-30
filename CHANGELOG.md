# Changelog

## [0.4.0](https://github.com/Liphium/neoroute/compare/v0.3.1...v0.4.0) (2026-06-30)


### ⚠ BREAKING CHANGES

* **neogate:** change messageEvent to panic instead of returning an error, as the error should never occur

### Features

* **neoroute:** remove Use function as use is implemented on each router ([bfafc73](https://github.com/Liphium/neoroute/commit/bfafc735c9dfeb523f02ab81e118d95f1607b740))


### Bug Fixes

* **neoroute:** actually initialize neos slice to avoid nil pointer ([bff3e79](https://github.com/Liphium/neoroute/commit/bff3e79048fb0cc45acd410db117c7e5f53700dd))
* **neoroute:** return self as well when getNeos is called on NeoRouter ([a6c8863](https://github.com/Liphium/neoroute/commit/a6c886352af34531b724a070005015767031f3e4))
* **neoroute:** the getNeos now return all neos even of sub neos ([cc9e839](https://github.com/Liphium/neoroute/commit/cc9e8392eea70e3d71ff67a9597b9fd902936820))


### Code Refactoring

* **neogate:** change messageEvent to panic instead of returning an error, as the error should never occur ([a553392](https://github.com/Liphium/neoroute/commit/a55339242954191544b267bc0e987588f48838ee))

## [0.3.1](https://github.com/Liphium/neoroute/compare/v0.3.0...v0.3.1) (2026-06-27)


### Bug Fixes

* force update root module files ([598938b](https://github.com/Liphium/neoroute/commit/598938b8a59898b0dcbf4da6adb5d1dadf7c55ca))

## [0.3.0](https://github.com/Liphium/neoroute/compare/v0.2.0...v0.3.0) (2026-06-27)


### ⚠ BREAKING CHANGES

* **web_transport:** rename type to shorter versions
* move transporter to extra modules, to avoid unused imports

### Features

* **websocket:** remove unneeded config parameters ([0a8a9db](https://github.com/Liphium/neoroute/commit/0a8a9dbf7f703c9eda599cf9a5f89d94c4a20e9a))


### Bug Fixes

* **client/websocket:** make module importable ([d2e3ea7](https://github.com/Liphium/neoroute/commit/d2e3ea7ec38fbdde1888e0a567cc4793192ae407))
* **transporter:** make transporter importable ([4db8ba3](https://github.com/Liphium/neoroute/commit/4db8ba3bb97c690e57a5c98e7aaa212fc9b0dd42))


### Code Refactoring

* move transporter to extra modules, to avoid unused imports ([652ccd7](https://github.com/Liphium/neoroute/commit/652ccd7c425245255240e5a2918352bfc8f75d2f))
* **web_transport:** rename type to shorter versions ([b8bb56c](https://github.com/Liphium/neoroute/commit/b8bb56ce217ae122878701761a14126d10a4b6c3))
