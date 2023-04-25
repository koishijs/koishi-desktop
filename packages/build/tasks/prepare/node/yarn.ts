import { series } from 'gulp'
import { info } from 'gulplog'
import * as fs from 'node:fs'
import { Exceptions } from '../../../utils/exceptions'
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
  const cachedFile = dir('buildCache', destFileYarn)
  const destFile = dir('buildPortableBin', destFileYarn)
  info('Checking destination cache.')
  if (await exists(destFile)) return

  if (!(await exists(cachedFile))) {
    throw Exceptions.fileNotFound(cachedFile)
  }

  await fs.promises.copyFile(cachedFile, destFile)
}

export const prepareNodeYarn = series(
  prepareNodeYarnDownload,
  prepareNodeYarnCopy
)
