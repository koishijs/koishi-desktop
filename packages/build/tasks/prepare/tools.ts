import { parallel } from 'gulp'
import {
  versionToolsGolangCILint,
  versionToolsVersioninfo,
} from '../../utils/config'
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

export const prepareTools =
  process.platform === 'win32'
    ? parallel(prepareToolsVersioninfo, prepareToolsGolangCILint)
    : parallel(prepareToolsGolangCILint)
