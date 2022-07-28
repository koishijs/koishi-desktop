import { parallel, series } from 'gulp'
import { prepareBoilerplate } from './boilerplate'
import { prepareFolder } from './folder'
import { prepareGoMod } from './gomod'
import { prepareNode } from './node'
import { prepareTools } from './tools'

export * from './boilerplate'
export * from './folder'
export * from './gomod'
export * from './node'
export * from './tools'

export const prepare = parallel(
  prepareTools,
  prepareGoMod,
  series(prepareFolder, parallel(prepareNode, prepareBoilerplate))
)
