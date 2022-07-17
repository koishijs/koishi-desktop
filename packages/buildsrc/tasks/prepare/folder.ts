import { parallel } from 'gulp'
import mkdirp from 'mkdirp'
import { BuildPaths, dir } from '../../utils/path'

const prepareFolderTasks = (['buildPortableData', 'dist'] as BuildPaths[]).map(
  (x) => () => mkdirp(dir(x)) as Promise<void>
)

export const prepareFolder = parallel(...prepareFolderTasks)
