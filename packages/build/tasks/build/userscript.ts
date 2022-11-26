import { analyzeMetafile, build } from 'esbuild'
import { info } from 'gulplog'
import { dir } from '../../utils/path'

const defineAgentMap = {
  win32: 'shellwin',
  darwin: 'shellmac',
  linux: 'shelllinux',
} as const

export const generateUserscript = async () => {
  const result = await build({
    entryPoints: [dir('srcUserscript', 'src/index.ts')],
    outfile: dir('buildResources', 'userscript.js'),

    define: {
      DEFINE_AGENT:
        defineAgentMap[process.platform as keyof typeof defineAgentMap],
    },

    bundle: true,
    platform: 'browser',
    target: ['safari14'],
    minify: true,
    metafile: true,
    color: true,
  })

  info(await analyzeMetafile(result.metafile))
}
