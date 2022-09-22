import mkdirp from 'mkdirp'
import { promises as fs } from 'node:fs'
import { koiVersion, macAppPlist } from '../../utils/config'
import { dir } from '../../utils/path'
import { tryExec } from '../../utils/spawn'

export const packMacApp = async () => {
  const appPath = dir('buildMac', 'Koishi.app/')
  const appInfoPlistPath = dir('buildMac', 'Koishi.app/Contents/Info.plist')
  const appMacosPath = dir('buildMac', 'Koishi.app/Contents/MacOS/')
  const appResourcesPath = dir('buildMac', 'Koishi.app/Contents/Resources/')
  const appIconPath = dir(
    'buildMac',
    'Koishi.app/Contents/Resources/koishi.icns'
  )

  await mkdirp(appMacosPath)
  await mkdirp(appResourcesPath)

  await fs.cp(dir('buildPortable'), appMacosPath, { recursive: true })
  await fs.copyFile(dir('buildAssets', 'koishi.icns'), appIconPath)
  await fs.writeFile(appInfoPlistPath, macAppPlist)

  await tryExec('yarn', ['create-dmg', appPath, dir('buildMac'), '--overwrite'])
  await fs.rename(
    dir('buildMac', `Koishi ${koiVersion}.dmg`),
    dir('dist', 'koishi.dmg')
  )
}
