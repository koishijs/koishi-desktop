import { series } from 'gulp'
import { boil } from './boil'
import { cleanTemp } from './clean'
import { compile } from './compile'
import { generateAfter, generateBefore } from './generate'

export * from './boil'
export * from './clean'
export * from './compile'
export * from './generate'
export * from './assets'

export const build = series(
  generateBefore,
  compile,
  boil,
  generateAfter,
  cleanTemp
)
