import { info } from 'gulplog'
import { join } from 'node:path'
import { exists } from '../../../utils/fs'
import { download } from '../../../utils/net'
import { dir } from '../../../utils/path'
import { destFileYarn, srcYarn } from './path'

export const prepareNodeYarn = async () => {
  const destYarn = dir(
    'buildPortableData',
    process.platform === 'win32' ? 'node' : 'node/bin'
  )

  info('Checking temporary cache.')
  if (await exists(join(destYarn, destFileYarn))) return

  info('Now downloading Yarn.')
  await download(srcYarn, destYarn, destFileYarn)
}
