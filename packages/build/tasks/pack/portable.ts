import archiver from 'archiver'
import { error, warn } from 'gulplog'
import * as fs from 'node:fs'
import { dir } from '../../utils/path'

export const packPortable = async () => {
  const archive = archiver('zip', { zlib: { level: 9 } })
  archive.on('warning', warn)
  archive.on('error', error)
  archive.on('error', (x) => {
    throw x
  })

  archive.pipe(fs.createWriteStream(dir('dist', 'koishi.zip')))
  archive.directory(dir('buildPortable'), false)

  await archive.finalize()
}
