import del from 'del'
import fs from 'node:fs/promises'
import { dir } from '../../utils/path'
import { exec, tryExec } from '../../utils/spawn'

export const i18nExtract = async () => {
  await del(dir('locales'))

  await exec(
    'gotext',
    ['extract', '--lang=en-US,zh-CN', 'gopkg.ilharper.com/...'],
    dir('root')
  )
}

export const i18nGenerate = async () => {
  await tryExec(
    'gotext',
    [
      '--srclang=en-US',
      'update',
      '--lang=en-US,zh-CN',
      '--out=packages/app/catalog.go',
      'gopkg.ilharper.com/...',
    ],
    dir('root')
  )

  const pathCatalog = dir('src', 'catalog.go')
  const catalog = (await fs.readFile(pathCatalog)).toString()
  catalog.replace('package ', 'package main')
  await fs.writeFile(pathCatalog, catalog)
}
