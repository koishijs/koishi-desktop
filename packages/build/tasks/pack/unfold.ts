import { series } from 'gulp'
import mkdirp from 'mkdirp'
import fs from 'node:fs/promises'
import { join } from 'node:path'
import { sleep } from '../../utils/common'
import { zip } from '../../utils/compress'
import { versionMSVC } from '../../utils/config'
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

  if (process.platform === 'win32') {
    // Copy Edge/WV2 setup
    await fs.copyFile(
      dir('buildCache', 'MicrosoftEdgeSetup.exe'),
      dir('buildUnfoldBinary', 'MicrosoftEdgeSetup.exe')
    )
    await fs.copyFile(
      dir('buildCache', 'Webview2Setup.exe'),
      dir('buildUnfoldBinary', 'Webview2Setup.exe')
    )

    // Copy VCRedist
    await fs.copyFile(
      join(
        process.env['PROGRAMFILES']!,
        `Microsoft Visual Studio/2022/Enterprise/VC/Redist/MSVC/${versionMSVC}/MergeModules/Microsoft_VC143_CRT_x64.msm`
      ),
      dir('buildUnfoldBinary', 'Microsoft_VC143_CRT_x64.msm')
    )
  }
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
