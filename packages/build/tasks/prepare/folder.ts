import { parallel } from 'gulp'
import mkdirp from 'mkdirp'
import { dir } from '../../utils/path'

const prepareFolderTasks = [
  dir('buildCache'),
  dir('buildVendor'),
  dir('buildPortableData'),
  dir('dist'),
  dir('buildPortableData', process.platform === 'win32' ? 'node' : 'node/bin'),
].map((x) => () => mkdirp(x) as Promise<void>)

export const prepareFolder = parallel(...prepareFolderTasks)
