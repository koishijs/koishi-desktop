import { series } from 'gulp'
import { build, clean, lint, pack, prepare, stop } from './tasks'

export * from './tasks'

export const dev = series(prepare, stop, build, stop, lint)

export const full = series(clean, prepare, stop, build, stop, lint, pack)

export const defaultTask = dev
