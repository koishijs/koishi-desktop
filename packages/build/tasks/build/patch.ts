import * as fs from 'node:fs'
import { koishiManifest } from '../../templates'
import { koishiVersionStrings } from '../../utils/config'
import { dir } from '../../utils/path'
import { exec } from '../../utils/spawn'

export const patchNodeRcedit = async () => {
  const koishiManifestPath = dir('buildCache', 'koishi.exe.manifest')
  const exePath = dir('buildPortableData', 'node/koishi.exe')

  await fs.promises.writeFile(koishiManifestPath, koishiManifest)

  const args = [
    exePath,
    '--set-icon',
    dir('buildAssets', 'koishi.ico'),
    '--application-manifest',
    koishiManifestPath,
  ]

  ;(
    Object.keys(koishiVersionStrings) as (keyof typeof koishiVersionStrings)[]
  ).forEach((x) => {
    args.push('--set-version-string', x, koishiVersionStrings[x])
  })

  await exec('rcedit.exe', args, dir('buildCache'))

  // Change subsystem to GUI
  // https://learn.microsoft.com/windows/win32/api/winnt/ns-winnt-image_optional_header64
  const buf = await fs.promises.readFile(exePath)
  const peOffset = buf.readUint32LE(0x3c) // IMAGE_DOS_HEADER.e_lfanew, the offset of PE Header
  const subsystemOffset = peOffset + 0x5c // Offset of IMAGE_OPTIONAL_HEADER64.Subsystem
  buf.writeUInt8(2, subsystemOffset) // IMAGE_SUBSYSTEM_WINDOWS_GUI
  await fs.promises.writeFile(exePath, buf)
}

export const patch =
  process.platform === 'win32'
    ? patchNodeRcedit
    : async () => {
        // Ignore
      }
