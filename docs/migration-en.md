# 从低版本迁移

## Config "mode" Migration

Migration roadmap:

- 必须迁移: At v0.2.0

::: tip
This migration is covered by a newer migration below and can be ignored.
:::

The `mode` config has new values in v0.2.0. When upgrading, you need to manually delete the `mode` config in `koi.yml` if you changed it before.

## Start Script Migration

Migration roadmap:

- 必须迁移: At v0.7.1

As of v0.7.1, Koishi Desktop no longer calls Koishi CLI directly, but instead uses the startup script in the instance `package.json` file. When upgrading, you need to manually add the following line to the `package.json` of all instances:

```diff
 {
   "name": "@koishijs/boilerplate",
   "version": "1.1.0",
   "private": true,
+  "scripts": { "start": "koishi start" },
   "files": [
     ".env",
     "koishi.yml"
   ],
```

使用新版本 Koishi Desktop 或新版本整合包则无需进行迁移。

## 数据目录迁移

Migration roadmap:

- 开始迁移: From v0.9.0 (自动)
- 必须迁移: 无需

自 v0.9.0 起，我们将本机安装版本的 Koishi Desktop 的默认数据目录由 `Il Harper/Koishi` 变更为 `Koishi/Desktop`，以规范产品。迁移会自动进行，无需手动操作。

在迁移成功后，你可以手动移除用户数据目录内的 `Il Harper/Koishi` 文件夹。

## 执行环境迁移

Migration roadmap:

- 开始迁移: From v0.10.0 (自动)
- 必须迁移: 无需

自 v0.10.0 起，我们使用位于 Koishi Desktop 程序内的 Node 环境运行 Koishi，而非数据内的。迁移会自动进行，无需手动操作。

在迁移成功后，你可以手动移除数据文件夹内的 `node` 文件夹。

## Remove Config "mode"

Migration roadmap:

- Deprecated: From v0.2.0
- 不再有效: From v1.0.0
- 必须迁移: From v1.0.1

As of v0.2.0, the `mode` configuration item has been deprecated and is scheduled to be removed. To start daemon, use `koi run daemon`.

You need to delete the `mode` config in `koi.yml` as soon as possible if you changed it before.

自 v1.0.0 起，这一配置项已失效。你无法使用这一配置项改变 Koishi Desktop 的工作模式。

自 v1.0.1 起，这一配置项已被移除。使用带有这一配置项的配置文件启动 Koishi Desktop 将会使 Koishi Desktop 报错。

## 包管理器迁移

Migration roadmap:

- 生效: From v0.2.0
- 开始迁移: From v0.11.1
- 必须迁移: 无需

在以往版本的 Koishi Desktop，下载和解析包所使用的包管理器 Yarn 由 Koishi Desktop 本身提供。这导致了一个问题：如果实例文件夹在不同计算机之间迁移，且两台计算机间的 Koishi Desktop 提供的包管理器的版本不同的话，下载和解析包很有可能失败。「实例内包管理器」是 Yarn 提出的一种解决方案，通过在 Koishi 项目实例内保存包管理器程序的方法解决了不同环境下包管理器版本不同的问题。

自 v0.2.0 起，Koishi Desktop 会优先使用实例内的包管理器。

自 v0.11.1 起，Koishi Desktop 推荐任何实例均使用实例内包管理器，并且为默认实例内置了包管理器。

如果你的现有实例尚未启用实例内包管理器，则可以在 Koishi 终端内输入下面的命令启用：

```sh
koi yarn -n <实例名> set version berry
```

这样，该实例将会使用固定版本的包管理器，不会在计算机间迁移时遇到问题。
