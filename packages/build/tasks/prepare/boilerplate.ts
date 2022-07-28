import { info } from 'gulplog'
import { join } from 'node:path'
import {
  repoBoilerplate,
  sourceGitHub,
  versionBoilerplate,
  versionNode,
} from '../../utils/config'
import { exists } from '../../utils/fs'
import { download } from '../../utils/net'
import { dir } from '../../utils/path'

export const prepareBoilerplate = async () => {
  const filename = 'boilerplate.zip'
  const platform = process.platform === 'win32' ? 'windows' : process.platform
  const src = `${sourceGitHub}/${repoBoilerplate}/releases/download/${versionBoilerplate}/boilerplate-${versionBoilerplate}-${platform}-amd64-node${
    versionNode.split('.')[0]
  }.zip`
  const dest = dir('buildCache')

  info('Checking temporary cache.')
  if (await exists(join(dest, filename))) return

  info('Now downloading boilerplate.')
  await download(src, dest, filename)
}
