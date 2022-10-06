import { parallel, series } from 'gulp'
import mkdirp from 'mkdirp'
import { promises as fs } from 'node:fs'
import { msiWxs } from '../../templates'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'

export const packMsiMkdir = () => mkdirp(dir('buildMsi'))

export const packMsiIndex = () =>
  fs.writeFile(dir('buildMsi', 'index.wxs'), msiWxs)

export const packMsiFiles = async () => {
  const dirSource = dir('buildMsi', 'SourceDir/')
  await mkdirp(dirSource)
  await fs.cp(dir('buildUnfoldBinary'), dirSource, { recursive: true })
}

export const packMsiCandle = () =>
  exec(
    dir('buildVendor', 'wix/candle.exe'),
    ['-nologo', 'index.wxs'],
    dir('buildMsi')
  )

export const packMsiLight = () =>
  exec(
    dir('buildVendor', 'wix/light.exe'),
    [
      '-nologo',
      '-sice:ICE61',
      '-sice:ICE69',
      '-spdb',
      '-out',
      dir('dist', 'koishi.msi'),
      '-ext',
      'WixUIExtension',
      'index.wixobj',
    ],
    dir('buildMsi')
  )

export const packMsi = series(
  packMsiMkdir,
  parallel(packMsiIndex, packMsiFiles),
  packMsiCandle,
  packMsiLight
)
