import { parallel, series } from 'gulp'
import { prepareFolder } from './folder'
import { prepareNode } from './node'
import { prepareTools } from './tools'

export * from './folder'
export * from './node'
export * from './tools'

export const prepare = parallel(
  prepareTools,
  series(prepareFolder, prepareNode)
)
