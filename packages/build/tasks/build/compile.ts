import { series } from 'gulp'
import { koiVersion } from '../../utils/config'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'

export const compileVersioninfo = () =>
  exec('goversioninfo', ['-64'], dir('src'))

export const compileAppDebug = () =>
  exec(
    'go',
    [
      'build',
      '-o',
      dir('buildPortable', process.platform === 'win32' ? 'koi.exe' : 'koi'),
      '-trimpath',
      '-ldflags',
      `-X gopkg.ilharper.com/koi/app/util.AppVersion=${koiVersion.slice(1)}`,
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
      `-w -s -X gopkg.ilharper.com/koi/app/util.AppVersion=${koiVersion.slice(
        1
      )}`,
    ],
    dir('src')
  )

export const compileApp = process.env.CI ? compileAppRelease : compileAppDebug

export const compile =
  process.platform === 'win32'
    ? series(compileVersioninfo, compileApp)
    : series(compileApp)
