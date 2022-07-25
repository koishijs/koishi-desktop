import { series } from 'gulp'
import { prepareFolder } from './folder'
import { prepareNode } from './node'

export * from './folder'
export * from './node'

export const prepare = series(prepareFolder, prepareNode)
