import { koiVersion } from '../utils/config'
import { Exceptions } from '../utils/exceptions'
import { dir } from '../utils/path'
import { exec2 } from '../utils/spawn'

export const dirty = async () => {
  if (!koiVersion.includes('-dirty')) return

  await exec2('git', ['status'], dir('root'))

  throw Exceptions.dirty()
}
