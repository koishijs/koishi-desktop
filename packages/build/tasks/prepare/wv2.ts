import { info } from 'gulplog'
import { exists } from '../../utils/fs'
import { download } from '../../utils/net'
import { dir } from '../../utils/path'

const src = `https://go.microsoft.com/fwlink/p/?LinkId=2124703`
const destFile = 'Webview2Setup.exe'

export const prepareWebView2Download = async () => {
  info('Checking temporary cache.')
  if (await exists(dir('buildCache', destFile))) return

  info('Now downloading WebView2.')
  await download(src, dir('buildCache'), destFile)
}

export const prepareWebView2 =
  process.platform === 'win32'
    ? prepareWebView2Download
    : async () => {
        /* No need to do anything */
      }
