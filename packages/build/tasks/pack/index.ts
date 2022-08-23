import { parallel } from 'gulp'
import { packPortable } from './portable'

export * from './portable'

export const pack = parallel(packPortable)
