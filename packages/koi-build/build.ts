import del from 'del'
import * as fs from 'fs'
import { series } from 'gulp'
import { error } from 'gulplog'
import {
  boilerplateVersion,
  defaultInstance,
  defaultKoiConfig,
  defaultNpmrc,
  getKoiVersion,
} from './config'
import { resolve } from './path'
import { spawnAsync } from './utils'

export async function goModDownload() {
  const result = await spawnAsync('go', ['mod', 'download'], {
    cwd: resolve('.', 'koi'),
  })
  if (result) {
    const err = `'go mod download' exited with error code: ${result}`
    error(err)
    throw new Error(err)
  }
}

export async function goBuild() {
  const result = await spawnAsync(
    'go',
    [
      'build',
      '-o',
      resolve(process.platform === 'win32' ? 'koi.exe' : '', 'dist'),
      '-x',
      '-v',
      '-trimpath',
      '-ldflags',
      `-w -s -X koi/config.Version=${await getKoiVersion()}`,
    ],
    { cwd: resolve('.', 'koi') }
  )
  if (result) {
    const err = `'go build' exited with error code: ${result}`
    error(err)
    throw new Error(err)
  }
}

export const buildExe = series(goModDownload, goBuild)

export async function writeConfig() {
  await fs.promises.writeFile(resolve('home/.npmrc', 'distData'), defaultNpmrc)
  await fs.promises.writeFile(resolve('koi.yml', 'dist'), defaultKoiConfig)
}

export async function createDefaultInstance() {
  const result = await spawnAsync(
    process.platform === 'win32' ? 'koi' : './koi',
    [
      'instance',
      'create',
      '-n',
      defaultInstance,
      '-p',
      '@koishijs/plugin-database-sqlite',
      '-r',
      boilerplateVersion,
    ],
    { cwd: resolve('.', 'dist') }
  )
  if (result) {
    const err = `'koi instance create' exited with error code: ${result}`
    error(err)
    throw new Error(err)
  }
}

export async function cleanupDefaultInstance(): Promise<void> {
  del(resolve('home', 'distData'))
  del(resolve('tmp', 'distData'))
}

export async function run() {
  const result = await spawnAsync(
    process.platform === 'win32' ? 'koi' : './koi',
    [],
    { cwd: resolve('.', 'dist') }
  )
  if (result) {
    const err = `'koi' exited with error code: ${result}`
    error(err)
    throw new Error(err)
  }
}

export const build = series(
  buildExe,
  writeConfig,
  createDefaultInstance,
  cleanupDefaultInstance,
  writeConfig
)
