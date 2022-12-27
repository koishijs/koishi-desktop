import { series } from 'gulp'
import { boil } from './boil'
import { compile } from './compile'
import { generate } from './generate'
import { patch } from './patch'
import { stopAndClean } from './stop'

export * from './assets'
export * from './boil'
export * from './compile'
export * from './compileShell'
export * from './generate'
export * from './patch'
export * from './stop'
export * from './userscript'

export const build = series(generate, compile, patch, boil, stopAndClean)
