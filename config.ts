//            KOISHI DESKTOP BUILD CONFIG
// ---------------------------------------------------
// Here are all config used during building. You can
// freely change configs like mirrors and build paths,
// but do not commit/push them.

//#region Sources
// These are sources for downloading toolchains.
// Try using mirrors if you cannot download some of
// them.

export const sourceNode = 'https://nodejs.org/dist'
export const sourceYarn = 'https://repo.yarnpkg.com'
export const sourceGitHub = 'https://github.com'

//#endregion

//#region Toolchain
// These are used to prepare toolchains.
// Remember to test all tasks before upgrading
// versions.

export const versionNode = '16.17.0'
export const versionYarn = '3.2.2'

export const versionToolsVersioninfo = 'v1.4.0'
export const versionToolsGolangCILint = 'v1.48.0'

//#endregion

//#region Defaults
// These are defaults for koishi-desktop.

export const repoBoilerplate = 'koishijs/boilerplate'
export const versionBoilerplate = 'v1.0.5'

/**
 * ID of the default instance.
 */
export const defaultInstance = 'default'

//#endregion

//#region Overrides
// Enable these config overrides to test edge case
// behaviors.

export const overrideKoiVersion = '' // '0.1.0'

//#endregion
