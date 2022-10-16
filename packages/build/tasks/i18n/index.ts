import del from 'del'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'

export const i18nExtract = async () => {
  await del(dir('locales'))

  await exec(
    'gotext',
    ['extract', '--lang=en-US,zh-CN', 'gopkg.ilharper.com/...'],
    dir('root')
  )
}
