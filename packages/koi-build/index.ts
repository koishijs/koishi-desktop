export * from './prepareNode'

import del from 'del'
import { series } from 'gulp'
import mkdirp from 'mkdirp'
import { resolve } from './path'
import { prepareNode } from './prepareNode'

export function clean(): Promise<string[]> {
  return del(resolve('build', 'root'))
}

export function cleanDist(): Promise<string[]> {
  return del(resolve('build/koi', 'root'))
}

export async function mkdir() {
  await mkdirp(resolve('.', 'buildTemp'))
  await mkdirp(resolve('.', 'dist'))
  await mkdirp(resolve('.', 'distData'))
  await mkdirp(resolve('node', 'distData'))
}

export const all = series(clean, mkdir, prepareNode)

export default all
