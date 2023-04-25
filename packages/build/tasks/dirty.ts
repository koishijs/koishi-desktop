import { koiVersion } from '../utils/config'
import { Exceptions } from '../utils/exceptions'
import { dir } from '../utils/path'
import { exec } from '../utils/spawn'

export const dirty = async () => {
  if (!koiVersion.includes('-dirty')) return

  await exec('git', ['status'], dir('root'))

  throw Exceptions.dirty()
}
