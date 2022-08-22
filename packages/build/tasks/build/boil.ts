import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'

export const boil = () =>
  exec(
    process.platform === 'win32' ? 'koi' : './koi',
    [
      'import',
      '--name',
      'default',
      '--force',
      dir('buildCache', 'boilerplate.zip'),
    ],
    dir('buildPortable')
  )
