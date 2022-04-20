import { path7za } from '7zip-bin'
import { error } from 'gulplog'
import { resolve } from './path'
import { spawnAsync } from './utils'

export async function pack(): Promise<void> {
  const result = await spawnAsync(path7za, ['a', 'koi.7z', 'koi', '-mx9'], {
    cwd: resolve('.', 'build'),
  })
  if (result) {
    const err = `'7za a' exited with error code: ${result}`
    error(err)
    throw new Error(err)
  }
}
