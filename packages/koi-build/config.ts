import { info } from 'gulplog'
import { spawnOut } from './utils'

export const nodeVersion = '14.19.1'
export const yarnVersion = '3.2.0'

export const defaultInstance = 'adace8ea4130c619a7376e8e117780102e67dca7'

export const defaultYarnrc = `
npmRegistryServer: https://registry.npmmirror.com/
nodeLinker: node-modules
plugins:
  - path: ../../node/plugin-workspace-tools.cjs
    spec: "@yarnpkg/plugin-workspace-tools"
`.trim()

export const defaultKoiConfig = `
mode: portable
target: ${defaultInstance}
`.trim()

export const boilerplateVersion = 'a48d1e81653656509a6364ba50959398b86eb5e8'

let koiVersionTemp = ''

export async function getKoiVersion(): Promise<string> {
  if (koiVersionTemp) return koiVersionTemp
  try {
    koiVersionTemp = (await spawnOut('git', ['describe', '--tags'])).trim()
  } catch (error) {
    koiVersionTemp = 'v0.0.1'
  }
  info(`Use koi version ${koiVersionTemp}`)
  return koiVersionTemp
}

export interface KoiSemVer {
  major: number
  minor: number
  patch: number
  build: number
}

let koiSemverTemp: KoiSemVer | null = null

export async function getKoiSemVer(): Promise<KoiSemVer> {
  if (koiSemverTemp) return koiSemverTemp

  const raw = await getKoiVersion()
  const spl = raw.split('-')
  let build = 0

  if (spl.length > 1) {
    // Parse build number
    build = Number(spl[1])
  }

  const majorMinorPatch = spl[0].slice(1).split('.')

  koiSemverTemp = {
    major: Number(majorMinorPatch[0]),
    minor: Number(majorMinorPatch[1]),
    patch: Number(majorMinorPatch[2]),
    build,
  }

  return koiSemverTemp
}
