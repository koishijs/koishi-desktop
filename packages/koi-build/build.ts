import del from 'del'
import * as fs from 'fs'
import { series } from 'gulp'
import { error } from 'gulplog'
import mkdirp from 'mkdirp'
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
  await del(resolve('home', 'distData'))
  await del(resolve('tmp', 'distData'))
  await del(resolve('node/lib/node_modules/npm', 'distData'))
  await del(resolve('node/lib/node_modules/corepack', 'distData'))
  await del(resolve('node/bin/corepack', 'distData'))
  await del(resolve('node/bin/npm', 'distData'))
  await del(resolve('node/bin/npx', 'distData'))
  await del(resolve('node/share', 'distData'))
  const cwd = process.cwd()
  process.chdir(resolve('node/bin', 'distData'))
  fs.symlinkSync(
    `../../instances/${defaultInstance}/node_modules/.bin/yarn`,
    'yarn'
  )
  process.chdir(cwd)
  await mkdirp(resolve('home', 'distData'))
  await mkdirp(resolve('tmp', 'distData'))
  const deps = JSON.parse(
    await fs.promises.readFile(
      resolve(`package.json`, 'defaultInstance'),
      'utf-8'
    )
  )
  delete deps.devDependencies
  await fs.promises.writeFile(
    resolve('package.json', 'defaultInstance'),
    JSON.stringify(deps, null, 2)
  )
  await spawnAsync('npx', ['yarn', '--production'], {
    cwd: resolve('.', 'defaultInstance'),
  })
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
function execSync() {
  throw new Error('Function not implemented.')
}
