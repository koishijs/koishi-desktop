import { parallel, series } from 'gulp'
import { prepareAppimagetool } from './appimagetool'
import { prepareBoilerplate } from './boilerplate'
import { prepareEdge } from './edge'
import { prepareFolder } from './folder'
import { prepareGoMod } from './gomod'
import { prepareNode } from './node'
import { prepareTools } from './tools'
import { prepareWix } from './wix'

export * from './boilerplate'
export * from './folder'
export * from './gomod'
export * from './node'
export * from './tools'
export * from './wix'
export * from './appimagetool'
export * from './edge'

export const prepare = parallel(
  prepareTools,
  prepareGoMod,
  series(prepareFolder, parallel(prepareNode, prepareBoilerplate)),
  prepareEdge,
  prepareAppimagetool,
  prepareWix
)
