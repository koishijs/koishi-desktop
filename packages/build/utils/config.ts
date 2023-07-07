import { overrideKoiVersion } from '../../../config'
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

const buildKoiBuildNumber = () => {
  try {
    return spawnSyncOutput('git', ['rev-list', '--count', 'HEAD']).trim()
  } catch (error) {
    return '0'
  }
}

export const koiBuildNumber = buildKoiBuildNumber()

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

export const koiVersionStrings = {
  Comments: 'Koishi Desktop',
  CompanyName: 'Koishi.js',
  FileDescription: 'Koishi Desktop',
  FileVersion: koiVersion,
  InternalName: 'Koishi Desktop',
  LegalCopyright: '2022 Koishi.js Team',
  LegalTrademarks: '2022 Koishi.js Team',
  OriginalFilename: 'koi',
  PrivateBuild: koiVersion,
  ProductName: 'Koishi Desktop',
  ProductVersion: koiVersion,
  SpecialBuild: koiVersion,
}

export const koishiVersionStrings = {
  Comments: 'Koishi',
  CompanyName: 'Koishi.js',
  FileDescription: 'Koishi',
  FileVersion: koiVersion,
  InternalName: 'Koishi',
  LegalCopyright: '2022 Koishi.js Team',
  LegalTrademarks: '2022 Koishi.js Team',
  OriginalFilename: 'koishi',
  PrivateBuild: koiVersion,
  ProductName: 'Koishi',
  ProductVersion: koiVersion,
  SpecialBuild: koiVersion,
}

export const koiVersionStringsJson = JSON.stringify(koiVersionStrings)

//#endregion

export const goEnv = {
  GOOS: spawnSyncOutput('go', ['env', 'GOOS']),
  GOARCH: spawnSyncOutput('go', ['env', 'GOARCH']),
}
