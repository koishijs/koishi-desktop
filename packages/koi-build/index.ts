export * from './build'
export * from './common'
export * from './pack'
export * from './prepareNode'
export * from './menifest'

import { series } from 'gulp'
import { build, buildExe, run } from './build'
import { clean, mkdir } from './common'
import { pack } from './pack'
import { prepareNode } from './prepareNode'

export const dev = series(mkdir, prepareNode, build)

export const test = series(mkdir, prepareNode, buildExe, run)

export const all = series(clean, dev, pack)

export default dev
