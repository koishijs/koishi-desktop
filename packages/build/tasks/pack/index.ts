import { series } from 'gulp'
import { packPortable } from './portable'

export * from './portable'

export const pack = series(packPortable)
