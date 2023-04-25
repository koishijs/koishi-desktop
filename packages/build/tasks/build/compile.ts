import { parallel, series } from 'gulp'
import { koiVersion } from '../../utils/config'
import { Exceptions } from '../../utils/exceptions'
import { dir } from '../../utils/path'
import { exec2 } from '../../utils/spawn'
import { compileShell } from './compileShell'

export const compileVersioninfo = () =>
  exec2('goversioninfo', ['-64'], dir('src'))

export const compileAppDebug = () =>
  exec2(
    'go',
    [
      'build',
      '-o',
      dir('buildPortable', process.platform === 'win32' ? 'koi.exe' : 'koi'),
      '-ldflags',
      `-X gopkg.ilharper.com/koi/app/util.AppVersion=${koiVersion.slice(1)}`,
    ],
    dir('src')
  )

export const compileAppRelease = () =>
  exec2(
    'go',
    [
      'build',
      '-o',
      dir('buildPortable', process.platform === 'win32' ? 'koi.exe' : 'koi'),
      '-trimpath',
      '-ldflags',
      `-w -s -X gopkg.ilharper.com/koi/app/util.AppVersion=${koiVersion.slice(
        1
      )}`,
    ],
    dir('src')
  )

export const compileApp = process.env.CI ? compileAppRelease : compileAppDebug

const buildCompile = () => {
  switch (process.platform) {
    case 'win32':
      return parallel(compileShell, series(compileVersioninfo, compileApp))
    case 'darwin':
      return parallel(compileShell, compileApp)
    case 'linux':
      return parallel(compileShell, compileApp)
    default:
      throw Exceptions.platformNotSupported()
  }
}

export const compile = buildCompile()
