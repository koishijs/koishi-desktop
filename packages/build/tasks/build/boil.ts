import { series } from 'gulp'
import { sleep } from '../../utils/common'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'
import { stopAndClean } from './stop'

const boilIntl = () =>
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

export const boil = series(
  boilIntl,
  // Intended to have daemon started after boil, so stop it.
  stopAndClean,
  sleep
)
