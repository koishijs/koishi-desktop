import del from 'del'
import { info } from 'gulplog'
import * as lzma from 'lzma-native'
import StreamZip from 'node-stream-zip'
import * as fs from 'node:fs'
import stream from 'node:stream'
import { promisify } from 'node:util'
import * as tar from 'tar'
import { koishiManifest, koishiVersionStrings } from '../../../utils/config'
import { Exceptions } from '../../../utils/exceptions'
import { exists } from '../../../utils/fs'
import { dir } from '../../../utils/path'
import { exec } from '../../../utils/spawn'
import { destFileLinux, destFileMac, destFileWin, nameWin } from './path'

export const prepareNodeExtractWin = async () => {
  const nodeFolder = dir('buildPortableData', 'node')
  const cachedFile = dir('buildCache', destFileWin)

  info('Checking destination cache.')
  if (await exists(dir('buildPortableData', 'node/node.exe'))) return

  if (!(await exists(cachedFile))) {
    throw Exceptions.fileNotFound(cachedFile)
  }

  const zip = new StreamZip.async({ file: cachedFile })
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

export const prepareNodeRcedit = async () => {
  const koishiManifestPath = dir('buildCache', 'koishi.exe.manifest')

  await fs.promises.writeFile(koishiManifestPath, koishiManifest)

  const args = [
    dir('buildPortableData', 'node/koishi.exe'),
    '--set-icon',
    dir('src', 'resources/koi.ico'),
    '--application-manifest',
    koishiManifestPath,
  ]

  ;(
    Object.keys(koishiVersionStrings) as (keyof typeof koishiVersionStrings)[]
  ).forEach((x) => {
    args.push('--set-version-string', x, koishiVersionStrings[x])
  })

  info(args)

  await exec('rcedit.exe', args, dir('buildCache'))
}

export const prepareNodeExtractMac = async () => {
  const nodeFolder = dir('buildPortableData', 'node')
  const cachedFile = dir('buildCache', destFileMac)

  info('Checking destination cache.')
  if (await exists(dir('buildPortableData', 'node/bin/node'))) return

  if (!(await exists(cachedFile))) {
    throw Exceptions.fileNotFound(cachedFile)
  }

  await promisify(stream.finished)(
    fs
      .createReadStream(cachedFile)
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
  const cachedFile = dir('buildCache', destFileLinux)

  info('Checking destination cache.')
  if (await exists(dir('buildPortableData', 'node/bin/node'))) return

  if (!(await exists(cachedFile))) {
    throw Exceptions.fileNotFound(cachedFile)
  }

  await promisify(stream.finished)(
    fs
      .createReadStream(cachedFile)
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
