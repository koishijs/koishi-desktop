import { TaskFunction } from 'gulp'
import { info } from 'gulplog'
import * as lzma from 'lzma-native'
import StreamZip from 'node-stream-zip'
import { createReadStream } from 'node:fs'
import { copyFile } from 'node:fs/promises'
import { join } from 'node:path'
import stream from 'node:stream'
import { promisify } from 'node:util'
import * as tar from 'tar'
import { Exceptions } from '../../../utils/exceptions'
import { exists } from '../../../utils/fs'
import { dir } from '../../../utils/path'
import {
  destFileLinux,
  destFileMac,
  destFileWin,
  extractCachePath,
  nameWin,
} from './path'

export const prepareNodeExtractWin = async () => {
  const cachedFile = dir('buildCache', destFileWin)

  info('Checking destination cache.')
  if (await exists(dir('buildPortableBin', 'koishi.exe'))) return

  if (!(await exists(cachedFile))) {
    throw Exceptions.fileNotFound(cachedFile)
  }

  const zip = new StreamZip.async({ file: cachedFile })
  await zip.extract(nameWin, extractCachePath)
  await zip.close()

  await copyFile(
    join(extractCachePath, 'node.exe'),
    dir('buildPortableBin', 'koishi.exe')
  )
}

export const prepareNodeExtractMac = async () => {
  const cachedFile = dir('buildCache', destFileMac)

  info('Checking destination cache.')
  if (await exists(dir('buildPortableBin', 'koishi'))) return

  if (!(await exists(cachedFile))) {
    throw Exceptions.fileNotFound(cachedFile)
  }

  await promisify(stream.finished)(
    createReadStream(cachedFile).pipe(
      tar.extract({ cwd: extractCachePath, strip: 1 })
    )
  )

  await copyFile(
    join(extractCachePath, 'bin/node'),
    dir('buildPortableBin', 'koishi')
  )
}

export const prepareNodeExtractLinux = async () => {
  const cachedFile = dir('buildCache', destFileLinux)

  info('Checking destination cache.')
  if (await exists(dir('buildPortableBin', 'koishi'))) return

  if (!(await exists(cachedFile))) {
    throw Exceptions.fileNotFound(cachedFile)
  }

  await promisify(stream.finished)(
    createReadStream(cachedFile)
      .pipe(lzma.createDecompressor())
      .pipe(tar.extract({ cwd: extractCachePath, strip: 1 }))
  )

  await copyFile(
    join(extractCachePath, 'bin/node'),
    dir('buildPortableBin', 'koishi')
  )
}

const buildPrepareNodeExtract = (): TaskFunction => {
  switch (process.platform) {
    case 'win32':
      return prepareNodeExtractWin
    case 'darwin':
      return prepareNodeExtractMac
    case 'linux':
      return prepareNodeExtractLinux
    default:
      throw Exceptions.platformNotSupported()
  }
}

export const prepareNodeExtract = buildPrepareNodeExtract()
