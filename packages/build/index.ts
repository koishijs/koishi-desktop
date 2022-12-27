import { series } from 'gulp'
import {
  build,
  clean,
  compileApp,
  dirty,
  lint,
  pack,
  prepare,
  startApp,
  stopAndClean,
} from './tasks'

export * from './tasks'

export const start = series(compileApp, startApp)

export const dev = series(prepare, stopAndClean, build, stopAndClean, lint)

export const full = series(
  clean,
  prepare,
  stopAndClean,
  build,
  stopAndClean,
  lint,
  pack
)

export const ciBuild = series(prepare, build, stopAndClean, dirty, pack)

export const ciLint = series(prepare, build, stopAndClean, lint)

export const defaultTask = dev
