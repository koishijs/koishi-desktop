import { parallel } from 'gulp'
import { info } from 'gulplog'
import mkdirp from 'mkdirp'
import StreamZip from 'node-stream-zip'
import { createReadStream } from 'node:fs'
import { join } from 'node:path'
import stream from 'node:stream'
import { promisify } from 'node:util'
import * as tar from 'tar'
import {
  goEnv,
  sourceGitHub,
  versionToolsGolangCILint,
  versionToolsGoText,
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

export const prepareToolsGolangCILint = async () => {
  const isWindows = goEnv.GOOS === 'windows'

  const entryName = `golangci-lint-${versionToolsGolangCILint.slice(1)}-${
    goEnv.GOOS
  }-${goEnv.GOARCH}`
  const src = `${sourceGitHub}/golangci/golangci-lint/releases/download/${versionToolsGolangCILint}/${entryName}.${
    isWindows ? 'zip' : 'tar.gz'
  }`

  const destFile = `golangci-lint.${isWindows ? 'zip' : 'tar.gz'}`
  const destDir = dir('buildCache', 'golangci-lint')

  info('Checking temporary cache.')
  if (
    await exists(
      join(destDir, isWindows ? 'golangci-lint.exe' : 'golangci-lint')
    )
  )
    return

  info('Now downloading golangci-lint.')
  await download(src, dir('buildCache'), destFile)

  info('Now extracting golangci-lint.')
  await mkdirp(destDir)
  if (isWindows) {
    const zip = new StreamZip.async({ file: dir('buildCache', destFile) })
    await zip.extract(entryName, dir('buildCache', destDir))
    await zip.close()
  } else {
    await promisify(stream.finished)(
      createReadStream(dir('buildCache', destFile)).pipe(
        tar.extract({ cwd: dir('buildCache', destDir), strip: 1 })
      )
    )
  }
}

export const prepareToolsGoText = buildPrepareTool(
  'golang.org/x/text/cmd/gotext',
  versionToolsGoText
)

export const prepareToolsRcedit = async () => {
  const src = `${sourceGitHub}/electron/rcedit/releases/download/${versionToolsRcedit}/rcedit-x64.exe`
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
        prepareToolsGoText,
        prepareToolsRcedit
      )
    : parallel(prepareToolsGolangCILint, prepareToolsGoText)
