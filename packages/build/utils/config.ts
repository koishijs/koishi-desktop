import Handlebars from 'handlebars'
import * as fs from 'node:fs'
import { overrideKoiVersion } from '../../../config'
import { defaultInstance } from './config'
import { dir } from './path'
import { spawnSyncOutput } from './spawn'

export * from '../../../config'

//#region Version

export interface KoiSemver {
  major: number
  minor: number
  patch: number
  build: number
}

const buildKoiVersion = () => {
  if (overrideKoiVersion) return overrideKoiVersion

  try {
    return spawnSyncOutput('git', ['describe', '--tags', '--dirty']).trim()
  } catch (error) {
    return '0.0.1'
  }
}

export const koiVersion = buildKoiVersion()

const buildKoiSemver = () => {
  const spl = koiVersion.split('-')
  const build = spl.length > 1 ? Number(spl[1]) : 0
  const majorMinorPatch = spl[0].slice(1).split('.')

  return {
    major: Number(majorMinorPatch[0]),
    minor: Number(majorMinorPatch[1]),
    patch: Number(majorMinorPatch[2]),
    build,
  }
}

export const koiSemver = buildKoiSemver()

//#endregion

//#region Templates

export const koiConfigBefore = Handlebars.compile(
  fs
    .readFileSync(dir('templates', 'koi-config-before.yml.hbs'))
    .toString('utf-8')
)({})

export const koiConfig = Handlebars.compile(
  fs.readFileSync(dir('templates', 'koi-config.yml.hbs')).toString('utf-8')
)({ defaultInstance })

export const koiVersionInfo = Handlebars.compile(
  fs.readFileSync(dir('templates', 'versioninfo.json.hbs')).toString('utf-8')
)({ koiVersion, koiSemver })

export const koiManifest = Handlebars.compile(
  fs.readFileSync(dir('templates', 'koi.exe.manifest.hbs')).toString('utf-8')
)({ koiSemver })

export const koishiManifest = Handlebars.compile(
  fs.readFileSync(dir('templates', 'koishi.exe.manifest.hbs')).toString('utf-8')
)({ koiSemver })

//#endregion
