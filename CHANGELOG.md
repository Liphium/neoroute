# Changelog

## [2.0.0](https://github.com/Liphium/neoroute/compare/v1.0.1...v2.0.0) (2026-09-06)


### ⚠ BREAKING CHANGES

* Upgrade everything to Go 1.27 generic methods
* **neogen:** Generate Uint8Array instead of number[]
* **neogen:** Make the WebSocketConfig one parameter
* **neogen:** Correctly extend transporters for TS
* **neogen:** Let WebSocketConfig be configurable
* **neogen:** rename args to command for more control
* **neoroute:** add error return to disconnect function on session
* Remove error return from event creation function on
* **neoroute:** allow multiple middlewares for the same route
* **neoroute:** remove uppercase letters from allowed characters
* **neoroute:** change route separator from dot to slash
* **neoroute:** rename the RouteNoop function to RoutePing
* **neogate:** change messageEvent to panic instead of returning an error, as the error should never occur

### Features

* add ability to disconnect user using an adapter ([28e95da](https://github.com/Liphium/neoroute/commit/28e95da9adab1166a0b8d4077cecf7fda3a056d8))
* add adapters ([345d731](https://github.com/Liphium/neoroute/commit/345d73190af5b2ff6c2decb8c9508db6028152a9))
* add client ([1865a72](https://github.com/Liphium/neoroute/commit/1865a72bfe17567ae289ba6fa2ce6a50cb8fe8ee))
* add debug print to help detect faulty uuid generators ([f117ed1](https://github.com/Liphium/neoroute/commit/f117ed1fd7f25456d44e97338a5d2a56e3f03b59))
* Add debugger for neoroute transporters ([#28](https://github.com/Liphium/neoroute/issues/28)) ([47b32a2](https://github.com/Liphium/neoroute/commit/47b32a2540cfbf0747d568bfa4d840faed390428))
* add event registry to let users create events they want to send ([e1171c7](https://github.com/Liphium/neoroute/commit/e1171c7bceab639299d3e79dd0f54389e9a7449c))
* Add functions for adapter management on registry ([249d2d0](https://github.com/Liphium/neoroute/commit/249d2d073c98bcf4734710fc01000e24788993ee))
* add generation to websocket server example ([60f5994](https://github.com/Liphium/neoroute/commit/60f5994b37e86e67505fa6861e9934ee3987f08e))
* add generator support for all default transporters ([15a2992](https://github.com/Liphium/neoroute/commit/15a29929f5397f064657eaf9b2b1c4d0d601c02d))
* add middleware support ([fc67ba4](https://github.com/Liphium/neoroute/commit/fc67ba47b83527eda99e2b4a2a379b89df817d3c))
* add new function that allows users to update session data concurrency safe ([54e9e2f](https://github.com/Liphium/neoroute/commit/54e9e2f48854d8f23b149d8d2c51b6337f93fe6d))
* add new RouteOk, RouteOkNoResponse and rename other Route functions to be more descriptive and move RouterGroup functionality to AddRouters function in Router interface. ([bc6caf3](https://github.com/Liphium/neoroute/commit/bc6caf3e5e5b833ba800a43edc894d2768c09db7))
* add panic protection to websocket and http neoroute transporter ([42aad5a](https://github.com/Liphium/neoroute/commit/42aad5a32297dc8998b75aeea1ef6608aa4dcb94))
* add runAfter funcs that get executed after the handler returns ([8f70934](https://github.com/Liphium/neoroute/commit/8f70934b661e0cfa00e75b301e507efc0688beb9))
* add runtime check for if a event is actually registered with transporter ([539f436](https://github.com/Liphium/neoroute/commit/539f436e0d3dbbcbe7bb834ebd69c8f68082b386))
* add schema generation for transporters ([604cdb9](https://github.com/Liphium/neoroute/commit/604cdb9660594678998b88eb30ea28688484f97a))
* add session and handshake ([b42fbfe](https://github.com/Liphium/neoroute/commit/b42fbfe5a9645ea5a505068ad3cbf7a33509ae1d))
* add testing util ([7015e8b](https://github.com/Liphium/neoroute/commit/7015e8b3edbbebfd095a7a811d2ee34d0d5ef590))
* add unmarshal  event testing helper ([f418622](https://github.com/Liphium/neoroute/commit/f418622743239892530b0af3647b52b2e656cbd8))
* add websocket transporter ([df53b69](https://github.com/Liphium/neoroute/commit/df53b6964ff771ed3a69e57b0e178f598a087153))
* basic generation of models, routes and events ([aeee208](https://github.com/Liphium/neoroute/commit/aeee208260d1d2cffd88e558a960e005d29a4ce6))
* begin to add basic router and mount ([8f75420](https://github.com/Liphium/neoroute/commit/8f754206ccc293489301b8990f4ff6c0111ebc55))
* change go version to 1.26 everywhere ([6e07054](https://github.com/Liphium/neoroute/commit/6e07054adc68aa4fd7141fc3bdc9184ce6d0e838))
* change structure slightly to make sure we can test in the future ([751c0e6](https://github.com/Liphium/neoroute/commit/751c0e6acd38f28da83852aaaf59b2cf4cec0027))
* finish websocket generator ([4f9c92f](https://github.com/Liphium/neoroute/commit/4f9c92fe322ee293d8f12aed71a555f5d21b7302))
* implement router and simple http transporter ([8f7150e](https://github.com/Liphium/neoroute/commit/8f7150ea62c1ad5d00116d55824e5c06870dd2d3))
* make generation simpler using text templates ([e260973](https://github.com/Liphium/neoroute/commit/e26097355cc63e75bb72c4b7dba72e8fece1c1df))
* make request private as user doesn't need access to it ([86d465a](https://github.com/Liphium/neoroute/commit/86d465a4b20fc529070d112ec862a8e19ab49656))
* More safety checks for handlers ([9d47b0a](https://github.com/Liphium/neoroute/commit/9d47b0adfd06db5f9314736c79b7c55b08c57c3a))
* move err logging and conversion from router to context ([40507df](https://github.com/Liphium/neoroute/commit/40507dfdd4e9e617d5fed844a7bd8106644ef9c7))
* **neodebug:** Add debugger for neoroute servers ([47b32a2](https://github.com/Liphium/neoroute/commit/47b32a2540cfbf0747d568bfa4d840faed390428))
* **neogen:** add typescript support ([c9bd471](https://github.com/Liphium/neoroute/commit/c9bd47108578eeda1f94c4f4d37f4c278500439b))
* **neogen:** add typescript support ([0a4634a](https://github.com/Liphium/neoroute/commit/0a4634a0db4d35d1b18f7663bab6d3c17bd4af27))
* **neogen:** generate websocket transporter properly ([10c8cc9](https://github.com/Liphium/neoroute/commit/10c8cc95d6ff6d16772af448c804410e9c1c77af))
* **neogen:** Make the WebSocketConfig one parameter ([4fa21f5](https://github.com/Liphium/neoroute/commit/4fa21f550d9b09eb416c8a5f7e57fa6e2c1901d9))
* **neogen:** rename args to command for more control ([e2965fb](https://github.com/Liphium/neoroute/commit/e2965fb78f7ab1602452da4ff7aee6f992791c3f))
* **neogen:** support maps ([dad4af3](https://github.com/Liphium/neoroute/commit/dad4af3b1405024efb7ea88e09a66beb9a3776ae))
* **neoroute:** add error return to disconnect function on session ([091d50a](https://github.com/Liphium/neoroute/commit/091d50a3d4b9f2aa08f6ff5f4ad89b7e4f18bd9b))
* **neoroute:** add helper function to create a RouterGroup directly from routers ([164aab1](https://github.com/Liphium/neoroute/commit/164aab1f89d0a38d76a8467734261ecf12f3dfb5))
* **neoroute:** add new errors ([45830ed](https://github.com/Liphium/neoroute/commit/45830ed059da8326091f9763efb9d145694ff668))
* **neoroute:** add NoData type for users that want no SessionData ([6bb33ee](https://github.com/Liphium/neoroute/commit/6bb33eeddf2210929ce0d2c14c232f3297a1caf5))
* **neoroute:** allow multiple middlewares for the same route ([29ab932](https://github.com/Liphium/neoroute/commit/29ab9321325402c6901e8c3f2053377c1ffcaf87))
* **neoroute:** remove unneeded generics ([d26e477](https://github.com/Liphium/neoroute/commit/d26e4774e3b1623d66b676e02361d71b4a314efe))
* **neoroute:** remove unneeded interface in generics ([400650b](https://github.com/Liphium/neoroute/commit/400650b8d1a5575d3e1f74a4ce4e481d276b6c1f))
* **neoroute:** remove uppercase letters from allowed characters ([2ee9738](https://github.com/Liphium/neoroute/commit/2ee97384076d001eceb8a74812e96c12951672f3))
* **neoroute:** remove Use function as use is implemented on each router ([bfafc73](https://github.com/Liphium/neoroute/commit/bfafc735c9dfeb523f02ab81e118d95f1607b740))
* **neoschema:** Support generation for slices ([266aa3c](https://github.com/Liphium/neoroute/commit/266aa3c08d5a3fb26881971b9e82a306432d469c))
* nullable support for pointer types ([afd13d1](https://github.com/Liphium/neoroute/commit/afd13d1eba8148276e6b7a461070d79505049c23))
* prepare or type ([647b185](https://github.com/Liphium/neoroute/commit/647b185b8476e7db7b86f177df7a03572ae4784c))
* remove Data() func from cxt, because users shouldn't have access to it ([d5bf8e0](https://github.com/Liphium/neoroute/commit/d5bf8e06635fa077c0e061fcd0de2f46121b50ac))
* remove unused generic from adapter registry ([f9e423d](https://github.com/Liphium/neoroute/commit/f9e423d57e6cfa2c1ede7120def39895a2ab89c8))
* return new RouteRouter when using Route to allow easier middleware mounting ([31ef86e](https://github.com/Liphium/neoroute/commit/31ef86e103b672cf9d0782f765aab60029023f99))
* rework ctx to work with middleware ([7f50795](https://github.com/Liphium/neoroute/commit/7f507953dc1495ee814639e25315d70048ad87e4))
* routing structure and add WebTransport support and middlewares ([03e0702](https://github.com/Liphium/neoroute/commit/03e0702e9c634d888e13c9a501b2a5e8fa831572))
* schema for generation ([6651330](https://github.com/Liphium/neoroute/commit/6651330c3602fe8eaf865b5dd925ae99d6797a6d))
* schema generation start ([02b1a2e](https://github.com/Liphium/neoroute/commit/02b1a2e50dbe7ad5d6debe287c0d0adf5097c762))
* session ids can be managed by neogate ([16940de](https://github.com/Liphium/neoroute/commit/16940defafce0769cffc47442968119256a0ea10))
* simplify generator by using text/template ([63e410b](https://github.com/Liphium/neoroute/commit/63e410bc69f6d8754c84a569ae746eb256097ac9))
* start interface support ([0b671f0](https://github.com/Liphium/neoroute/commit/0b671f015517cbcb90aa2ca26a92718f313308e7))
* support http transporter and full generation ([cb3e08e](https://github.com/Liphium/neoroute/commit/cb3e08ee12d71465f70e766dd31852c441136567))
* support recursion + start on command ([cfe04e1](https://github.com/Liphium/neoroute/commit/cfe04e1597017f4d3056850ef59bc1e6e6c3389a))
* update context to make it easier to use ([b32b8c5](https://github.com/Liphium/neoroute/commit/b32b8c53c6e68e2f2b0cdd82867a8724ef87cb7e))
* update example to use handshake ([7cb4493](https://github.com/Liphium/neoroute/commit/7cb4493bf6cc88422d67927dcee7f4d7796cdbe1))
* update to SendPing ([7ebdd90](https://github.com/Liphium/neoroute/commit/7ebdd900a2439a1fd7a11785ae35261b73585f50))
* use string from respondError instead of error ([17c48e8](https://github.com/Liphium/neoroute/commit/17c48e8fb231bf179b79ccb22b58ad4eea3adcac))
* Use writers instead of reading the complete byte stream ([e578fcf](https://github.com/Liphium/neoroute/commit/e578fcfad3b99c748aeff3c3e39f70dbd1ec5fc8))
* users no longer create their own sessions and ids and don't interact with ids at all anymore ([4a30453](https://github.com/Liphium/neoroute/commit/4a304535cf3806495c94467ff55f3e256295de5c))
* **workflow:** update push trigger to include all branches ([4cffb36](https://github.com/Liphium/neoroute/commit/4cffb36c596adb175e303606a82cef09c1825059))


### Bug Fixes

* allow pointer for marshaler too ([426732e](https://github.com/Liphium/neoroute/commit/426732e587a7d3f16fe84f2e923b2a5ec26c0cb2))
* calculate coverage properly ([1f068ab](https://github.com/Liphium/neoroute/commit/1f068ab96d670af4bc92ba6bf776cf39a401fc5f))
* change response format to remove error on unmarshal on the client ([2c7d4c6](https://github.com/Liphium/neoroute/commit/2c7d4c64e11b8abb5a67e3bf6c4386b4e625a770))
* fix buildSubroutes to not leave out the full route and clean the route before ([ca4cbc2](https://github.com/Liphium/neoroute/commit/ca4cbc23d170094daf1b1d384729ca2dc078801b))
* force update root module files ([598938b](https://github.com/Liphium/neoroute/commit/598938b8a59898b0dcbf4da6adb5d1dadf7c55ca))
* go pack to previous + better coverage ([3640c2b](https://github.com/Liphium/neoroute/commit/3640c2b345c73473b7157f1771bac65ace36ddab))
* gofmt alignment in router.go ([63430ed](https://github.com/Liphium/neoroute/commit/63430edf7d82315724eeab6a78842873bc146b5c))
* **http:** return error if handshake fails or body cant be read ([5337467](https://github.com/Liphium/neoroute/commit/5337467f9474ed53384cef4ea038b292acf26e39))
* install msgp in pipeline ([b80be7e](https://github.com/Liphium/neoroute/commit/b80be7e4c6c85919fb244edec559ad625bc2bcb1))
* make new function usable by providing bytes directly ([a28b0cc](https://github.com/Liphium/neoroute/commit/a28b0cc829570f30129e7b26dd2fe83e7a5eb126))
* Make routes support pointers due to new generics ([#32](https://github.com/Liphium/neoroute/issues/32)) ([97e4cc2](https://github.com/Liphium/neoroute/commit/97e4cc2f4d4952226303df1e03a677a17eb94e65))
* Make sure middlewares work properly across groups ([47b32a2](https://github.com/Liphium/neoroute/commit/47b32a2540cfbf0747d568bfa4d840faed390428))
* neogen typescript generation errors ([bde5470](https://github.com/Liphium/neoroute/commit/bde5470d161651982cd0d5069391415215ca20ab))
* **neogen:** Correctly extend transporters for TS ([ca9c09f](https://github.com/Liphium/neoroute/commit/ca9c09f988341163da94a2a4b905ed941ac2e8de))
* **neogen:** Generate Uint8Array instead of number[] ([4c2659f](https://github.com/Liphium/neoroute/commit/4c2659f88575f885f73d1156f28d206588d9a6b1))
* **neogen:** Let WebSocketConfig be configurable ([78b7f71](https://github.com/Liphium/neoroute/commit/78b7f71e3809acdbc8e427bc0d022a108fb1aab7))
* **neogen:** nullable generation ([82b7ab8](https://github.com/Liphium/neoroute/commit/82b7ab8a63a33b6ec70a72b838cc1d96be5843cc))
* **neogen:** Structs now generate properly in slices ([86b0de4](https://github.com/Liphium/neoroute/commit/86b0de466d66b3b21f71a9287c09542c39f523ad))
* **neoroute:** actually initialize neos slice to avoid nil pointer ([bff3e79](https://github.com/Liphium/neoroute/commit/bff3e79048fb0cc45acd410db117c7e5f53700dd))
* **neoroute:** return self as well when getNeos is called on NeoRouter ([a6c8863](https://github.com/Liphium/neoroute/commit/a6c886352af34531b724a070005015767031f3e4))
* **neoroute:** the getNeos now return all neos even of sub neos ([cc9e839](https://github.com/Liphium/neoroute/commit/cc9e8392eea70e3d71ff67a9597b9fd902936820))
* properly generate the coverage report for gobadge ([86116b1](https://github.com/Liphium/neoroute/commit/86116b148156bb5e4d73c74db71fc08a77bdbdf5))
* pull rebase ([7470fa9](https://github.com/Liphium/neoroute/commit/7470fa948231cecae11c9233f33309f4a6000310))
* push to correct thingy ([922b494](https://github.com/Liphium/neoroute/commit/922b4943c9883199a151c31971fd999bcc4ce886))
* remove unwanted print statements ([40e20dc](https://github.com/Liphium/neoroute/commit/40e20dc89795c99489be80407871b890765da339))
* runafter was improperly done ([ebefc39](https://github.com/Liphium/neoroute/commit/ebefc3989ea45629c47c23f9725388722ae89000))
* **schema:** Ignore pointers on first type ([97e4cc2](https://github.com/Liphium/neoroute/commit/97e4cc2f4d4952226303df1e03a677a17eb94e65))
* unmarshal message first before unmarshalling event ([a626e33](https://github.com/Liphium/neoroute/commit/a626e336830a139dd694b945341fc95fcf0dc736))
* update .gitignore to not gen files as they are needed when the module is imported ([d30d5ba](https://github.com/Liphium/neoroute/commit/d30d5ba141d60342ed2f1d0590abddbb945f0e11))
* update WebTransport transporter to use the sessions context ([84cf147](https://github.com/Liphium/neoroute/commit/84cf1472c1fe27962df8d573a3e1d44ef20c64a5))
* use newer go version for pipeline ([93d135f](https://github.com/Liphium/neoroute/commit/93d135ff6d905d62bf23646405cc3bc82179ff1f))
* use stuff from docs ([c4e6ee9](https://github.com/Liphium/neoroute/commit/c4e6ee9ec9742b06ca51ae11a1633efbd0c62deb))
* Various issues introduced through 1.0.0 ([97e4cc2](https://github.com/Liphium/neoroute/commit/97e4cc2f4d4952226303df1e03a677a17eb94e65))


### Code Refactoring

* **neogate:** change messageEvent to panic instead of returning an error, as the error should never occur ([a553392](https://github.com/Liphium/neoroute/commit/a55339242954191544b267bc0e987588f48838ee))
* **neoroute:** change route separator from dot to slash ([063b8fc](https://github.com/Liphium/neoroute/commit/063b8fc3caf0ecb6b1ed321efd52f6815729928e))
* **neoroute:** rename the RouteNoop function to RoutePing ([6a6eb5e](https://github.com/Liphium/neoroute/commit/6a6eb5ed993ed0d78062758e87a14da72251c1fe))
* Remove error return from event creation function on ([b128954](https://github.com/Liphium/neoroute/commit/b128954247c5c9f98554504babb1d1a68eb763a9))
* Upgrade everything to Go 1.27 generic methods ([67f0a91](https://github.com/Liphium/neoroute/commit/67f0a91e478ef21047696c4c210d9ea4c2ae4d9a))

## [1.0.1](https://github.com/Liphium/neoroute/compare/v1.0.0...v1.0.1) (2026-09-06)


### Bug Fixes

* Make routes support pointers due to new generics ([#32](https://github.com/Liphium/neoroute/issues/32)) ([97e4cc2](https://github.com/Liphium/neoroute/commit/97e4cc2f4d4952226303df1e03a677a17eb94e65))
* **schema:** Ignore pointers on first type ([97e4cc2](https://github.com/Liphium/neoroute/commit/97e4cc2f4d4952226303df1e03a677a17eb94e65))
* Various issues introduced through 1.0.0 ([97e4cc2](https://github.com/Liphium/neoroute/commit/97e4cc2f4d4952226303df1e03a677a17eb94e65))

## [1.0.0](https://github.com/Liphium/neoroute/compare/v0.8.0...v1.0.0) (2026-08-29)


### ⚠ BREAKING CHANGES

* Upgrade everything to Go 1.27 generic methods

### Features

* Add debugger for neoroute transporters ([#28](https://github.com/Liphium/neoroute/issues/28)) ([47b32a2](https://github.com/Liphium/neoroute/commit/47b32a2540cfbf0747d568bfa4d840faed390428))
* More safety checks for handlers ([9d47b0a](https://github.com/Liphium/neoroute/commit/9d47b0adfd06db5f9314736c79b7c55b08c57c3a))
* **neodebug:** Add debugger for neoroute servers ([47b32a2](https://github.com/Liphium/neoroute/commit/47b32a2540cfbf0747d568bfa4d840faed390428))
* **neoroute:** remove unneeded generics ([d26e477](https://github.com/Liphium/neoroute/commit/d26e4774e3b1623d66b676e02361d71b4a314efe))
* **neoroute:** remove unneeded interface in generics ([400650b](https://github.com/Liphium/neoroute/commit/400650b8d1a5575d3e1f74a4ce4e481d276b6c1f))
* Use writers instead of reading the complete byte stream ([e578fcf](https://github.com/Liphium/neoroute/commit/e578fcfad3b99c748aeff3c3e39f70dbd1ec5fc8))
* **workflow:** update push trigger to include all branches ([4cffb36](https://github.com/Liphium/neoroute/commit/4cffb36c596adb175e303606a82cef09c1825059))


### Bug Fixes

* Make sure middlewares work properly across groups ([47b32a2](https://github.com/Liphium/neoroute/commit/47b32a2540cfbf0747d568bfa4d840faed390428))


### Code Refactoring

* Upgrade everything to Go 1.27 generic methods ([67f0a91](https://github.com/Liphium/neoroute/commit/67f0a91e478ef21047696c4c210d9ea4c2ae4d9a))

## [0.8.0](https://github.com/Liphium/neoroute/compare/v0.7.0...v0.8.0) (2026-08-06)


### ⚠ BREAKING CHANGES

* **neogen:** Generate Uint8Array instead of number[]
* **neogen:** Make the WebSocketConfig one parameter
* **neogen:** Correctly extend transporters for TS
* **neogen:** Let WebSocketConfig be configurable

### Features

* Add functions for adapter management on registry ([249d2d0](https://github.com/Liphium/neoroute/commit/249d2d073c98bcf4734710fc01000e24788993ee))
* **neogen:** Make the WebSocketConfig one parameter ([4fa21f5](https://github.com/Liphium/neoroute/commit/4fa21f550d9b09eb416c8a5f7e57fa6e2c1901d9))


### Bug Fixes

* neogen typescript generation errors ([bde5470](https://github.com/Liphium/neoroute/commit/bde5470d161651982cd0d5069391415215ca20ab))
* **neogen:** Correctly extend transporters for TS ([ca9c09f](https://github.com/Liphium/neoroute/commit/ca9c09f988341163da94a2a4b905ed941ac2e8de))
* **neogen:** Generate Uint8Array instead of number[] ([4c2659f](https://github.com/Liphium/neoroute/commit/4c2659f88575f885f73d1156f28d206588d9a6b1))
* **neogen:** Let WebSocketConfig be configurable ([78b7f71](https://github.com/Liphium/neoroute/commit/78b7f71e3809acdbc8e427bc0d022a108fb1aab7))
* **neogen:** Structs now generate properly in slices ([86b0de4](https://github.com/Liphium/neoroute/commit/86b0de466d66b3b21f71a9287c09542c39f523ad))

## [0.7.0](https://github.com/Liphium/neoroute/compare/v0.6.0...v0.7.0) (2026-07-21)


### ⚠ BREAKING CHANGES

* **neogen:** rename args to command for more control
* **neoroute:** add error return to disconnect function on session
* Remove error return from event creation function on

### Features

* **neogen:** add typescript support ([c9bd471](https://github.com/Liphium/neoroute/commit/c9bd47108578eeda1f94c4f4d37f4c278500439b))
* **neogen:** add typescript support ([0a4634a](https://github.com/Liphium/neoroute/commit/0a4634a0db4d35d1b18f7663bab6d3c17bd4af27))
* **neogen:** rename args to command for more control ([e2965fb](https://github.com/Liphium/neoroute/commit/e2965fb78f7ab1602452da4ff7aee6f992791c3f))
* **neoroute:** add error return to disconnect function on session ([091d50a](https://github.com/Liphium/neoroute/commit/091d50a3d4b9f2aa08f6ff5f4ad89b7e4f18bd9b))
* **neoschema:** Support generation for slices ([266aa3c](https://github.com/Liphium/neoroute/commit/266aa3c08d5a3fb26881971b9e82a306432d469c))


### Code Refactoring

* Remove error return from event creation function on ([b128954](https://github.com/Liphium/neoroute/commit/b128954247c5c9f98554504babb1d1a68eb763a9))

## [0.6.0](https://github.com/Liphium/neoroute/compare/v0.5.0...v0.6.0) (2026-07-17)


### ⚠ BREAKING CHANGES

* **neoroute:** allow multiple middlewares for the same route

### Features

* **neoroute:** allow multiple middlewares for the same route ([29ab932](https://github.com/Liphium/neoroute/commit/29ab9321325402c6901e8c3f2053377c1ffcaf87))


### Bug Fixes

* gofmt alignment in router.go ([63430ed](https://github.com/Liphium/neoroute/commit/63430edf7d82315724eeab6a78842873bc146b5c))

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
