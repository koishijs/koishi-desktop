import Handlebars from 'handlebars'
import * as fs from 'node:fs'
import * as path from 'node:path'
import {
  defaultInstance,
  koiSemver,
  koiVersionStringsJson,
} from '../../utils/config'

export const koiConfig = Handlebars.compile(
  fs.readFileSync(path.join(__dirname, 'koi-config.yml.hbs')).toString('utf-8')
)({ defaultInstance })

export const koiVersionInfo = Handlebars.compile(
  fs
    .readFileSync(path.join(__dirname, 'versioninfo.json.hbs'))
    .toString('utf-8')
)({ koiSemver, koiVersionStringsJson })

export const koiManifest = Handlebars.compile(
  fs
    .readFileSync(path.join(__dirname, 'koi.exe.manifest.hbs'))
    .toString('utf-8')
)({ koiSemver })

export const koishiManifest = Handlebars.compile(
  fs
    .readFileSync(path.join(__dirname, 'koishi.exe.manifest.hbs'))
    .toString('utf-8')
)({ koiSemver })

export const koiVisualElementsManifest = Handlebars.compile(
  fs
    .readFileSync(path.join(__dirname, 'koi.VisualElementsManifest.xml.hbs'))
    .toString('utf-8')
)({})
