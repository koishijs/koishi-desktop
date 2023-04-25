import { dir } from '../utils/path'
import { exec2 } from '../utils/spawn'

export const startApp = () =>
  exec2(
    process.platform === 'win32' ? 'koi' : './koi',
    [],
    dir('buildPortable')
  )
