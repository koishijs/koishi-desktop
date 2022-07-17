import { resolve } from 'path'

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
