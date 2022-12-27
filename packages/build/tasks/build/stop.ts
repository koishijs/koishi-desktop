import del from 'del'
import { series } from 'gulp'
import { info } from 'gulplog'
import { sleep } from '../../utils/common'
import { exists } from '../../utils/fs'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'

export const cleanTemp = async () => {
  const dirs = [
    dir('buildPortableData', 'home'),
    dir('buildPortableData', 'lock'),
    dir('buildPortableData', 'logs'),
    dir('buildPortableData', 'tmp'),
  ]
  for (const d of dirs) await del(d)
}

export const stop = async () => {
  if (
    !(await exists(
      dir('buildPortable', process.platform === 'win32' ? 'koi.exe' : 'koi')
    ))
  )
    return

  await exec(
    process.platform === 'win32' ? 'koi' : './koi',
    ['daemon', 'stop'],
    dir('buildPortable')
  )

  await sleep()

  info('Remaining process stopped.')
}

export const stopAndClean = series(stop, cleanTemp)

export const kill = async () => {
  if (
    !(await exists(
      dir('buildPortable', process.platform === 'win32' ? 'koi.exe' : 'koi')
    ))
  )
    return

  await exec(
    process.platform === 'win32' ? 'koi' : './koi',
    ['daemon', 'kill'],
    dir('buildPortable')
  )

  await sleep()

  info('Remaining process killed.')
}
