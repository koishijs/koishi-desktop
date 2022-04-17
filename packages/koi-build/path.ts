import { resolve as pathResolve } from 'path'

const resolveMap = {
  root: '../../',

  build: '../../build/',
  buildTemp: '../../build/tmp/',
  dist: '../../build/koi/',
  distPre: '../../build/dist-pre/',

  koi: '../koi/',
}

export type ResolveRoot = keyof typeof resolveMap

export function resolve(path: string, root: ResolveRoot = 'root') {
  return pathResolve(__dirname, resolveMap[root], path)
}
