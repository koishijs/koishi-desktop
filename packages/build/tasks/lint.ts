import { tryEachModule } from '../utils/module'
import { dir } from '../utils/path'
import { exec } from '../utils/spawn'

export const lint = () =>
  tryEachModule((pkg) => exec('golangci-lint', ['run'], dir('packages', pkg)))

export const lintFix = () =>
  tryEachModule((pkg) =>
    exec('golangci-lint', ['run', '--fix'], dir('packages', pkg))
  )
