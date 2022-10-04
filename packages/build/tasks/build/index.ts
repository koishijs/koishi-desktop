import { series } from 'gulp'
import { boil } from './boil'
import { cleanTemp } from './clean'
import { compile } from './compile'
import { generate } from './generate'

export * from './assets'
export * from './boil'
export * from './clean'
export * from './compile'
export * from './generate'

export const build = series(generate, compile, boil, cleanTemp)
