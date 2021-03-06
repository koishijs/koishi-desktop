import del from 'del'
import { info } from 'gulplog'
import * as lzma from 'lzma-native'
import StreamZip from 'node-stream-zip'
import * as fs from 'node:fs'
import stream from 'node:stream'
import { promisify } from 'node:util'
import * as tar from 'tar'
import { Exceptions } from '../../../utils/exceptions'
import { exists } from '../../../utils/fs'
import { dir } from '../../../utils/path'
import { destFileLinux, destFileMac, destFileWin, nameWin } from './path'

export const prepareNodeExtractWin = async () => {
  const nodeFolder = dir('buildPortableData', 'node')

  info('Checking temporary cache.')
  if (await exists(nodeFolder)) return

  if (!(await exists(dir('buildCache', destFileWin)))) {
    throw Exceptions.fileNotFound(destFileWin)
  }

  const zip = new StreamZip.async({ file: destFileWin })
  await zip.extract(nameWin, nodeFolder)
  await zip.close()

  await del(dir('buildPortableData', 'node/CHANGELOG.md'))
  await del(dir('buildPortableData', 'node/README.md'))
  await del(dir('buildPortableData', 'node/node_modules'))
  await del(dir('buildPortableData', 'node/node_etw_provider.man'))
  await del(dir('buildPortableData', 'node/install_tools.bat'))
  await del(dir('buildPortableData', 'node/nodevars.bat'))
  await del(dir('buildPortableData', 'node/corepack'))
  await del(dir('buildPortableData', 'node/corepack.cmd'))
  await del(dir('buildPortableData', 'node/npm'))
  await del(dir('buildPortableData', 'node/npm.cmd'))
  await del(dir('buildPortableData', 'node/npx'))
  await del(dir('buildPortableData', 'node/npx.cmd'))
}

export const prepareNodeExtractMac = async () => {
  const nodeFolder = dir('buildPortableData', 'node')

  info('Checking temporary cache.')
  if (await exists(nodeFolder)) return

  if (!(await exists(dir('buildCache', destFileMac)))) {
    throw Exceptions.fileNotFound(destFileMac)
  }

  await promisify(stream.finished)(
    fs
      .createReadStream(destFileMac)
      .pipe(tar.extract({ cwd: nodeFolder, strip: 1 }))
  )

  await del(dir('buildPortableData', 'node/CHANGELOG.md'))
  await del(dir('buildPortableData', 'node/README.md'))
  await del(dir('buildPortableData', 'node/bin/corepack'))
  await del(dir('buildPortableData', 'node/bin/npm'))
  await del(dir('buildPortableData', 'node/bin/npx'))
  await del(dir('buildPortableData', 'node/include'))
  await del(dir('buildPortableData', 'node/lib'))
  await del(dir('buildPortableData', 'node/share'))
}

export const prepareNodeExtractLinux = async () => {
  const nodeFolder = dir('buildPortableData', 'node')

  info('Checking temporary cache.')
  if (await exists(nodeFolder)) return

  if (!(await exists(dir('buildCache', destFileLinux)))) {
    throw Exceptions.fileNotFound(destFileLinux)
  }

  await promisify(stream.finished)(
    fs
      .createReadStream(destFileLinux)
      .pipe(lzma.createDecompressor())
      .pipe(tar.extract({ cwd: nodeFolder, strip: 1 }))
  )

  await del(dir('buildPortableData', 'node/CHANGELOG.md'))
  await del(dir('buildPortableData', 'node/README.md'))
  await del(dir('buildPortableData', 'node/bin/corepack'))
  await del(dir('buildPortableData', 'node/bin/npm'))
  await del(dir('buildPortableData', 'node/bin/npx'))
  await del(dir('buildPortableData', 'node/include'))
  await del(dir('buildPortableData', 'node/lib'))
  await del(dir('buildPortableData', 'node/share'))
}

const buildPrepareNodeExtract = (): (() => Promise<void>) => {
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
