import Handlebars from 'handlebars'
import * as fs from 'node:fs'
import * as path from 'node:path'
import { koiSemver, koiVersion } from '../../utils/config'

export const koiShellResources = Handlebars.compile(
  fs.readFileSync(path.join(__dirname, 'koishell.rc.hbs')).toString('utf-8')
)({
  date: {
    year: new Date().getFullYear(),
  },
  koiVersion,
  koiSemver,
})

export const koiShellManifest = Handlebars.compile(
  fs
    .readFileSync(path.join(__dirname, 'koishell.exe.manifest.hbs'))
    .toString('utf-8')
)({
  date: {
    year: new Date().getFullYear(),
  },
  koiSemver,
})
