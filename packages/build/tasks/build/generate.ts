import * as fs from 'fs'
import { parallel } from 'gulp'
import { koiConfig, koiManifest, koiVersionInfo } from '../../templates'
import { dir } from '../../utils/path'
import { i18nGenerate } from '../i18n'
import { generateAssets } from './assets'

export const generateKoiConfig = () =>
  fs.promises.writeFile(dir('buildPortable', 'koi.yml'), koiConfig)

export const generateKoiVersionInfo = () =>
  fs.promises.writeFile(dir('src', 'versioninfo.json'), koiVersionInfo)

export const generateKoiManifest = () =>
  fs.promises.writeFile(dir('src', 'koi.exe.manifest'), koiManifest)

export const generate = parallel(
  generateKoiConfig,
  generateKoiVersionInfo,
  generateKoiManifest,
  i18nGenerate,
  generateAssets
)
