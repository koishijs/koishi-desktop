import { dir } from '../utils/path'
import { exec } from '../utils/spawn'

export const start = () =>
  exec(process.platform === 'win32' ? 'koi' : './koi', [], dir('buildPortable'))
