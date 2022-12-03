import { info } from 'gulplog'
import { exists } from '../../utils/fs'
import { download } from '../../utils/net'
import { dir } from '../../utils/path'

const files = [
  {
    src: 'https://go.microsoft.com/fwlink/?linkid=2109047&Channel=Stable&language=en&Consent=1&brand=M100',
    destFile: 'MicrosoftEdgeSetup.exe',
  },
  {
    src: 'https://go.microsoft.com/fwlink/p/?LinkId=2124703',
    destFile: 'Webview2Setup.exe',
  },
]

export const prepareEdge = async () => {
  if (process.platform !== 'win32') return

  for (const file of files) {
    info('Checking temporary cache.')
    if (await exists(dir('buildCache', file.destFile))) return

    info(`Now downloading ${file.destFile}.`)
    await download(file.src, dir('buildCache'), file.destFile)
  }
}
