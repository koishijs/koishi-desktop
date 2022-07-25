import { info } from 'gulplog'
import { join } from 'node:path'
import { Exceptions } from '../../../utils/exceptions'
import { exists } from '../../../utils/fs'
import { download } from '../../../utils/net'
import { dir } from '../../../utils/path'
import {
  destFileLinux,
  destFileMac,
  destFileWin,
  srcLinux,
  srcMac,
  srcWin,
} from './path'

const buildPrepareNodeDownloadImpl =
  (srcPath: string, dest: string, filename: string): (() => Promise<void>) =>
  async () => {
    info('Checking temporary cache.')
    if (await exists(join(dest, filename))) return

    info('Now downloading Node.js.')
    await download(srcPath, dest, filename)
  }

export const prepareNodeDownloadWin = buildPrepareNodeDownloadImpl(
  srcWin,
  dir('buildCache'),
  destFileWin
)

export const prepareNodeDownloadMac = buildPrepareNodeDownloadImpl(
  srcMac,
  dir('buildCache'),
  destFileMac
)

export const prepareNodeDownloadLinux = buildPrepareNodeDownloadImpl(
  srcLinux,
  dir('buildCache'),
  destFileLinux
)

const buildPrepareNodeDownload = (): (() => Promise<void>) => {
  switch (process.platform) {
    case 'win32':
      return prepareNodeDownloadWin
    case 'darwin':
      return prepareNodeDownloadMac
    case 'linux':
      return prepareNodeDownloadLinux
    default:
      throw Exceptions.platformNotSupportedException
  }
}

export const prepareNodeDownload = buildPrepareNodeDownload()
