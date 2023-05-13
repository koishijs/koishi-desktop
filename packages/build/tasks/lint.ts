import { join } from 'node:path'
import { goEnv } from '../utils/config'
import { tryEachModule } from '../utils/module'
import { dir } from '../utils/path'
import { exec } from '../utils/spawn'

const buildLintArgs = (pkg: string, fix?: boolean) => {
  const isGh = process.env.GITHUB_ACTIONS
  const args = ['run']
  if (isGh) args.push('--out-format=github-actions')
  else args.push('--out-format=colored-line-number')
  args.push(`--path-prefix=${dir('packages', pkg)}`)
  if (fix) args.push('--fix')
  return args
}

const lintCommand = join(
  dir('buildCache', 'golangci-lint'),
  goEnv.GOOS === 'windows' ? 'golangci-lint.exe' : 'golangci-lint'
)

export const lint = () =>
  tryEachModule((pkg) =>
    exec(lintCommand, buildLintArgs(pkg), dir('packages', pkg))
  )

export const lintFix = () =>
  tryEachModule((pkg) =>
    exec(lintCommand, buildLintArgs(pkg, true), dir('packages', pkg))
  )
