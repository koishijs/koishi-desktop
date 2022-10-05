import { zip } from '../../utils/compress'
import { dir } from '../../utils/path'

export const packPortable = () =>
  zip(dir('buildPortable'), dir('dist', 'koishi.zip'))
