import { Exceptions } from '../../utils/exceptions'
import { exec, spawnOutput } from '../../utils/spawn'
import path from 'node:path'
import fs from 'node:fs/promises'
import { dir } from '../../utils/path'

const buildCompileShellMac = (isRelease: boolean) => async () => {
  const conf = isRelease ? 'release' : 'debug'

  const buildPath = await spawnOutput(
    'swift',
    ['build', '--show-bin-path', '-c', conf],
    {
      cwd: dir('srcShellMac'),
    }
  )

  await exec('swift', ['build', '-c', conf], dir('srcShellMac'))

  await fs.copyFile(
    path.join(buildPath, 'KoiShell'),
    dir('buildPortable', 'koishell')
  )
}

const buildCompileShell = () => {
  switch (process.platform) {
    case 'win32':
      return async () => {
        /* Ignore */
      }
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
