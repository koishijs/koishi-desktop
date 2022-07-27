import { parallel } from 'gulp'
import * as fs from 'fs'
import { dir } from '../../utils/path'
import { koiConfig, koiManifest, koiVersionInfo } from '../../utils/config'
import mkdirp from 'mkdirp'

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

export const generate = parallel(
  generateKoiConfig,
  generateKoiVersionInfo,
  generateKoiManifest
)
