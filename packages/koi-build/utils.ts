import { SpawnOptions } from 'child_process'
import { spawn } from 'cross-spawn'
import * as fs from 'fs'

export async function exists(path: fs.PathLike): Promise<boolean> {
  try {
    await fs.promises.stat(path)
  } catch (_) {
    return false
  }
  return true
}

export async function notEmpty(path: fs.PathLike): Promise<boolean> {
  return Boolean((await fs.promises.readdir(path)).length)
}

export async function spawnOut(
  command: string,
  args?: ReadonlyArray<string>,
  options?: SpawnOptions
): Promise<string> {
  const parsedArgs = args ?? []
  const parsedOptions: SpawnOptions = Object.assign(
    {},
    { stdio: 'ignore' },
    options
  )
  const child = spawn(command, parsedArgs, parsedOptions)
  let stdout = ''
  child.stdout?.on('data', (x) => (stdout += x))
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
  const parsedOptions: SpawnOptions = Object.assign(
    {},
    { stdio: 'inherit' },
    options
  )
  const child = spawn(command, parsedArgs, parsedOptions)
  return new Promise<number>((resolve) => {
    child.on('close', resolve)
  })
}
