import { eachModule } from '../utils/module'
import { dir } from '../utils/path'
import { exec } from '../utils/spawn'

export const lint = () =>
  eachModule((pkg) =>
    exec('golangci-lint', ['run', `packages/${pkg}/...`], dir('root'))
  )

export const lintFix = () =>
  eachModule((pkg) =>
    exec('golangci-lint', ['run', `packages/${pkg}/...`, '--fix'], dir('root'))
  )
