import { parallel } from 'gulp'
import mkdirp from 'mkdirp'
import { dir } from '../../utils/path'

const prepareFolderTasks = [
  dir('buildCache'),
  dir('buildPortableData'),
  dir('dist'),
  dir('buildPortableData', 'node'),
].map((x) => () => mkdirp(x) as Promise<void>)

export const prepareFolder = parallel(...prepareFolderTasks)
