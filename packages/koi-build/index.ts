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

export async function mkdir(): Promise<void> {
  await mkdirp(resolve('.', 'buildTemp'))
  await mkdirp(resolve('home', 'distData'))
  await mkdirp(resolve('node', 'distData'))
  await mkdirp(resolve('temp', 'distData'))
  await mkdirp(resolve('.', 'defaultInstance'))
}

export const dev = series(mkdir, prepareNode)

export const all = series(clean, dev)

export default dev
