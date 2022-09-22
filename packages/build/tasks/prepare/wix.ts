import { series } from 'gulp'
import { info } from 'gulplog'
import StreamZip from 'node-stream-zip'
import { sourceGitHub, versionToolsWix } from '../../utils/config'
import { Exceptions } from '../../utils/exceptions'
import { exists } from '../../utils/fs'
import { download } from '../../utils/net'
import { dir } from '../../utils/path'

const src = `${sourceGitHub}/wixtoolset/wix3/releases/download/wix${versionToolsWix}2rtm/wix${versionToolsWix}-binaries.zip`
const destFile = 'wix.zip'

export const prepareWixDownload = async () => {
  info('Checking temporary cache.')
  if (await exists(dir('buildCache', destFile))) return

  info('Now downloading Wix Toolset.')
  await download(src, dir('buildCache'), destFile)
}

export const prepareWixExtract = async () => {
  const targetFolder = dir('buildVendor', 'wix')
  const cachedFile = dir('buildCache', destFile)

  info('Checking destination cache.')
  if (await exists(dir('buildVendor', 'wix/candle.exe'))) return

  if (!(await exists(cachedFile))) throw Exceptions.fileNotFound(cachedFile)

  const zip = new StreamZip.async({ file: cachedFile })
  await zip.extract(null, targetFolder)
  await zip.close()
}

export const prepareWix =
  process.platform === 'win32'
    ? series(prepareWixDownload, prepareWixExtract)
    : async () => {
        // Ignore
      }
