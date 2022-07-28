import { series } from 'gulp'
import { koiVersion } from '../../utils/config'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'

export const compileWire = () => exec('wire', [], dir('src'))

export const compileVersioninfo = () => exec('goversioninfo', [], dir('src'))

export const compileApp = () =>
  exec(
    'go',
    [
      'build',
      '-o',
      dir('buildPortable', process.platform === 'win32' ? 'koi.exe' : 'koi'),
      '-trimpath',
      '-ldflags',
      `-w -s -X gopkg.ilharper.com/koi/app/util/const.AppVersion=${koiVersion}`,
    ],
    dir('src')
  )

export const compile = series(compileWire, compileVersioninfo, compileApp)