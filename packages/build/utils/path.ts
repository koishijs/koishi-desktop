import { resolve } from 'path'

/**
 * Paths used in building.
 *
 * Please check tasks/prepare/folder.ts after
 * modification to ensure all folder created
 * during prepareFolder task.
 */
const buildPaths = {
  root: '.',

  packages: 'packages',
  src: 'packages/app',
  srcBuild: 'packages/build',
  templates: 'packages/build/templates',

  build: 'build/',
  buildCache: 'build/caches',

  buildPortable: 'build/varients/portable/',
  buildPortableData: 'build/varients/portable/data/',

  dist: 'build/dist/',
}

export type BuildPaths = keyof typeof buildPaths

export const dir = (base: BuildPaths, path = ''): string =>
  resolve(__dirname, '../../../', buildPaths[base], path)
