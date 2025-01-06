# Changelog

## 0.1.0 (2025-01-06)


### 🎉 Features

* add htmx support ([#2](https://github.com/zerok/samara/issues/2)) ([172691b](https://github.com/zerok/samara/commit/172691b1f5e33672871d34c3f4c208385e6325b6))
* add thread cache and configurable logging ([#9](https://github.com/zerok/samara/issues/9)) ([4c759e1](https://github.com/zerok/samara/commit/4c759e1077b47fe406bc64465dfb77ff0268c598))
* add Valkey support for oop caching ([#24](https://github.com/zerok/samara/issues/24)) ([0e94032](https://github.com/zerok/samara/commit/0e940323612bbe5bb4d4cc162f9cb93a66b79e58))
* expose otel traces through valkeyotel ([#26](https://github.com/zerok/samara/issues/26)) ([821acd5](https://github.com/zerok/samara/commit/821acd5c047dc67c6321f24b16e059410dd0c82b))
* integrate opentelemetry tracing ([#18](https://github.com/zerok/samara/issues/18)) ([afa4f66](https://github.com/zerok/samara/commit/afa4f66f709cfc942dd3de1b1cf3b774041c957c))
* proxy avatars ([#23](https://github.com/zerok/samara/issues/23)) ([4ad6f11](https://github.com/zerok/samara/commit/4ad6f11bb10a8951fa2c30b270c179376ac26ee8))
* remove favorited_by endpoint ([#19](https://github.com/zerok/samara/issues/19)) ([ce04dfc](https://github.com/zerok/samara/commit/ce04dfc50f7f11641e6b09d036ef7699a6770f64))
* render thread as HTML or JSON ([#20](https://github.com/zerok/samara/issues/20)) ([160e741](https://github.com/zerok/samara/commit/160e74128bb1ca33f3f0b1b2ecf36d3dbe3ae31d))


### 🐛 Bug Fixes

* allow raw handles for incoming URIs ([#7](https://github.com/zerok/samara/issues/7)) ([f134a08](https://github.com/zerok/samara/commit/f134a08aee4607b635d38b6da625456f342d2cbc))
* disable attestation until this repository is public ([#12](https://github.com/zerok/samara/issues/12)) ([52e749b](https://github.com/zerok/samara/commit/52e749bea558d89e384a2a02f0e545b6b1fc975b))
* docker image push + attestation ([#11](https://github.com/zerok/samara/issues/11)) ([ace285c](https://github.com/zerok/samara/commit/ace285cb05b59ce19122c7e2938ba9de4d3254ac))
* improve htmx template ([#8](https://github.com/zerok/samara/issues/8)) ([0a8a4ca](https://github.com/zerok/samara/commit/0a8a4caf093da46dde5b9fbab08502597b599fd9))
* trace context was wrong to bsky api ([#22](https://github.com/zerok/samara/issues/22)) ([e28f367](https://github.com/zerok/samara/commit/e28f367b76ae75bc1bb5ce265c4a3b6b0e2c31f8))


### 📝 Documentation

* update README for caching ([#25](https://github.com/zerok/samara/issues/25)) ([fcd8e72](https://github.com/zerok/samara/commit/fcd8e72599acc2d966e8fb7b8ffec182646d8b39))


### 🤖 Continuous Integration

* add dependabot configuration ([#3](https://github.com/zerok/samara/issues/3)) ([60aeda0](https://github.com/zerok/samara/commit/60aeda05551def9dc3241d14f388be24e4b45031))
* add release-please to automate releases ([#28](https://github.com/zerok/samara/issues/28)) ([1eb5ebd](https://github.com/zerok/samara/commit/1eb5ebd70eaf98b8993aea27f55ea107a159767e))
* add test workflow ([#1](https://github.com/zerok/samara/issues/1)) ([1cad130](https://github.com/zerok/samara/commit/1cad1308816c448a78960ee9c23c939870e5a888))
* build Docker image ([#10](https://github.com/zerok/samara/issues/10)) ([207dfc7](https://github.com/zerok/samara/commit/207dfc7dd281f56d8a622088a165fa8c30247bc1))
* fix typo in release-please workflow ([#29](https://github.com/zerok/samara/issues/29)) ([8bbe834](https://github.com/zerok/samara/commit/8bbe8342ccade201b73887ce95b47e000c23a722))
* give release-please packages-write permissions ([#30](https://github.com/zerok/samara/issues/30)) ([7c9f8c1](https://github.com/zerok/samara/commit/7c9f8c1fbbaf5e7401545b6d4ea97fae17c809c1))
* more permissions for release-please docker integration ([#31](https://github.com/zerok/samara/issues/31)) ([f8dd0ac](https://github.com/zerok/samara/commit/f8dd0ac3a395fc53a3aceb5ca42f3e10339bf554))
* run tests on more pull_request events ([#34](https://github.com/zerok/samara/issues/34)) ([d3fe3db](https://github.com/zerok/samara/commit/d3fe3dbcebfd6a0022cf7695f340603dbcf2da5c))


### 🔧 Miscellaneous Chores

* add license ([#6](https://github.com/zerok/samara/issues/6)) ([7931ef7](https://github.com/zerok/samara/commit/7931ef7010ea1802021787738992979c1ecd046f))
* **deps:** bump actions/setup-go from 5.1.0 to 5.2.0 ([#16](https://github.com/zerok/samara/issues/16)) ([4847796](https://github.com/zerok/samara/commit/48477964f76bad0b89527a9e595a31bd230e0448))
* **deps:** bump alpine from 3.20 to 3.21 ([#15](https://github.com/zerok/samara/issues/15)) ([cb1d5ec](https://github.com/zerok/samara/commit/cb1d5eccc3b9afa8428bd106c7f103c2e0a4925f))
* **deps:** bump docker/login-action from 3.2.0 to 3.3.0 ([#13](https://github.com/zerok/samara/issues/13)) ([1f40c46](https://github.com/zerok/samara/commit/1f40c46296bb246e1241e999e1ef178b95d38527))
* **deps:** bump docker/setup-buildx-action from 3.7.1 to 3.8.0 ([#17](https://github.com/zerok/samara/issues/17)) ([4a40819](https://github.com/zerok/samara/commit/4a40819853d7ec246dc53deb4f027dfae28ff2f7))
* **deps:** bump github.com/stretchr/testify from 1.9.0 to 1.10.0 ([#5](https://github.com/zerok/samara/issues/5)) ([17e631e](https://github.com/zerok/samara/commit/17e631e8c4ccd7cdc20b08b7223eea54535540da))
* **deps:** bump github.com/yuin/goldmark from 1.3.5 to 1.7.8 ([#4](https://github.com/zerok/samara/issues/4)) ([a7c1c1b](https://github.com/zerok/samara/commit/a7c1c1b02d95a77b2665c08935adba1d20bd1c9f))
* **deps:** bump golang from 1.23.3-alpine to 1.23.4-alpine ([#14](https://github.com/zerok/samara/issues/14)) ([fb81ed8](https://github.com/zerok/samara/commit/fb81ed89ed4c2ad2f5a68bc08c50f123b24c6d8c))
* **deps:** bump golang.org/x/crypto from 0.30.0 to 0.31.0 ([#27](https://github.com/zerok/samara/issues/27)) ([a9aec31](https://github.com/zerok/samara/commit/a9aec3131efaecc0a5f57ccffb625488cc7bb701))
* initial implementation ([e0926bd](https://github.com/zerok/samara/commit/e0926bd924e040062a27ee71ed3ee501f5aa95f8))
* release 0.1.0 ([#33](https://github.com/zerok/samara/issues/33)) ([26ca502](https://github.com/zerok/samara/commit/26ca5026c0bb5636baee49b09c8019e13c2b09f8))
* update README ([#21](https://github.com/zerok/samara/issues/21)) ([0608c9c](https://github.com/zerok/samara/commit/0608c9c0e348dd96794027f9d23b472fd7780a4c))
