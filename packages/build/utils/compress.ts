import archiver from 'archiver'
import { error, warn } from 'gulplog'
import * as fs from 'node:fs'

export const zip = async (src: string, dest: string) => {
  const archive = archiver('zip', { zlib: { level: 9 } })
  archive.on('warning', warn)
  archive.on('error', error)
  archive.on('error', (x) => {
    throw x
  })

  archive.pipe(fs.createWriteStream(dest))
  archive.directory(src, false)

  await archive.finalize()
}
