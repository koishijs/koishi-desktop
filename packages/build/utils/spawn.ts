import execa, { sync as execaSync } from 'execa'
import { Exceptions } from './exceptions'
import { dir } from './path'

export function spawnSyncOutput(
  command: string,
  args?: ReadonlyArray<string>,
  options?: execa.SyncOptions
): string {
  const parsedArgs = args ?? []
  const parsedOptions: execa.SyncOptions = Object.assign<
    execa.SyncOptions,
    execa.SyncOptions,
    execa.SyncOptions | undefined
  >({}, { stdio: 'pipe', shell: true }, options)
  const child = execaSync(command, parsedArgs, parsedOptions)
  return child.stdout.toString()
}

export async function spawnOutput(
  command: string,
  args?: ReadonlyArray<string>,
  options?: execa.SyncOptions
): Promise<string> {
  const parsedArgs = args ?? []
  const parsedOptions: execa.SyncOptions = Object.assign<
    execa.SyncOptions,
    execa.SyncOptions,
    execa.SyncOptions | undefined
  >({}, { stdio: 'pipe', shell: true }, options)
  const child = execa(command, parsedArgs, parsedOptions)
  let stdout = ''
  if (!child.stdout)
    throw Exceptions.runtime(
      `cannot get stdout of ${command} ${parsedArgs.join(' ')}`
    )
  child.stdout.on('data', (x) => (stdout += x))
  return new Promise<string>((resolve, reject) => {
    child.on('close', (x) => {
      if (x) reject(x)
      else resolve(stdout)
    })
  })
}

export async function spawnAsync(
  command: string,
  args?: ReadonlyArray<string>,
  options?: execa.SyncOptions
): Promise<number> {
  const parsedArgs = args ?? []
  const parsedOptions: execa.SyncOptions = Object.assign<
    execa.SyncOptions,
    execa.SyncOptions,
    execa.SyncOptions | undefined
  >({}, { stdio: 'inherit', shell: true }, options)
  const child = execa(command, parsedArgs, parsedOptions)
  return new Promise<number>((resolve) => {
    child.on('close', resolve)
  })
}

export async function exec(
  command: string,
  args?: ReadonlyArray<string>,
  cwd?: string,
  options?: execa.SyncOptions
): Promise<void> {
  const parsedArgs = (args ?? []).map((x) =>
    process.platform === 'win32' ? `"${x}"` : `'${x}'`
  )
  const parsedCwd = cwd ?? dir('root')
  const parsedOptions: execa.SyncOptions = Object.assign<
    execa.SyncOptions,
    execa.SyncOptions,
    execa.SyncOptions | undefined
  >({}, { stdio: 'inherit', shell: true, cwd: parsedCwd }, options)
  const child = execa(command, parsedArgs, parsedOptions)
  const result = await new Promise<number>((resolve) => {
    child.on('close', resolve)
  })
  if (result) {
    throw new Error(
      `'${child.spawnargs.join(' ')}' exited with error code: ${result}`
    )
  }
}
