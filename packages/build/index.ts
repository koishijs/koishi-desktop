import { series } from 'gulp'
import { build, clean, lint, pack, prepare } from './tasks'

export const dev = series(prepare, lint, build)

export const full = series(clean, prepare, lint, build, pack)

export const defaultTask = dev
