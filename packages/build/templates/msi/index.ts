import Handlebars from 'handlebars'
import fs from 'node:fs'
import path from 'node:path'

export const msiWxsHbs = Handlebars.compile(
  fs.readFileSync(path.join(__dirname, 'index.wxs.hbs')).toString('utf-8')
)
