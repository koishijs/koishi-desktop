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

  packages: 'packages/',
  src: 'packages/app/',
  srcCore: 'packages/core/',
  srcAssets: 'packages/assets/',
  srcIcon: 'packages/core/ui/icon/',
  srcBuild: 'packages/build/',
  templates: 'packages/build/templates/',

  build: 'build/',
  buildAssets: 'build/assets',
  buildCache: 'build/caches/',
  buildVendor: 'build/vendor/',

  buildPortable: 'build/varients/portable/',
  buildPortableData: 'build/varients/portable/data/',
  buildMac: 'build/varients/mac/',

  dist: 'build/dist/',
}

export type BuildPaths = keyof typeof buildPaths

export const dir = (base: BuildPaths, path = ''): string =>
  resolve(__dirname, '../../../', buildPaths[base], path)
