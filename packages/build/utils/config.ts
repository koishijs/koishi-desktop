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
  const build = spl.length > 1 ? (spl[1] === 'dirty' ? 0 : Number(spl[1])) : 0
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
  Comments: 'Cordis Desktop',
  CompanyName: 'Cordis',
  FileDescription: 'Cordis Desktop',
  FileVersion: koiVersion,
  InternalName: 'Cordis Desktop',
  LegalCopyright: `${new Date().getFullYear()} Cordis Team`,
  LegalTrademarks: `${new Date().getFullYear()} Cordis Team`,
  OriginalFilename: 'koi',
  PrivateBuild: koiVersion,
  ProductName: 'Cordis Desktop',
  ProductVersion: koiVersion,
  SpecialBuild: koiVersion,
}

export const koishiVersionStrings = {
  Comments: 'Cordis',
  CompanyName: 'Cordis',
  FileDescription: 'Cordis',
  FileVersion: koiVersion,
  InternalName: 'Cordis',
  LegalCopyright: `${new Date().getFullYear()} Cordis Team`,
  LegalTrademarks: `${new Date().getFullYear()} Cordis Team`,
  OriginalFilename: 'cordis',
  PrivateBuild: koiVersion,
  ProductName: 'Cordis',
  ProductVersion: koiVersion,
  SpecialBuild: koiVersion,
}

export const koiVersionStringsJson = JSON.stringify(koiVersionStrings)

//#endregion

export const goEnv = {
  GOOS: spawnSyncOutput('go', ['env', 'GOOS']),
  GOARCH: spawnSyncOutput('go', ['env', 'GOARCH']),
}
