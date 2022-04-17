import axios from 'axios'
import * as fs from 'fs'
// @ts-ignore
import { error, info } from 'gulplog'
import StreamZip from 'node-stream-zip'
import stream from 'stream'
import { promisify } from 'util'
import { nodeVersion } from './config'
import { resolve } from './path'
import { exists } from './utils'

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
    ;(res.data as stream.Readable).pipe(writeStream)
    await promisify(stream.finished)(writeStream)
  }

async function extractNodeWin() {
  info('Checking temporary cache.')
  if (!(await exists(destPathWin))) {
    const err = "Node.js dist not found. Try 'gulp clean && gulp prepareNode'."
    error(err)
    throw new Error(err)
  }

  info('Now extracting.')
  const zip = new StreamZip.async({ file: destPathWin })
  await zip.extract(nodeFolderWin, resolve('node', 'distData'))
  await zip.close()
}

export async function prepareNode(): Promise<void> {
  info(`Downloading Node.js for ${process.platform} on ${process.arch}.`)

  switch (process.platform) {
    case 'win32':
      await buildDownloadNode(srcPathWin, destPathWin)()
      await extractNodeWin()
      break
    case 'darwin':
      await buildDownloadNode(srcPathMac, destPathMac)()
      break
    case 'linux':
      buildDownloadNode(srcPathLinux, destPathLinux)()
      break
    default:
      const err = `Platform ${process.platform} not supported yet`
      error(err)
      throw new Error(err)
  }
}
