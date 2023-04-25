import { series } from 'gulp'
import mkdirp from 'mkdirp'
import fs from 'node:fs/promises'
import { sleep } from '../../utils/common'
import { zip } from '../../utils/compress'
import { exists } from '../../utils/fs'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'

export const packUnfoldDataCopy = async () => {
  await mkdirp(dir('buildUnfoldData', 'data'))
  await mkdirp(dir('buildUnfoldBinary'))

  // Copy all data to unfold 'data' folder
  await fs.cp(dir('buildPortableData'), dir('buildUnfoldData', 'data'), {
    recursive: true,
  })

  // Copy koi.yml to unfold root folder
  await fs.copyFile(
    dir('buildPortable', 'koi.yml'),
    dir('buildUnfoldData', 'koi.yml')
  )

  // Collect all binaries
  const binaries = (await fs.readdir(dir('buildPortable'))).filter(
    (x) => x !== 'data' && x !== 'koi.yml'
  )
  // And copy them to binary folder
  for (const b of binaries)
    await fs.cp(dir('buildPortable', b), dir('buildUnfoldBinary', b), {
      recursive: true,
    })

  // Copy Edge/WV2 setup for Windows
  if (await exists(dir('buildCache', 'MicrosoftEdgeSetup.exe')))
    await fs.copyFile(
      dir('buildCache', 'MicrosoftEdgeSetup.exe'),
      dir('buildUnfoldBinary', 'MicrosoftEdgeSetup.exe')
    )
  if (await exists(dir('buildCache', 'Webview2Setup.exe')))
    await fs.copyFile(
      dir('buildCache', 'Webview2Setup.exe'),
      dir('buildUnfoldBinary', 'Webview2Setup.exe')
    )
}

export const packUnfoldDataZip = async () => {
  await zip(dir('buildUnfoldData'), dir('srcUnfold', 'portabledata.zip'))
  await sleep()
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
      '-ldflags',
      `${process.platform === 'win32' ? '-H=windowsgui ' : ''}`,
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
      `-w -s ${process.platform === 'win32' ? '-H=windowsgui' : ''}`,
    ],
    dir('srcUnfold')
  )

export const compileUnfold = process.env.CI
  ? compileUnfoldRelease
  : compileUnfoldDebug

export const packUnfold = series(
  packUnfoldDataCopy,
  packUnfoldDataZip,
  compileUnfold
)
