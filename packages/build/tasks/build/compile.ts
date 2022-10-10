import { parallel, series } from 'gulp'
import { koiVersion } from '../../utils/config'
import { Exceptions } from '../../utils/exceptions'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'
import { compileShell } from './compileShell'

export const compileVersioninfo = () =>
  exec('goversioninfo', ['-64'], dir('src'))

export const compileAppDebug = () =>
  exec(
    'go',
    [
      'build',
      '-o',
      dir('buildPortable', process.platform === 'win32' ? 'koi.exe' : 'koi'),
      '-ldflags',
      `${
        process.platform === 'win32' ? '-H=windowsgui ' : ''
      }-X gopkg.ilharper.com/koi/app/util.AppVersion=${koiVersion.slice(1)}`,
    ],
    dir('src')
  )

export const compileAppRelease = () =>
  exec(
    'go',
    [
      'build',
      '-o',
      dir('buildPortable', process.platform === 'win32' ? 'koi.exe' : 'koi'),
      '-trimpath',
      '-ldflags',
      `-w -s ${
        process.platform === 'win32' ? '-H=windowsgui' : ''
      } -X gopkg.ilharper.com/koi/app/util.AppVersion=${koiVersion.slice(1)}`,
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
