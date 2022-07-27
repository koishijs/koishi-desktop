import {
  sourceNode,
  sourceYarn,
  versionNode,
  versionYarn,
} from '../../../utils/config'

export const nameWin = `node-v${versionNode}-win-${process.arch}`
export const nameMac = `node-v${versionNode}-darwin-${process.arch}`
export const nameLinux = `node-v${versionNode}-linux-${process.arch}`

export const srcWin = `${sourceNode}/v${versionNode}/${nameWin}.zip`
export const srcMac = `${sourceNode}/v${versionNode}/${nameMac}.tar.gz`
export const srcLinux = `${sourceNode}/v${versionNode}/${nameLinux}.tar.xz`
export const srcYarn = `${sourceYarn}/${versionYarn}/packages/yarnpkg-cli/bin/yarn.js`

export const destFileWin = 'node.zip'
export const destFileMac = 'node.tar.gz'
export const destFileLinux = 'node.tar.xz'
export const destFileYarn = 'yarn.cjs'
