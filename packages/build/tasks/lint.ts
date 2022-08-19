import { dir } from '../utils/path'
import { exec } from '../utils/spawn'

export const lint = () => exec('golangci-lint', ['run'], dir('root'))

export const lintFix = () =>
  exec('golangci-lint', ['run', '--fix'], dir('root'))
