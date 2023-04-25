import { parallel } from 'gulp'
import { mkdir } from 'node:fs/promises'
import { dir } from '../../utils/path'
import { prepareAppimagetool } from './appimagetool'
import { prepareBoilerplate } from './boilerplate'
import { prepareEdge } from './edge'
import { prepareGoMod } from './gomod'
import { prepareNode } from './node'
import { prepareTools } from './tools'
import { prepareWix } from './wix'

export * from './appimagetool'
export * from './boilerplate'
export * from './edge'
export * from './gomod'
export * from './node'
export * from './tools'
export * from './wix'

export const prepare = parallel(
  () =>
    mkdir(dir('buildAssets'), {
      recursive: true,
    }),
  () =>
    mkdir(dir('buildVendor'), {
      recursive: true,
    }),
  () =>
    mkdir(dir('buildResources'), {
      recursive: true,
    }),
  () =>
    mkdir(dir('dist'), {
      recursive: true,
    }),

  prepareTools,
  prepareGoMod,
  prepareNode,
  prepareBoilerplate,
  prepareEdge,
  prepareAppimagetool,
  prepareWix
)
