import { series } from 'gulp'
import { build, clean, kill, lint, pack, prepare } from './tasks'

export * from './tasks'

export const dev = series(prepare, kill, build, kill, lint)

export const full = series(clean, prepare, kill, build, kill, lint, pack)

export const defaultTask = dev
