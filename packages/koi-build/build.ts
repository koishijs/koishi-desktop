import * as fs from 'fs'
import { series } from 'gulp'
import { error } from 'gulplog'
import {
  defaultInstance,
  defaultKoiConfig,
  defaultNpmrc,
  getKoiVersion,
} from './config'
import { resolve } from './path'
import { spawnAsync } from './utils'

export async function buildExe() {
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
      `-w -s -X koi/main.Version=${await getKoiVersion()}`,
    ],
    { cwd: resolve('.', 'koi') }
  )
  if (result) {
    const err = `'go build' exited with error code: ${result}`
    error(err)
    throw new Error(err)
  }
}

export async function writeConfig() {
  await fs.promises.writeFile(resolve('home/.npmrc', 'distData'), defaultNpmrc)
  await fs.promises.writeFile(resolve('koi.yml', 'dist'), defaultKoiConfig)
}

export async function createDefaultInstance() {
  const result = await spawnAsync(
    'koi',
    [
      'instance',
      'create',
      '-n',
      defaultInstance,
      '-p',
      '@koishijs/plugin-database-sqlite',
    ],
    { cwd: resolve('.', 'dist') }
  )
  if (result) {
    const err = `'koi instance create' exited with error code: ${result}`
    error(err)
    throw new Error(err)
  }
}

export const build = series(buildExe, writeConfig, createDefaultInstance)
