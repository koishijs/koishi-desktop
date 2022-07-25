import { parallel, series } from 'gulp'
import { prepareNodeDownload } from './download'
import { prepareNodeExtract } from './extract'
import { prepareNodeYarn } from './yarn'

export * from './download'
export * from './extract'
export * from './yarn'

export const prepareNode = parallel(
  series(prepareNodeDownload, prepareNodeExtract),
  prepareNodeYarn
)
