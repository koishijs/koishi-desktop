import { SpawnOptions } from 'child_process'
import { spawn } from 'cross-spawn'
import { Exceptions } from './exceptions'

export async function spawnOutput(
  command: string,
  args?: ReadonlyArray<string>,
  options?: SpawnOptions
): Promise<string> {
  const parsedArgs = args ?? []
  const parsedOptions: SpawnOptions = Object.assign<
    SpawnOptions,
    SpawnOptions,
    SpawnOptions | undefined
  >({}, { stdio: 'pipe', shell: true }, options)
  const child = spawn(command, parsedArgs, parsedOptions)
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
  options?: SpawnOptions
): Promise<number> {
  const parsedArgs = args ?? []
  const parsedOptions: SpawnOptions = Object.assign<
    SpawnOptions,
    SpawnOptions,
    SpawnOptions | undefined
  >({}, { stdio: 'inherit', shell: true }, options)
  const child = spawn(command, parsedArgs, parsedOptions)
  return new Promise<number>((resolve) => {
    child.on('close', resolve)
  })
}
