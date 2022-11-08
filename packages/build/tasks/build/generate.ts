import { parallel, series } from 'gulp'
import fs from 'node:fs/promises'
import {
  koiConfig,
  koiManifest,
  koiShellResources,
  koiVersionInfo,
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

export const generateKoiShellResources = () =>
  fs.writeFile(dir('srcShellWin', 'src/koishell.rc'), koiShellResources)

export const generateVisualElementsManifest = async () => {
  await fs.copyFile(
    dir('templates', 'portable/koi.VisualElementsManifest.xml'),
    dir('buildPortable', 'koi.VisualElementsManifest.xml')
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
    generateKoiShellResources,
    process.platform === 'win32'
      ? series(generateAssets, generateVisualElementsManifest)
      : generateAssets
  ),
  i18nGenerate
)
