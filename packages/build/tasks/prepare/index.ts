import { parallel, series } from 'gulp'
import { prepareAppimagetool } from './appimagetool'
import { prepareBoilerplate } from './boilerplate'
import { prepareFolder } from './folder'
import { prepareGoMod } from './gomod'
import { prepareNode } from './node'
import { prepareTools } from './tools'
import { prepareWix } from './wix'
import { prepareWebView2 } from './wv2'

export * from './boilerplate'
export * from './folder'
export * from './gomod'
export * from './node'
export * from './wv2'
export * from './tools'
export * from './wix'
export * from './appimagetool'

export const prepare = parallel(
  prepareTools,
  prepareGoMod,
  series(prepareFolder, parallel(prepareNode, prepareBoilerplate)),
  prepareWebView2,
  prepareAppimagetool,
  prepareWix
)
