import del from 'del'
import { series } from 'gulp'
import fs from 'node:fs/promises'
import { dir } from '../../utils/path'
import { exec2, tryExec } from '../../utils/spawn'

export const i18nExtract = () =>
  exec2(
    'gotext',
    ['extract', '--lang=en-US,zh-CN', 'gopkg.ilharper.com/...'],
    dir('root')
  )

export const i18nGenerate = async () => {
  const pathCatalog = dir('src', 'catalog.go')
  await del(pathCatalog)

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

  let catalog = (await fs.readFile(pathCatalog)).toString()
  catalog = catalog.replace('package ', 'package main')
  await fs.writeFile(pathCatalog, catalog)
}

export const i18n = series(i18nExtract, i18nGenerate)
