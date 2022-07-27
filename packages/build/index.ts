import { series } from 'gulp'
import { build, clean, pack, prepare } from './tasks'

export const dev = series(prepare, build)

export const full = series(clean, prepare, build, pack)

export const defaultTask = dev
