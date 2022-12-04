// https://stackoverflow.com/a/24470998

import { Exceptions } from '../../utils/exceptions'
import { exec, spawnOutput } from '../../utils/spawn'
import path from 'node:path'
import fs from 'node:fs/promises'
import { dir } from '../../utils/path'
import mkdirp from 'mkdirp'

const buildCompileShellWin = (isRelease: boolean) => async () => {
  const conf = isRelease ? 'MinSizeRel' : 'Debug'

  await mkdirp(dir('buildShellWin'))

  await exec(
    'cmake',
    ['-G', 'Visual Studio 17 2022', dir('srcShellWin')],
    dir('buildShellWin')
  )

  await exec(
    'cmake',
    [`--build`, '.', '--target', 'koishell', '--config', conf],
    dir('buildShellWin')
  )

  await fs.copyFile(
    dir('buildShellWin', `${conf}/koishell.exe`),
    dir('buildPortable', 'koishell.exe')
  )
}

const buildCompileShellMac = (isRelease: boolean) => async () => {
  const conf = isRelease ? 'release' : 'debug'
  const distPath = dir('buildPortable', 'koishell')
  const distBundlePath = dir('buildPortable', 'KoiShell_KoiShell.bundle/')

  await fs.rm(distPath, { force: true })

  const buildPath = (
    await spawnOutput('swift', ['build', '--show-bin-path', '-c', conf], {
      cwd: dir('srcShellMac'),
    })
  ).trim()

  await exec('swift', ['build', '-c', conf], dir('srcShellMac'))

  await fs.copyFile(path.join(buildPath, 'KoiShell'), distPath)
  await mkdirp(distBundlePath)
  await fs.cp(
    path.join(buildPath, 'KoiShell_KoiShell.bundle/'),
    distBundlePath,
    { recursive: true }
  )
}

const buildCompileShellLinux = (isRelease: boolean) => async () => {
  const conf = isRelease ? 'MinSizeRel' : 'Debug'

  await mkdirp(dir('buildShellLinux'))

  await exec(
    'cmake',
    [
      `-DCMAKE_BUILD_TYPE=${conf}`,
      '-G',
      'Unix Makefiles',
      dir('srcShellLinux'),
    ],
    dir('buildShellLinux')
  )

  await exec(
    'cmake',
    [`--build`, '.', '--target', 'koishell'],
    dir('buildShellLinux')
  )

  await fs.copyFile(
    dir('buildShellLinux', 'koishell'),
    dir('buildPortable', 'koishell')
  )
}

const buildCompileShell = () => {
  switch (process.platform) {
    case 'win32':
      return buildCompileShellWin(process.env.CI ? true : false)
    case 'darwin':
      return buildCompileShellMac(process.env.CI ? true : false)
    case 'linux':
      return buildCompileShellLinux(process.env.CI ? true : false)
    default:
      throw Exceptions.platformNotSupported()
  }
}

export const compileShell = buildCompileShell()
