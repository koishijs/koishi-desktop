export * from './build'
export * from './common'
export * from './menifest'
export * from './pack'
export * from './prepareNode'

import { series } from 'gulp'
import { build } from './build'
import { clean, mkdir } from './common'
import { pack } from './pack'
import { prepareNode } from './prepareNode'

export const dev = series(mkdir, prepareNode, build)

export const all = series(clean, dev, pack)

export default dev
