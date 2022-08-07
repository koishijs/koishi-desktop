import { series } from 'gulp'
import { info } from 'gulplog'
import * as fs from 'node:fs'
import { exists } from '../../../utils/fs'
import { download } from '../../../utils/net'
import { dir } from '../../../utils/path'
import { destFileYarn, srcYarn } from './path'

export const prepareNodeYarnDownload = async () => {
  info('Checking temporary cache.')
  if (await exists(dir('buildCache', destFileYarn))) return

  info('Now downloading Yarn.')
  await download(srcYarn, dir('buildCache'), destFileYarn)
}

export const prepareNodeYarnCopy = async () => {
  const destFile = dir('buildCache', destFileYarn)
  info('Checking destination cache.')
  if (await exists(destFile)) return

  await fs.promises.copyFile(
    destFile,
    dir(
      'buildPortableData',
      (process.platform === 'win32' ? 'node/' : 'node/bin/') + destFileYarn
    )
  )
}

export const prepareNodeYarn = series(
  prepareNodeYarnDownload,
  prepareNodeYarnCopy
)
