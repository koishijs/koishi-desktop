import del from 'del'
import mkdirp from 'mkdirp'
import { resolve } from './path'

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
