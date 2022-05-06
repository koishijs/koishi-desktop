import * as fs from 'fs'
import { getKoiSemVer, getKoiVersion } from './config'
import { resolve } from './path'

export async function genManifest() {
  const koiVersion = await getKoiVersion()
  const koiSemVer = await getKoiSemVer()

  const versionInfo = `
{
    "FixedFileInfo": {
        "FileVersion": {
            "Major": ${koiSemVer.major},
            "Minor": ${koiSemVer.minor},
            "Patch": ${koiSemVer.patch},
            "Build": ${koiSemVer.build}
        },
        "ProductVersion": {
            "Major": ${koiSemVer.major},
            "Minor": ${koiSemVer.minor},
            "Patch": ${koiSemVer.patch},
            "Build": ${koiSemVer.build}
        },
        "FileFlagsMask": "3f",
        "FileFlags ": "00",
        "FileOS": "040004",
        "FileType": "01",
        "FileSubType": "00"
    },
    "StringFileInfo": {
        "Comments": "koi",
        "CompanyName": "Koishi.js",
        "FileDescription": "koi",
        "FileVersion": "${koiVersion}",
        "InternalName": "koi",
        "LegalCopyright": "2022 Il Harper",
        "LegalTrademarks": "2022 Il Harper",
        "OriginalFilename": "koi",
        "PrivateBuild": "${koiVersion}",
        "ProductName": "koi",
        "ProductVersion": "${koiVersion}",
        "SpecialBuild": "${koiVersion}"
    },
    "VarFileInfo": {
        "Translation": {
            "LangID": "0409",
            "CharsetID": "04B0"
        }
    },
    "IconPath": "resources/koi.ico",
    "ManifestPath": "resources/koi.exe.manifest"
}
`.trim()

  const manifest = `
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <assemblyIdentity
    type="win32"
    name="koi"
    version="${koiSemVer.major}.${koiSemVer.minor}.${koiSemVer.patch}.${koiSemVer.build}"
    processorArchitecture="*"/>
 <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
   <security>
     <requestedPrivileges>
       <requestedExecutionLevel
         level="asInvoker"
         uiAccess="false"/>
       </requestedPrivileges>
   </security>
 </trustInfo>
</assembly>
  `.trim()

  await fs.promises.writeFile(resolve('versioninfo.json', 'koi'), versionInfo)
  await fs.promises.writeFile(
    resolve('resources/koi.exe.manifest', 'koi'),
    manifest
  )
}
