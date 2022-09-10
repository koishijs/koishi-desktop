import { parallel } from 'gulp'
import { info } from 'gulplog'
import {
  versionToolsGolangCILint,
  versionToolsRcedit,
  versionToolsVersioninfo,
} from '../../utils/config'
import { exists } from '../../utils/fs'
import { download } from '../../utils/net'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'

const buildPrepareTool = (name: string, version: string) => () =>
  exec('go', ['install', `${name}@${version}`])

export const prepareToolsVersioninfo = buildPrepareTool(
  'github.com/josephspurrier/goversioninfo/cmd/goversioninfo',
  versionToolsVersioninfo
)

export const prepareToolsGolangCILint = buildPrepareTool(
  'github.com/golangci/golangci-lint/cmd/golangci-lint',
  versionToolsGolangCILint
)

export const prepareToolsRcedit = async () => {
  const src = `https://github.com/electron/rcedit/releases/download/${versionToolsRcedit}/rcedit-x64.exe`
  const destFile = 'rcedit.exe'

  info('Checking temporary cache.')
  if (await exists(dir('buildCache', destFile))) return

  info('Now downloading Rcedit.')
  await download(src, dir('buildCache'), destFile)
}

export const prepareTools =
  process.platform === 'win32'
    ? parallel(
        prepareToolsVersioninfo,
        prepareToolsGolangCILint,
        prepareToolsRcedit
      )
    : parallel(prepareToolsGolangCILint)
