import { series } from 'gulp'
import { boil } from './boil'
import { compile } from './compile'
import { generate } from './generate'

export * from './boil'
export * from './compile'
export * from './generate'

export const build = series(generate, compile, boil)
