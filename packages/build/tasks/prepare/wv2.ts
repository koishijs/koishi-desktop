import { series } from 'gulp'
import { info } from 'gulplog'
import mkdirp from 'mkdirp'
import StreamZip from 'node-stream-zip'
import * as fs from 'node:fs'
import { versionWebView2 } from '../../utils/config'
import { Exceptions } from '../../utils/exceptions'
import { exists } from '../../utils/fs'
import { download } from '../../utils/net'
import { dir } from '../../utils/path'

const src = `https://api.nuget.org/v3-flatcontainer/microsoft.web.webview2/${versionWebView2}/microsoft.web.webview2.${versionWebView2}.nupkg`
const destFile = 'Microsoft.Web.WebView2.nupkg'

export const prepareWebView2Download = async () => {
  info('Checking temporary cache.')
  if (await exists(dir('buildCache', destFile))) return

  info('Now downloading WebView2.')
  await download(src, dir('buildCache'), destFile)
}

export const prepareWebView2Extract = async () => {
  const targetFolder = dir('buildCache', 'WebView2')
  const cachedFile = dir('buildCache', destFile)

  info('Checking destination cache.')
  if (await exists(dir('buildCache', 'WebView2/Microsoft.Web.WebView2.nuspec')))
    return

  if (!(await exists(cachedFile))) throw Exceptions.fileNotFound(cachedFile)

  const zip = new StreamZip.async({ file: cachedFile })
  await zip.extract(null, targetFolder)
  await zip.close()
}

export const prepareWebView2CopyVendor = async () => {
  const cacheDir = dir('buildCache', 'WebView2/build/native')
  const targetDir = dir('buildVendor', 'WebView2')

  if (!(await exists(cacheDir))) throw Exceptions.fileNotFound(cacheDir)

  await mkdirp(targetDir)
  await fs.promises.cp(cacheDir, targetDir, { recursive: true })
}

export const prepareWebView2CopyDll = async () => {
  const cachedFile = dir(
    'buildCache',
    'WebView2/build/native/x64/WebView2Loader.dll'
  )
  const destFile = dir('buildPortable', 'WebView2Loader.dll')

  if (!(await exists(cachedFile))) throw Exceptions.fileNotFound(cachedFile)

  await fs.promises.copyFile(cachedFile, destFile)
}

export const prepareWebView2Copy = series(prepareWebView2CopyVendor)

export const prepareWebView2 =
  process.platform === 'win32'
    ? series(
        prepareWebView2Download,
        prepareWebView2Extract,
        prepareWebView2Copy
      )
    : async () => {
        /* No need to do anything */
      }
