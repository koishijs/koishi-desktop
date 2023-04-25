import { parallel, series } from 'gulp'
import { mkdir } from 'node:fs/promises'
import { dir } from '../../../utils/path'
import { prepareNodeDownload } from './download'
import { prepareNodeExtract } from './extract'
import { extractCachePath } from './path'
import { prepareNodeYarn } from './yarn'

export * from './download'
export * from './extract'
export * from './yarn'

export const prepareNodeFolder = parallel(
  () =>
    mkdir(dir('buildPortableBin'), {
      recursive: true,
    }),

  () =>
    mkdir(extractCachePath, {
      recursive: true,
    })
)

export const prepareNode = series(
  prepareNodeFolder,
  parallel(series(prepareNodeDownload, prepareNodeExtract), prepareNodeYarn)
)
