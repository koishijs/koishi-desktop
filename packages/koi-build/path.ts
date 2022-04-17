import { resolve as pathResolve } from 'path'

const resolveMap = {
  root: '../../',
  koi: '../koi/',
}

export type ResolveRoot = keyof typeof resolveMap

export function resolve(path: string, root: ResolveRoot = 'root') {
  return pathResolve(__dirname, resolveMap[root], path)
}
