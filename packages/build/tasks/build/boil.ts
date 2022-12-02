import { sleep } from '../../utils/common'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'

export const boil = async () => {
  await exec(
    process.platform === 'win32' ? 'koi' : './koi',
    [
      'import',
      '--name',
      'default',
      '--force',
      dir('buildCache', 'boilerplate.zip'),
    ],
    dir('buildPortable')
  )

  await sleep()
}
