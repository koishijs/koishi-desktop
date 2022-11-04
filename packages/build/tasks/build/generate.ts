import { parallel, series } from 'gulp'
import fs from 'node:fs/promises'
import {
  koiConfig,
  koiManifest,
  koiVersionInfo,
  koiVisualElementsManifest,
} from '../../templates'
import { dir } from '../../utils/path'
import { i18nGenerate } from '../i18n'
import { generateAssets } from './assets'

export const generateKoiConfig = () =>
  fs.writeFile(dir('buildPortable', 'koi.yml'), koiConfig)

export const generateKoiVersionInfo = () =>
  fs.writeFile(dir('src', 'versioninfo.json'), koiVersionInfo)

export const generateKoiManifest = () =>
  fs.writeFile(dir('src', 'koi.exe.manifest'), koiManifest)

export const generateVisualElementsManifest = async () => {
  await fs.writeFile(
    dir('buildPortable', 'koi.VisualElementsManifest.xml'),
    koiVisualElementsManifest
  )
  await fs.copyFile(
    dir('buildAssets', 'koishi-tile.png'),
    dir('buildPortable', 'koishi.png')
  )
}

export const generate = series(
  parallel(
    generateKoiConfig,
    generateKoiVersionInfo,
    generateKoiManifest,
    series(generateAssets, generateVisualElementsManifest)
  ),
  i18nGenerate
)
