import {
  sourceNode,
  sourceYarn,
  versionNode,
  versionYarn,
} from '../../../utils/config'
import { dir } from '../../../utils/path'

export const nameWin = `node-v${versionNode}-win-x64`
export const nameMac = `node-v${versionNode}-darwin-x64`
export const nameLinux = `node-v${versionNode}-linux-x64`

export const srcWin = `${sourceNode}/v${versionNode}/${nameWin}.zip`
export const srcMac = `${sourceNode}/v${versionNode}/${nameMac}.tar.gz`
export const srcLinux = `${sourceNode}/v${versionNode}/${nameLinux}.tar.xz`
export const srcYarn = `${sourceYarn}/${versionYarn}/packages/yarnpkg-cli/bin/yarn.js`

export const destFileWin = 'node.zip'
export const destFileMac = 'node.tar.gz'
export const destFileLinux = 'node.tar.xz'
export const destFileYarn = 'yarn.cjs'

export const extractCachePath = dir('buildCache', 'node-extract')
