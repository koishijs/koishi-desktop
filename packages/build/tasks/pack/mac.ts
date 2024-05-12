import { series } from 'gulp'
import mkdirp from 'mkdirp'
import { promises as fs } from 'node:fs'
import { macAppPlist, macPkgDistribution } from '../../templates'
import { koiVersion } from '../../utils/config'
import { dir } from '../../utils/path'
import { exec, tryExec } from '../../utils/spawn'

const appPath = dir('buildMac', 'Koishi.app/')

export const packMacApp = async () => {
  const appInfoPlistPath = dir('buildMac', 'Koishi.app/Contents/Info.plist')
  const appMacosPath = dir('buildMac', 'Koishi.app/Contents/MacOS/')
  const appResourcesPath = dir('buildMac', 'Koishi.app/Contents/Resources/')
  const appIconPath = dir(
    'buildMac',
    'Koishi.app/Contents/Resources/koishi-app.icns'
  )

  await mkdirp(appMacosPath)
  await mkdirp(appResourcesPath)

  await fs.cp(dir('buildUnfoldBinary'), appMacosPath, { recursive: true })
  await fs.rename(
    dir('buildMac', 'Koishi.app/Contents/MacOS/KoiShell_KoiShell.bundle/'),
    dir('buildMac', 'Koishi.app/KoiShell_KoiShell.bundle/')
  )
  await fs.copyFile(dir('buildAssets', 'koishi-app.icns'), appIconPath)
  await fs.writeFile(appInfoPlistPath, macAppPlist)
}

export const packMacDmg = async () => {
  await fs.writeFile(
    dir('buildMac', 'Koishi.app/Contents/MacOS/koi.yml'),
    'redirect: USERDATA'
  )
  await tryExec('yarn', [
    'workspace',
    'koi-build',
    'create-dmg',
    appPath,
    dir('buildMac'),
    '--overwrite',
  ])
  await fs.rename(
    dir('buildMac', `Koishi ${koiVersion}.dmg`),
    dir('dist', 'koishi.dmg')
  )
}

export const packMacPkg = async () => {
  const scriptsPath = dir('buildMac', 'scripts/')
  const postinstallPath = dir('buildMac', 'scripts/postinstall')

  await mkdirp(scriptsPath)

  await fs.writeFile(dir('buildMac', 'distribution.xml'), macPkgDistribution)
  await fs.copyFile(dir('templates', 'mac/postinstall.sh'), postinstallPath)
  await fs.chmod(postinstallPath, 0o755)

  await exec(
    'pkgbuild',
    [
      '--identifier',
      'chat.koishi.desktop',
      '--component',
      'Koishi.app',
      '--scripts',
      'scripts',
      '--install-location',
      '/Applications',
      'koishi-app.pkg',
    ],
    dir('buildMac')
  )

  await exec(
    'productbuild',
    [
      '--distribution',
      'distribution.xml',
      '--package-path',
      '.',
      dir('dist', 'koishi.pkg'),
    ],
    dir('buildMac')
  )
}

// Series: packMacApp => packMacPkg => packMacDmg
// packMacDmg writes `koi.yml` to dir('buildMac') so it needs to be executed at last.
export const packMac = series(packMacApp, packMacPkg, packMacDmg)
