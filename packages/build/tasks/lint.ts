import { tryEachModule } from '../utils/module'
import { dir } from '../utils/path'
import { tryExec } from '../utils/spawn'

const buildLintArgs = (pkg: string, fix?: boolean) => {
  const isGh = process.env.GITHUB_ACTIONS
  const args = ['run']
  if (isGh) args.push('--out-format=github-actions')
  else args.push('--out-format=colored-line-number')
  args.push(`--path-prefix=${dir('packages', pkg)}`)
  if (fix) args.push('--fix')
  return args
}

export const lint = () =>
  tryEachModule((pkg) =>
    tryExec('golangci-lint', buildLintArgs(pkg), dir('packages', pkg))
  )

export const lintFix = () =>
  tryEachModule((pkg) =>
    tryExec('golangci-lint', buildLintArgs(pkg, true), dir('packages', pkg))
  )
