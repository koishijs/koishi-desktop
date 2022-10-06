import { series } from 'gulp'
import mkdirp from 'mkdirp'
import { promises as fs } from 'node:fs'
import { macAppPlist, macPkgDistribution } from '../../templates'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'

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
  await fs.copyFile(dir('buildAssets', 'koishi-app.icns'), appIconPath)
  await fs.writeFile(appInfoPlistPath, macAppPlist)
}

export const packMacPkg = async () => {
  const scriptsPath = dir('buildMac', 'scripts/')
  const postinstallPath = dir('buildMac', 'scripts/postinstall')

  await mkdirp(scriptsPath)

  await fs.writeFile(dir('buildMac', 'distribution.xml'), macPkgDistribution)
  await fs.writeFile(
    postinstallPath,
    `
#!/bin/bash
echo "Starting post-install process..."
echo "Removing com.apple.quarantine..."
sudo xattr -d com.apple.quarantine /Applications/Koishi.app/ || true
echo "Setting chmod for unfold..."
sudo chmod -R 777 /Applications/Koishi.app/ || true
echo "Starting unfold..."
sudo /Applications/Koishi.app/Contents/MacOS/unfold ensure
echo "Setting chmod for app..."
sudo chmod -R 755 /Applications/Koishi.app/ || true
echo "Setting chmod for user data..."
sudo chmod -R 777 ~/Library/Application\\ Support/Il\\ Harper/Koishi/ || true
echo "Post-install process finished."
`.trim()
  )
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

export const packMac = series(packMacApp, packMacPkg)
