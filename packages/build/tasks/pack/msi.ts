import { parallel, series } from 'gulp'
import { info } from 'gulplog'
import mkdirp from 'mkdirp'
import fs from 'node:fs/promises'
import { msiWxsHbs } from '../../templates'
import { koiSemver, koiVersion } from '../../utils/config'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'
import path from 'node:path'

const dirVarients = dir('buildMsi', 'varients/')
const pathNeutralMsi = dir('buildMsi', 'koishi.msi')

const langFilenameRegexp = /WixUI_.*\.wxl/i

// https://learn.microsoft.com/ja-jp/windows-hardware/manufacture/desktop/available-language-packs-for-windows
const langs = [
  { lcid: '1025', locale: 'ar-SA', codepage: '1256' },
  { lcid: '1026', locale: 'bg-BG', codepage: '1251' },
  { lcid: '1027', locale: 'ca-ES', codepage: '1252' },
  { lcid: '1029', locale: 'cs-CZ', codepage: '1250' },
  { lcid: '1030', locale: 'da-DK', codepage: '1252' },
  { lcid: '1031', locale: 'de-de', codepage: '1252' },
  { lcid: '1032', locale: 'el-GR', codepage: '1253' },
  { lcid: '1033', locale: 'en-us', codepage: '1252' },
  { lcid: '3082', locale: 'es-es', codepage: '1252' },
  { lcid: '1061', locale: 'et-EE', codepage: '1257' },
  { lcid: '1035', locale: 'fi-FI', codepage: '1252' },
  { lcid: '1036', locale: 'fr-fr', codepage: '1252' },
  { lcid: '1037', locale: 'he-IL', codepage: '1255' },
  { lcid: '1081', locale: 'hi-IN', codepage: /* '0' 0????? */ '1252' },
  { lcid: '1050', locale: 'hr-HR', codepage: '1250' },
  { lcid: '1038', locale: 'hu-HU', codepage: '1250' },
  { lcid: '1040', locale: 'it-it', codepage: '1252' },
  { lcid: '1041', locale: 'ja-jp', codepage: '932' },
  { lcid: '1087', locale: 'kk-KZ', codepage: '1251' },
  { lcid: '1042', locale: 'ko-KR', codepage: '949' },
  { lcid: '1063', locale: 'lt-LT', codepage: '1257' },
  { lcid: '1062', locale: 'lv-LV', codepage: '1257' },
  { lcid: '1044', locale: 'nb-NO', codepage: '1252' },
  { lcid: '1043', locale: 'nl-NL', codepage: '1252' },
  { lcid: '1045', locale: 'pl-pl', codepage: '1250' },
  { lcid: '1046', locale: 'pt-BR', codepage: '1252' },
  { lcid: '2070', locale: 'pt-PT', codepage: '1252' },
  { lcid: '1048', locale: 'ro-RO', codepage: '1250' },
  { lcid: '1049', locale: 'ru-ru', codepage: '1251' },
  { lcid: '1051', locale: 'sk-SK', codepage: '1250' },
  { lcid: '1060', locale: 'sl-SI', codepage: '1250' },
  { lcid: '1052', locale: 'sq-AL', codepage: '1252' },
  { lcid: '2074', locale: 'sr-Latn-CS', codepage: '1250' },
  { lcid: '1053', locale: 'sv-SE', codepage: '1252' },
  { lcid: '1054', locale: 'th-TH', codepage: '874' },
  { lcid: '1055', locale: 'tr-TR', codepage: '1254' },
  { lcid: '1058', locale: 'uk-UA', codepage: '1251' },
  { lcid: '2052', locale: 'zh-CN', codepage: '936' },
  { lcid: '3076', locale: 'zh-HK', codepage: '950' },
  { lcid: '1028', locale: 'zh-TW', codepage: '950' },
] as const

type Lang = typeof langs[number]

export const packMsiGenerateLangs = async () => {
  const dirWixlib = dir(
    'buildVendor',
    'wix3-develop/src/ext/UIExtension/wixlib'
  )

  const langFilenames = (await fs.readdir(dirWixlib)).filter((x) =>
    langFilenameRegexp.test(x)
  )

  const result: { lcid: string; locale: string; codepage: string }[] = []

  for (const langFilename of langFilenames) {
    const langFile = (
      await fs.readFile(path.join(dirWixlib, langFilename))
    ).toString()
    result.push({
      lcid: '',
      locale: langFilename.slice(6, langFilename.length - 4),
      codepage: langFile.split('" Codepage="')[1].split('" xmlns="')[0],
    })
  }

  info(result)
}

export const packMsiMkdir = () => mkdirp(dirVarients)

export const packMsiIndex = () =>
  fs.writeFile(
    dir('buildMsi', 'index.wxs'),
    msiWxsHbs({
      koiVersion,
      koiSemver,
      iconPath: dir('buildAssets', 'koishi.ico'),
      languages: langs.map((x) => x.lcid).join(),
    })
  )

export const packMsiFiles = async () => {
  const dirSource = dir('buildMsi', 'SourceDir/')
  await mkdirp(dirSource)
  await fs.cp(dir('buildUnfoldBinary'), dirSource, { recursive: true })
}

export const packMsiCandle = () =>
  exec(
    dir('buildVendor', 'wix/candle.exe'),
    ['-nologo', 'index.wxs'],
    dir('buildMsi')
  )

export const packMsiLight = () =>
  exec(
    dir('buildVendor', 'wix/light.exe'),
    [
      '-nologo',
      '-sice:ICE61',
      '-sice:ICE69',
      '-spdb',
      '-out',
      pathNeutralMsi,
      '-ext',
      'WixUIExtension',
      'index.wixobj',
    ],
    dir('buildMsi')
  )

const buildPackMsiTrans = (lang: Lang) => async () => {
  const pathVarientMsi = path.join(dirVarients, `${lang.lcid}.msi`)
  const pathMst = path.join(dirVarients, `${lang.lcid}.mst`)

  await fs.copyFile(pathNeutralMsi, pathVarientMsi)

  await exec(
    'cscript',
    ['wilangid.vbs', pathVarientMsi, 'Product', lang.lcid],
    dirVarients
  )

  await exec(
    'MsiTran.exe',
    ['-g', pathNeutralMsi, pathVarientMsi, pathMst],
    dirVarients
  )
}

const buildPackMsiSubstorage = (lang: Lang) => async () => {
  const pathMst = path.join(dirVarients, `${lang.lcid}.mst`)

  await exec(
    'cscript',
    ['wisubstg.vbs', pathNeutralMsi, pathMst, lang.lcid],
    dirVarients
  )
}

export const packMsiListStorage = () =>
  exec('cscript', ['wisubstg.vbs', pathNeutralMsi], dirVarients)

export const packMsiCopyDist = () =>
  fs.copyFile(pathNeutralMsi, dir('dist', 'koishi.msi'))

export const packMsi = series(
  packMsiMkdir,
  parallel(packMsiIndex, packMsiFiles),
  packMsiCandle,
  packMsiLight,
  parallel(langs.map(buildPackMsiTrans)),
  series(langs.map(buildPackMsiSubstorage)),
  packMsiListStorage,
  packMsiCopyDist
)
