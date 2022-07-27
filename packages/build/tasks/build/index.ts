import { series } from 'gulp'
import { generate } from './generate'

export const build = series(generate)
