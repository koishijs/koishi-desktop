import { series } from 'gulp'
import mkdirp from 'mkdirp'
import fs from 'node:fs/promises'
import { zip } from '../../utils/compress'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'

export const packUnfoldData = async () => {
  await mkdirp(dir('buildUnfoldData', 'data'))
  await mkdirp(dir('buildUnfoldBinary'))

  await fs.cp(dir('buildPortableData'), dir('buildUnfoldData', 'data'), {
    recursive: true,
  })
  await fs.copyFile(
    dir('buildPortable', 'koi.yml'),
    dir('buildUnfoldData', 'koi.yml')
  )

  const binaries = (await fs.readdir(dir('buildPortable'))).filter(
    (x) => x !== 'data' && x !== 'koi.yml'
  )
  for (const b of binaries)
    await fs.copyFile(dir('buildPortable', b), dir('buildUnfoldBinary', b))

  await zip(dir('buildUnfoldData'), dir('srcUnfold', 'portabledata.zip'))
}

export const compileUnfoldDebug = () =>
  exec(
    'go',
    [
      'build',
      '-o',
      dir(
        'buildUnfoldBinary',
        process.platform === 'win32' ? 'unfold.exe' : 'unfold'
      ),
    ],
    dir('srcUnfold')
  )

export const compileUnfoldRelease = () =>
  exec(
    'go',
    [
      'build',
      '-o',
      dir(
        'buildUnfoldBinary',
        process.platform === 'win32' ? 'unfold.exe' : 'unfold'
      ),
      '-trimpath',
      '-ldflags',
      '-w -s',
    ],
    dir('srcUnfold')
  )

export const compileUnfold = process.env.CI
  ? compileUnfoldRelease
  : compileUnfoldDebug

export const packUnfold = series(packUnfoldData, compileUnfold)
