import { analyzeMetafile, build } from 'esbuild'
import fs from 'fs/promises'
import { info } from 'gulplog'
import mkdirp from 'mkdirp'
import { dir } from '../../utils/path'

const defineAgentMap = {
  win32: 'shellwin',
  darwin: 'shellmac',
  linux: 'shelllinux',
} as const

export const generateUserscript = async () => {
  const outfile = dir('buildResources', 'userscript.js')

  const result = await build({
    entryPoints: [dir('srcUserscript', 'src/index.ts')],
    outfile,

    define: {
      DEFINE_AGENT: `"${
        defineAgentMap[process.platform as keyof typeof defineAgentMap]
      }"`,
      DEFINE_SUPPORTS: 'KOISHELL_RUNTIME_SUPPORTS',
    },

    bundle: true,
    platform: 'browser',
    target: ['safari14'],
    minify: true,
    metafile: true,
    color: true,
  })

  info(await analyzeMetafile(result.metafile))

  if (process.platform === 'darwin') {
    await mkdirp(dir('srcShellMac', 'Sources/KoiShell/Resources'))

    fs.copyFile(
      outfile,
      dir('srcShellMac', 'Sources/KoiShell/Resources/userscript.js')
    )
  }
}
