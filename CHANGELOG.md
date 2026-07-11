# Changelog

## [0.5.0](https://github.com/Liphium/neoroute/compare/v0.4.0...v0.5.0) (2026-07-11)


### ⚠ BREAKING CHANGES

* **neoroute:** remove uppercase letters from allowed characters
* **neoroute:** change route separator from dot to slash
* **neoroute:** rename the RouteNoop function to RoutePing

### Features

* add panic protection to websocket and http neoroute transporter ([42aad5a](https://github.com/Liphium/neoroute/commit/42aad5a32297dc8998b75aeea1ef6608aa4dcb94))
* add schema generation for transporters ([604cdb9](https://github.com/Liphium/neoroute/commit/604cdb9660594678998b88eb30ea28688484f97a))
* change go version to 1.26 everywhere ([6e07054](https://github.com/Liphium/neoroute/commit/6e07054adc68aa4fd7141fc3bdc9184ce6d0e838))
* change structure slightly to make sure we can test in the future ([751c0e6](https://github.com/Liphium/neoroute/commit/751c0e6acd38f28da83852aaaf59b2cf4cec0027))
* finish websocket generator ([4f9c92f](https://github.com/Liphium/neoroute/commit/4f9c92fe322ee293d8f12aed71a555f5d21b7302))
* make generation simpler using text templates ([e260973](https://github.com/Liphium/neoroute/commit/e26097355cc63e75bb72c4b7dba72e8fece1c1df))
* **neogen:** generate websocket transporter properly ([10c8cc9](https://github.com/Liphium/neoroute/commit/10c8cc95d6ff6d16772af448c804410e9c1c77af))
* **neogen:** support maps ([dad4af3](https://github.com/Liphium/neoroute/commit/dad4af3b1405024efb7ea88e09a66beb9a3776ae))
* **neoroute:** add helper function to create a RouterGroup directly from routers ([164aab1](https://github.com/Liphium/neoroute/commit/164aab1f89d0a38d76a8467734261ecf12f3dfb5))
* **neoroute:** add NoData type for users that want no SessionData ([6bb33ee](https://github.com/Liphium/neoroute/commit/6bb33eeddf2210929ce0d2c14c232f3297a1caf5))
* **neoroute:** remove uppercase letters from allowed characters ([2ee9738](https://github.com/Liphium/neoroute/commit/2ee97384076d001eceb8a74812e96c12951672f3))
* nullable support for pointer types ([afd13d1](https://github.com/Liphium/neoroute/commit/afd13d1eba8148276e6b7a461070d79505049c23))
* prepare or type ([647b185](https://github.com/Liphium/neoroute/commit/647b185b8476e7db7b86f177df7a03572ae4784c))
* simplify generator by using text/template ([63e410b](https://github.com/Liphium/neoroute/commit/63e410bc69f6d8754c84a569ae746eb256097ac9))
* start interface support ([0b671f0](https://github.com/Liphium/neoroute/commit/0b671f015517cbcb90aa2ca26a92718f313308e7))
* support http transporter and full generation ([cb3e08e](https://github.com/Liphium/neoroute/commit/cb3e08ee12d71465f70e766dd31852c441136567))
* update to SendPing ([7ebdd90](https://github.com/Liphium/neoroute/commit/7ebdd900a2439a1fd7a11785ae35261b73585f50))


### Bug Fixes

* calculate coverage properly ([1f068ab](https://github.com/Liphium/neoroute/commit/1f068ab96d670af4bc92ba6bf776cf39a401fc5f))
* go pack to previous + better coverage ([3640c2b](https://github.com/Liphium/neoroute/commit/3640c2b345c73473b7157f1771bac65ace36ddab))
* install msgp in pipeline ([b80be7e](https://github.com/Liphium/neoroute/commit/b80be7e4c6c85919fb244edec559ad625bc2bcb1))
* **neogen:** nullable generation ([82b7ab8](https://github.com/Liphium/neoroute/commit/82b7ab8a63a33b6ec70a72b838cc1d96be5843cc))
* properly generate the coverage report for gobadge ([86116b1](https://github.com/Liphium/neoroute/commit/86116b148156bb5e4d73c74db71fc08a77bdbdf5))
* pull rebase ([7470fa9](https://github.com/Liphium/neoroute/commit/7470fa948231cecae11c9233f33309f4a6000310))
* push to correct thingy ([922b494](https://github.com/Liphium/neoroute/commit/922b4943c9883199a151c31971fd999bcc4ce886))
* use newer go version for pipeline ([93d135f](https://github.com/Liphium/neoroute/commit/93d135ff6d905d62bf23646405cc3bc82179ff1f))
* use stuff from docs ([c4e6ee9](https://github.com/Liphium/neoroute/commit/c4e6ee9ec9742b06ca51ae11a1633efbd0c62deb))


### Code Refactoring

* **neoroute:** change route separator from dot to slash ([063b8fc](https://github.com/Liphium/neoroute/commit/063b8fc3caf0ecb6b1ed321efd52f6815729928e))
* **neoroute:** rename the RouteNoop function to RoutePing ([6a6eb5e](https://github.com/Liphium/neoroute/commit/6a6eb5ed993ed0d78062758e87a14da72251c1fe))

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
