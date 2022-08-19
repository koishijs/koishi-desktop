import { eachModule } from '../utils/module'
import { dir } from '../utils/path'
import { exec } from '../utils/spawn'

export const lint = () =>
  eachModule((pkg) => exec('golangci-lint', ['run'], dir('packages', pkg)))

export const lintFix = () =>
  eachModule((pkg) => exec('golangci-lint', ['run', '--fix'], dir('root', pkg)))
