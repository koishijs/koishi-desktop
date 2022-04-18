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
