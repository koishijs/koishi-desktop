import { resolve } from 'path'

/**
 * Paths used in building.
 *
 * Please check tasks/prepare/folder.ts after
 * modification to ensure all folder created
 * during prepareFolder task.
 */
const buildPaths = {
  root: './',
  locales: './locales/',

  packages: 'packages/',
  src: 'packages/app/',
  srcCore: 'packages/core/',
  srcUnfold: 'packages/unfold/',
  srcShellMac: 'packages/shellmac/',
  srcShellWin: 'packages/shellwin/',
  srcShellLinux: 'packages/shelllinux/',
  srcAssets: 'packages/assets/',
  srcIcon: 'packages/core/ui/icon/',
  srcBuild: 'packages/build/',
  srcUserscript: 'packages/userscript/',
  templates: 'packages/build/templates/',

  build: 'build/',
  buildAssets: 'build/assets',
  buildCache: 'build/caches/',
  buildVendor: 'build/vendor/',
  buildResources: 'build/resources/',

  buildPortable: 'build/varients/portable/',
  buildPortableBin: 'build/varients/portable/bin/',
  buildPortableData: 'build/varients/portable/data/',
  buildUnfold: 'build/varients/unfold/',
  buildUnfoldData: 'build/varients/unfold/portabledata/',
  buildUnfoldBinary: 'build/varients/unfold/binary/',
  buildShellWin: 'build/varients/shellwin/',
  buildShellLinux: 'build/varients/shelllinux/',
  buildMac: 'build/varients/mac/',
  buildMsi: 'build/varients/msi/',
  buildLinux: 'build/varients/linux/',

  dist: 'build/dist/',
}

export type BuildPaths = keyof typeof buildPaths

export const dir = (base: BuildPaths, path = ''): string =>
  resolve(__dirname, '../../../', buildPaths[base], path)
