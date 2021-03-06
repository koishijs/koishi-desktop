import { resolve } from 'path'

/**
 * Paths used in building.
 *
 * Please check tasks/prepare/folder.ts after
 * modification to ensure all folder created
 * during prepareFolder task.
 */
const buildPaths = {
  src: 'packages/cli',
  srcBuild: 'packages/buildsrc',

  build: 'build/',
  buildCache: 'build/caches',

  buildPortable: 'build/varients/portable/',
  buildPortableData: 'build/varients/portable/data/',

  dist: 'build/dist/',
}

export type BuildPaths = keyof typeof buildPaths

export const dir = (base: BuildPaths, path = ''): string =>
  resolve(__dirname, '../../../', base, path)
