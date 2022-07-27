import del from 'del'
import { dir } from '../utils/path'

export function clean(): Promise<string[]> {
  return del(dir('build'))
}

export function cleanDist(): Promise<string[]> {
  return del(dir('dist'))
}
