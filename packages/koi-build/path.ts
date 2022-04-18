import { resolve as pathResolve } from 'path'
import { defaultInstance } from './config'

const resolveMap = {
  root: '../../',

  build: '../../build/',
  buildTemp: '../../build/tmp/',
  dist: '../../build/koi/',
  distData: '../../build/koi/data/',
  defaultInstance: `../../build/koi/data/instances/${defaultInstance}/`,

  koi: '../koi/',
}

export type ResolveRoot = keyof typeof resolveMap

export function resolve(path: string, root: ResolveRoot = 'root') {
  return pathResolve(__dirname, resolveMap[root], path)
}
