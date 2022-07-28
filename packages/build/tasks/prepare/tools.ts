import { parallel } from 'gulp'
import { versionToolsVersioninfo, versionToolsWire } from '../../utils/config'
import { exec } from '../../utils/spawn'

const buildPrepareTool = (name: string, version: string) => () =>
  exec('go', ['install', `${name}@${version}`])

export const prepareToolsWire = buildPrepareTool(
  'github.com/google/wire/cmd/wire',
  versionToolsWire
)

export const prepareToolsVersioninfo = buildPrepareTool(
  'github.com/josephspurrier/goversioninfo/cmd/goversioninfo',
  versionToolsVersioninfo
)

export const prepareTools =
  process.platform === 'win32'
    ? parallel(prepareToolsWire, prepareToolsVersioninfo)
    : parallel(prepareToolsWire)
