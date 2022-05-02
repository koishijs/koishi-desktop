export * from './build'
export * from './pack'
export * from './prepareNode'

import del from 'del'
import { series } from 'gulp'
import mkdirp from 'mkdirp'
import { build, buildExe, run } from './build'
import { pack } from './pack'
import { resolve } from './path'
import { prepareNode } from './prepareNode'

export function clean(): Promise<string[]> {
  return del(resolve('.', 'build'))
}

export function cleanDist(): Promise<string[]> {
  return del(resolve('.', 'dist'))
}

export async function mkdir(): Promise<void> {
  await mkdirp(resolve('.', 'buildTemp'))
  await mkdirp(resolve('home', 'distData'))
  await mkdirp(resolve('instances', 'distData'))
  await mkdirp(resolve('node', 'distData'))
  await mkdirp(resolve('tmp', 'distData'))
}

export const dev = series(mkdir, prepareNode, build)

export const test = series(mkdir, prepareNode, buildExe, run)

export const all = series(clean, dev, pack)

export default dev
