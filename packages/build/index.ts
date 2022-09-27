import { series } from 'gulp'
import {
  build,
  clean,
  compileApp,
  lint,
  pack,
  prepare,
  startApp,
  stop,
} from './tasks'

export * from './tasks'

export const start = series(compileApp, startApp)

export const dev = series(prepare, stop, build, stop, lint)

export const full = series(clean, prepare, stop, build, stop, lint, pack)

export const ciBuild = series(prepare, build, stop, pack)

export const ciLint = series(prepare, build, stop, lint)

export const defaultTask = dev
