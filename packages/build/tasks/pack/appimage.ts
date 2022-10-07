import { parallel, series } from 'gulp'
import mkdirp from 'mkdirp'
import { dir } from '../../utils/path'
import path from 'node:path'
import fs from 'node:fs/promises'
import { exec } from '../../utils/spawn'
import { linuxAppImageDesktop } from '../../templates/linux'

const appDirPath = dir('buildLinux', 'Koishi.AppDir/')
const appBinaryPath = path.join(appDirPath, 'usr/bin/')
const appMetaInfoPath = path.join(appDirPath, 'usr/share/metainfo/')

export const packAppImageMkdir = parallel(
  () => mkdirp(appBinaryPath),
  () => mkdirp(appMetaInfoPath)
)

export const packAppImageCopyFiles = parallel(
  series(
    () =>
      fs.copyFile(
        dir('templates', 'linux/AppRun'),
        path.join(appDirPath, 'AppRun')
      ),
    () => fs.chmod(path.join(appDirPath, 'AppRun'), 0o755)
  ),
  () =>
    fs.writeFile(
      path.join(appDirPath, 'chat.koishi.desktop.desktop'),
      linuxAppImageDesktop
    ),
  () =>
    fs.copyFile(
      dir('templates', 'linux/chat.koishi.desktop.appdata.xml'),
      path.join(appMetaInfoPath, 'chat.koishi.desktop.appdata.xml')
    ),
  () =>
    fs.copyFile(
      dir('templates', 'linux/chat.koishi.desktop.appdata.xml'),
      path.join(appMetaInfoPath, 'chat.koishi.desktop.metainfo.xml')
    ),
  () =>
    fs.copyFile(
      dir('templates', 'linux/chat.koishi.desktop.appdata.xml'),
      path.join(appMetaInfoPath, 'chat.koishi.appdata.xml.appdata.xml')
    ),
  () =>
    fs.copyFile(
      dir('buildAssets', 'koishi.png'),
      path.join(appDirPath, '.DirIcon')
    ),
  () =>
    fs.copyFile(
      dir('buildAssets', 'koishi.png'),
      path.join(appDirPath, 'chat.koishi.desktop.png')
    ),
  () =>
    fs.copyFile(
      dir('buildAssets', 'koishi.svg'),
      path.join(appDirPath, 'chat.koishi.desktop.svg')
    ),
  () => fs.cp(dir('buildPortable'), appBinaryPath, { recursive: true })
)

export const packAppImageGenerate = () =>
  exec(
    dir('buildCache', 'appimagetool.AppImage'),
    [appDirPath],
    dir('buildLinux'),
    { env: { ARCH: 'x86_64' } }
  )

export const packAppImageCopyDist = () =>
  fs.copyFile(
    dir('buildLinux', 'Koishi-x86_64.AppImage'),
    dir('dist', 'Koishi.AppImage')
  )

export const packAppImage = series(
  packAppImageMkdir,
  packAppImageCopyFiles,
  packAppImageGenerate,
  packAppImageCopyDist
)
