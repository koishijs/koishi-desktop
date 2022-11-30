import del from 'del'
import { dir } from '../../utils/path'

export const cleanTemp = async () => {
  const dirs = [
    dir('buildPortableData', 'home'),
    dir('buildPortableData', 'lock'),
    dir('buildPortableData', 'logs'),
    dir('buildPortableData', 'tmp'),
  ]
  for (const d of dirs) await del(d)
}
