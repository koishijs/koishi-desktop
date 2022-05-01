import axios from 'axios'
import del from 'del'
import * as fs from 'fs'
import { error, info } from 'gulplog'
import * as lzma from 'lzma-native'
import StreamZip from 'node-stream-zip'
import stream from 'stream'
import * as tar from 'tar'
import { promisify } from 'util'
import { nodeVersion } from './config'
import { resolve } from './path'
import { exists, notEmpty } from './utils'

const nodeFolderWin = `node-v${nodeVersion}-win-${process.arch}`
const srcPathWin = `https://nodejs.org/dist/v${nodeVersion}/${nodeFolderWin}.zip`
const destPathWin = resolve('node.zip', 'buildTemp')
const nodeFolderMac = `node-v${nodeVersion}-darwin-${process.arch}`
const srcPathMac = `https://nodejs.org/dist/v${nodeVersion}/${nodeFolderMac}.tar.gz`
const destPathMac = resolve('node.tar.gz', 'buildTemp')
const nodeFolderLinux = `node-v${nodeVersion}-linux-${process.arch}`
const srcPathLinux = `https://nodejs.org/dist/v${nodeVersion}/${nodeFolderLinux}.tar.xz`
const destPathLinux = resolve('node.tar.xz', 'buildTemp')

const buildDownloadNode =
  (srcPath: string, destPath: string): (() => Promise<void>) =>
  async () => {
    info('Checking temporary cache.')
    if (await exists(destPath)) {
      info('Node.js exists. Skipping download.')
      info("If you want to re-download Node.js, use 'gulp clean'.")
      return
    }

    info('Now downloading.')
    const res = await axios.get(srcPath, { responseType: 'stream' })
    const writeStream = fs.createWriteStream(destPath)
    await promisify(stream.finished)(
      (res.data as stream.Readable).pipe(writeStream)
    )
  }

async function extractNode(destPath: string) {
  info('Checking temporary cache.')
  if (await notEmpty(resolve('node', 'distData'))) {
    info('Node.js exists. Skipping extract.')
    return
  }

  if (!(await exists(destPath))) {
    const err = "Node.js dist not found. Try 'gulp clean && gulp prepareNode'."
    error(err)
    throw new Error(err)
  }

  info('Now extracting.')
  switch (process.platform) {
    case 'win32': {
      const zip = new StreamZip.async({ file: destPath })
      await zip.extract(nodeFolderWin, resolve('node', 'distData'))
      await zip.close()
      break
    }
    case 'darwin':
      await promisify(stream.finished)(
        fs
          .createReadStream(destPath)
          .pipe(tar.extract({ cwd: resolve('node', 'distData'), strip: 1 }))
      )
      break
    case 'linux':
      await promisify(stream.finished)(
        fs
          .createReadStream(destPath)
          .pipe(lzma.createDecompressor())
          .pipe(tar.extract({ cwd: resolve('node', 'distData'), strip: 1 }))
      )
      break
    default: {
      const err = `Platform ${process.platform} not supported yet`
      error(err)
      throw new Error(err)
    }
  }
}

async function removeNpmWindows() {
  await del(resolve('node/CHANGELOG.md', 'distData'))
  await del(resolve('node/README.md', 'distData'))
  await del(resolve('node/node_modules', 'distData'))
  await del(resolve('node/node_etw_provider.man', 'distData'))
  await del(resolve('node/install_tools.bat', 'distData'))
  await del(resolve('node/nodevars.bat', 'distData'))
  await del(resolve('node/corepack', 'distData'))
  await del(resolve('node/corepack.cmd', 'distData'))
  await del(resolve('node/npm', 'distData'))
  await del(resolve('node/npm.cmd', 'distData'))
  await del(resolve('node/npx', 'distData'))
  await del(resolve('node/npx.cmd', 'distData'))
}

async function removeNpmUnix() {
  await del(resolve('node/CHANGELOG.md', 'distData'))
  await del(resolve('node/README.md', 'distData'))
  await del(resolve('node/bin/corepack', 'distData'))
  await del(resolve('node/bin/npm', 'distData'))
  await del(resolve('node/bin/npx', 'distData'))
  await del(resolve('node/include', 'distData'))
  await del(resolve('node/lib', 'distData'))
  await del(resolve('node/share', 'distData'))
}

export async function prepareNode(): Promise<void> {
  info(`Downloading Node.js for ${process.platform} on ${process.arch}.`)

  switch (process.platform) {
    case 'win32':
      await buildDownloadNode(srcPathWin, destPathWin)()
      await extractNode(destPathWin)
      await removeNpmWindows()
      break
    case 'darwin':
      await buildDownloadNode(srcPathMac, destPathMac)()
      await extractNode(destPathMac)
      await removeNpmUnix()
      break
    case 'linux':
      await buildDownloadNode(srcPathLinux, destPathLinux)()
      await extractNode(destPathLinux)
      await removeNpmUnix()
      break
    default: {
      const err = `Platform ${process.platform} not supported yet`
      error(err)
      throw new Error(err)
    }
  }
}
