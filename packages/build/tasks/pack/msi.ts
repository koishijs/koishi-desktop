import { parallel, series } from 'gulp'
import { info } from 'gulplog'
import mkdirp from 'mkdirp'
import fs from 'node:fs/promises'
import path from 'node:path'
import { v4 as uuid } from 'uuid'
import { msiWxsHbs } from '../../templates'
import { koiSemver, koiVersion } from '../../utils/config'
import { dir } from '../../utils/path'
import { exec2 } from '../../utils/spawn'

const buildProductGuid = () => uuid().toUpperCase()

const getPathWiLangId = () =>
  path.join(
    process.env['PROGRAMFILES(X86)']!,
    '/Windows Kits/10/bin/10.0.22000.0/x64/wilangid.vbs'
  )
const getPathWiSubStg = () =>
  path.join(
    process.env['PROGRAMFILES(X86)']!,
    '/Windows Kits/10/bin/10.0.22000.0/x64/wisubstg.vbs'
  )
const getPathMsiTran = () =>
  path.join(
    process.env['PROGRAMFILES(X86)']!,
    '/Windows Kits/10/bin/10.0.22000.0/x86/MsiTran.exe'
  )

const langFilenameRegexp = /WixUI_.*\.wxl/i

// https://learn.microsoft.com/ja-jp/windows-hardware/manufacture/desktop/available-language-packs-for-windows
const langs = [
  // { lcid: '1025', locale: 'ar-SA', codepage: '1256' },
  // { lcid: '1026', locale: 'bg-BG', codepage: '1251' },
  // { lcid: '1027', locale: 'ca-ES', codepage: '1252' },
  // { lcid: '1029', locale: 'cs-CZ', codepage: '1250' },
  // { lcid: '1030', locale: 'da-DK', codepage: '1252' },
  // { lcid: '1031', locale: 'de-DE', codepage: '1252' },
  // { lcid: '1032', locale: 'el-GR', codepage: '1253' },
  { lcid: '1033', locale: 'en-US', codepage: '1252' },
  // { lcid: '3082', locale: 'es-ES', codepage: '1252' },
  // { lcid: '1061', locale: 'et-EE', codepage: '1257' },
  // { lcid: '1035', locale: 'fi-FI', codepage: '1252' },
  { lcid: '1036', locale: 'fr-FR', codepage: '1252' },
  // { lcid: '1037', locale: 'he-IL', codepage: '1255' },
  // { lcid: '1081', locale: 'hi-IN', codepage: /* '0' 0????? */ '1252' },
  // { lcid: '1050', locale: 'hr-HR', codepage: '1250' },
  // { lcid: '1038', locale: 'hu-HU', codepage: '1250' },
  // { lcid: '1040', locale: 'it-IT', codepage: '1252' },
  { lcid: '1041', locale: 'ja-JP', codepage: '932' },
  // { lcid: '1087', locale: 'kk-KZ', codepage: '1251' },
  // { lcid: '1042', locale: 'ko-KR', codepage: '949' },
  // { lcid: '1063', locale: 'lt-LT', codepage: '1257' },
  // { lcid: '1062', locale: 'lv-LV', codepage: '1257' },
  // { lcid: '1044', locale: 'nb-NO', codepage: '1252' },
  // { lcid: '1043', locale: 'nl-NL', codepage: '1252' },
  // { lcid: '1045', locale: 'pl-PL', codepage: '1250' },
  // { lcid: '1046', locale: 'pt-BR', codepage: '1252' },
  // { lcid: '2070', locale: 'pt-PT', codepage: '1252' },
  // { lcid: '1048', locale: 'ro-RO', codepage: '1250' },
  // { lcid: '1049', locale: 'ru-RU', codepage: '1251' },
  // { lcid: '1051', locale: 'sk-SK', codepage: '1250' },
  // { lcid: '1060', locale: 'sl-SI', codepage: '1250' },
  // { lcid: '1052', locale: 'sq-AL', codepage: '1252' },
  // { lcid: '2074', locale: 'sr-Latn-CS', codepage: '1250' },
  // { lcid: '1053', locale: 'sv-SE', codepage: '1252' },
  // { lcid: '1054', locale: 'th-TH', codepage: '874' },
  // { lcid: '1055', locale: 'tr-TR', codepage: '1254' },
  // { lcid: '1058', locale: 'uk-UA', codepage: '1251' },
  { lcid: '2052', locale: 'zh-CN', codepage: '936' },
  { lcid: '3076', locale: 'zh-HK', codepage: '950' },
  { lcid: '1028', locale: 'zh-TW', codepage: '950' },
] as const

