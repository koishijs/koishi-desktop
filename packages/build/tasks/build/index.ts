import { series } from 'gulp'
import { compile } from './compile'
import { generate } from './generate'

export * from './generate'
export * from './compile'

export const build = series(generate, compile)
