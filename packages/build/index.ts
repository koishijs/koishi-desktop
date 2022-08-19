import { series } from 'gulp'
import { build, clean, lint, pack, prepare } from './tasks'

export const dev = series(prepare, build, lint)

export const full = series(clean, prepare, build, lint, pack)

export const defaultTask = dev
