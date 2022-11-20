import fs from 'fs/promises'
import { eachModule } from '../utils/module'
import { dir } from '../utils/path'
import { exec } from '../utils/spawn'

export const buildDep = (pkg: string) => async () => {
  const dirGoMod = dir('root', `packages/${pkg}/go.mod`)
  await fs.writeFile(
    dirGoMod,
    (await fs.readFile(dirGoMod)).toString().split('\n').slice(0, 4).join('\n')
  )
  await exec('go', ['mod', 'tidy', '-e', '-v'], dir('root', `packages/${pkg}`))
}

export const dep = () => eachModule(buildDep)
