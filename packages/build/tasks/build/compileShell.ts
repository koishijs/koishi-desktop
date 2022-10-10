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
    [
      `-DCMAKE_BUILD_TYPE=${conf}`,
      '-G',
      'Visual Studio 17 2022',
      dir('srcShellWin'),
    ],
    dir('buildShellWin')
  )

  await exec(
    'cmake',
    [`--build`, '.', '--target', 'koishell'],
    dir('buildShellWin')
  )

  await fs.copyFile(
    dir('buildShellWin', 'koishell.exe'),
    dir('buildPortable', 'koishell.exe')
  )
}

const buildCompileShellMac = (isRelease: boolean) => async () => {
  const conf = isRelease ? 'release' : 'debug'

  const buildPath = (
    await spawnOutput('swift', ['build', '--show-bin-path', '-c', conf], {
      cwd: dir('srcShellMac'),
    })
  ).trim()

  await exec('swift', ['build', '-c', conf], dir('srcShellMac'))

  await fs.copyFile(
    path.join(buildPath, 'KoiShell'),
    dir('buildPortable', 'koishell')
  )
}

const buildCompileShell = () => {
  switch (process.platform) {
    case 'win32':
      return buildCompileShellWin(Boolean(process.env.CI))
    case 'darwin':
      return buildCompileShellMac(Boolean(process.env.CI))
    case 'linux':
      return async () => {
        /* Ignore */
      }
    default:
      throw Exceptions.platformNotSupported()
  }
}

export const compileShell = buildCompileShell()
