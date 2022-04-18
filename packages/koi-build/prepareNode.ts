import axios from 'axios'
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
      await zip.extract(nodeFolderWin, resolve('data/node', 'dist'))
      await zip.close()
      break
    }
    case 'darwin':
      await promisify(stream.finished)(
        fs
          .createReadStream(destPath)
          .pipe(tar.extract({ cwd: resolve('data/node', 'dist'), strip: 1 }))
      )
      break
    case 'linux':
      await promisify(stream.finished)(
        fs
          .createReadStream(destPath)
          .pipe(lzma.createDecompressor())
          .pipe(tar.extract({ cwd: resolve('data/node', 'dist'), strip: 1 }))
      )
      break
    default: {
      const err = `Platform ${process.platform} not supported yet`
      error(err)
      throw new Error(err)
    }
  }
}

export async function prepareNode(): Promise<void> {
  info(`Downloading Node.js for ${process.platform} on ${process.arch}.`)

  switch (process.platform) {
    case 'win32':
      await buildDownloadNode(srcPathWin, destPathWin)()
      await extractNode(destPathWin)
      break
    case 'darwin':
      await buildDownloadNode(srcPathMac, destPathMac)()
      await extractNode(destPathMac)
      break
    case 'linux':
      await buildDownloadNode(srcPathLinux, destPathLinux)()
      await extractNode(destPathLinux)
      break
    default: {
      const err = `Platform ${process.platform} not supported yet`
      error(err)
      throw new Error(err)
    }
  }
}
