import Handlebars from 'handlebars'
import * as fs from 'node:fs'
import * as path from 'node:path'
import { koiSemver, koiVersion } from '../../utils/config'
import { dir } from '../../utils/path'

export const msiWxs = Handlebars.compile(
  fs.readFileSync(path.join(__dirname, 'index.wxs.hbs')).toString('utf-8')
)({ koiVersion, koiSemver, iconPath: dir('buildAssets', 'koishi.ico') })
