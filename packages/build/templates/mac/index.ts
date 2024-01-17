import Handlebars from 'handlebars'
import * as fs from 'node:fs'
import * as path from 'node:path'
import { koiSemver, koiVersion } from '../../utils/config'

export const macAppPlist = Handlebars.compile(
  fs.readFileSync(path.join(__dirname, 'mac-app.plist.hbs')).toString('utf-8')
)({
  date: {
    year: new Date().getFullYear(),
  },
  koiVersion,
})

export const macPkgDistribution = Handlebars.compile(
  fs
    .readFileSync(path.join(__dirname, 'distribution.xml.hbs'))
    .toString('utf-8')
)({
  date: {
    year: new Date().getFullYear(),
  },
  koiSemver,
})
