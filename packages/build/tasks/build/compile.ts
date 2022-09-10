import { series } from 'gulp'
import { koiVersion } from '../../utils/config'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'
import execa from 'execa'

export const compileVersioninfo = () =>
  exec('goversioninfo', ['-64'], dir('src'))

const compileAppExecOptions: execa.SyncOptions =
  process.platform === 'win32'
    ? {
        env: {
          CGO_CXXFLAGS: `-I${dir('buildVendor', 'WebView2/include')}`,
          CGO_LDFLAGS: `-L${dir('buildVendor', 'WebView2/x64')}`,
        },
      }
    : {}

export const compileAppDebug = () =>
  exec(
    'go',
    [
      'build',
      '-o',
      dir('buildPortable', process.platform === 'win32' ? 'koi.exe' : 'koi'),
      '-x',
      '-v',
      '-ldflags',
      `${
        process.platform === 'win32' ? '-H=windowsgui ' : ''
      }-X gopkg.ilharper.com/koi/app/util.AppVersion=${koiVersion.slice(1)}`,
    ],
    dir('src'),
    compileAppExecOptions
  )

export const compileAppRelease = () =>
  exec(
    'go',
    [
      'build',
      '-o',
      dir('buildPortable', process.platform === 'win32' ? 'koi.exe' : 'koi'),
      '-x',
      '-v',
      '-trimpath',
      '-ldflags',
      `-w -s ${
        process.platform === 'win32' ? '-H=windowsgui' : ''
      } -X gopkg.ilharper.com/koi/app/util.AppVersion=${koiVersion.slice(1)}`,
    ],
    dir('src'),
    compileAppExecOptions
  )

export const compileApp = process.env.CI ? compileAppRelease : compileAppDebug

export const compile =
  process.platform === 'win32'
    ? series(compileVersioninfo, compileApp)
    : series(compileApp)
