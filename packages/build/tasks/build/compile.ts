import { series } from 'gulp'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'

export const compileWire = () => exec('wire', [], dir('src'))

export const compileVersioninfo = () => exec('goversioninfo', [], dir('src'))

export const compile = series(compileWire, compileVersioninfo)
