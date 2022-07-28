import { path7za } from '7zip-bin'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'

export const packPortable = () =>
  exec(path7za, ['a', '../../dist/koi.7z', 'koi', '-mx9'], dir('buildPortable'))
