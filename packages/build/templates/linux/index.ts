import Handlebars from 'handlebars'
import * as fs from 'node:fs'
import * as path from 'node:path'
import { koiVersion } from '../../utils/config'

export const linuxAppImageDesktop = Handlebars.compile(
  fs
    .readFileSync(path.join(__dirname, 'chat.koishi.desktop.desktop.hbs'))
    .toString('utf-8')
)({ koiVersion })
