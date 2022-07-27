import axios from 'axios'
import { info } from 'gulplog'
import mkdirp from 'mkdirp'
import { createWriteStream } from 'node:fs'
import { join } from 'node:path'
import stream from 'node:stream'
import { promisify } from 'node:util'

export async function download(src: string, dest: string, filename: string) {
  await mkdirp(dest)
  const res = await axios.get(src, {
    responseType: 'stream',
    onDownloadProgress(progressEvent) {
      if (progressEvent.lengthComputable)
        info(
          `${filename} (${(
            (progressEvent.loaded * 100) /
            progressEvent.total
          ).toFixed(2)}%)`
        )
    },
  })
  const writeStream = createWriteStream(join(dest, filename))
  await promisify(stream.finished)(
    (res.data as stream.Readable).pipe(writeStream)
  )
}