const varientLangs = langs.filter((x) => x.lcid !== '1033')

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

export const packMsiMkdir = () => mkdirp(dir('buildMsi'))

export const packMsiFiles = async () => {
  const dirSource = dir('buildMsi', 'SourceDir/')
  await mkdirp(dirSource)
  await fs.cp(dir('buildUnfoldBinary'), dirSource, { recursive: true })
}

const buildPackMsi = (lang: Lang) => async () => {
  const pathWxs = dir('buildMsi', `${lang.lcid}.wxs`)
  const pathWixobj = dir('buildMsi', `${lang.lcid}.wixobj`)
  const pathMsi = dir('buildMsi', `${lang.lcid}.msi`)

  await fs.writeFile(
    pathWxs,
    msiWxsHbs({
      koiVersion,
      koiSemver,
      iconPath: dir('buildAssets', 'koishi.ico'),
      productGuid: buildProductGuid(),
      language: lang.lcid,
      codepage: lang.codepage,
    })
  )

  await exec2(
    dir('buildVendor', 'wix/candle.exe'),
    ['-nologo', pathWxs],
    dir('buildMsi')
  )

  await exec2(
    dir('buildVendor', 'wix/light.exe'),
    [
      '-nologo',
      '-sice:ICE61',
      '-sice:ICE69',
      '-spdb',
      '-out',
      pathMsi,
      `-cultures:${lang.locale.toLowerCase()}`,
      '-ext',
      'WixUIExtension',
      pathWixobj,
    ],
    dir('buildMsi')
  )
}

const pathNeutralMsi = dir('buildMsi', 'koishi.msi')

export const packMsiCopyNeutral = () =>
  fs.copyFile(dir('buildMsi', '1033.msi'), pathNeutralMsi)

const buildPackMsiTrans = (lang: Lang) => async () => {
  const pathMsi = dir('buildMsi', `${lang.lcid}.msi`)
  const pathMst = dir('buildMsi', `${lang.lcid}.mst`)

  await exec2(
    'cscript',
    [getPathWiLangId(), pathMsi, 'Product', lang.lcid],
    dir('buildMsi')
  )

  await exec2(
    getPathMsiTran(),
    ['-g', pathNeutralMsi, pathMsi, pathMst],
    dir('buildMsi')
  )
}

const buildPackMsiSubstorage = (lang: Lang) => async () => {
  const pathMst = dir('buildMsi', `${lang.lcid}.mst`)

  await exec2(
    'cscript',
    [getPathWiSubStg(), pathNeutralMsi, pathMst, lang.lcid],
    dir('buildMsi')
  )
}

export const packMsiWriteLangIds = () =>
  exec2(
    'cscript',
    [
      getPathWiLangId(),
      pathNeutralMsi,
      'Package',
      langs.map((x) => x.lcid).join(),
    ],
    dir('buildMsi')
  )

export const packMsiListStorage = () =>
  exec2('cscript', [getPathWiSubStg(), pathNeutralMsi], dir('buildMsi'))

export const packMsiCopyDist = () =>
  fs.copyFile(pathNeutralMsi, dir('dist', 'koishi.msi'))

export const packMsi = series(
  packMsiMkdir,
  packMsiFiles,
  parallel(langs.map(buildPackMsi)),
  packMsiCopyNeutral,
  parallel(varientLangs.map(buildPackMsiTrans)),
  series(varientLangs.map(buildPackMsiSubstorage)),
  packMsiWriteLangIds,
  packMsiListStorage,
  packMsiCopyDist
)
