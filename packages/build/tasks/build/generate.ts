import * as fs from 'fs'
import { parallel } from 'gulp'
import mkdirp from 'mkdirp'
import {
  koiConfig,
  koiConfigBefore,
  koiManifest,
  koiVersionInfo,
} from '../../templates'
import { dir } from '../../utils/path'
import { generateAssets } from './assets'

export const generateKoiConfigBefore = () =>
  fs.promises.writeFile(dir('buildPortable', 'koi.yml'), koiConfigBefore)

export const generateKoiConfig = () =>
  fs.promises.writeFile(dir('buildPortable', 'koi.yml'), koiConfig)

export const generateKoiVersionInfo = () =>
  fs.promises.writeFile(dir('src', 'versioninfo.json'), koiVersionInfo)

export const generateKoiManifest = async () => {
  await mkdirp(dir('src', 'resources'))
  await fs.promises.writeFile(
    dir('src', 'resources/koi.exe.manifest'),
    koiManifest
  )
}

export const generateBefore = parallel(
  generateKoiConfigBefore,
  generateKoiVersionInfo,
  generateKoiManifest,
  generateAssets
)

export const generateAfter = parallel(generateKoiConfig)
