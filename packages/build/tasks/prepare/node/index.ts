import { parallel, series } from 'gulp'
import { mkdir } from 'node:fs/promises'
import { dir } from '../../../utils/path'
import { prepareNodeDownload } from './download'
import { prepareNodeExtract } from './extract'
import { prepareNodeYarn } from './yarn'

export * from './download'
export * from './extract'
export * from './yarn'

export const prepareNodeFolder = () =>
  mkdir(
    dir(
      'buildPortableData',
      process.platform === 'win32' ? 'node' : 'node/bin'
    ),
    {
      recursive: true,
    }
  )

export const prepareNode = series(
  prepareNodeFolder,
  parallel(series(prepareNodeDownload, prepareNodeExtract), prepareNodeYarn)
)
