// Module "do" isn't here because it's treated as a vendor.
const modules = ['core', 'sdk', 'app']

export async function eachModule(
  fn:
    | ((module: string) => Promise<void>)
    | ((module: string) => () => Promise<void>)
    | ((module: string) => void)
) {
  for (const module of modules) {
    const fnResult = fn(module)

    if (!fnResult) return

    if (typeof fnResult === 'function') {
      await fnResult()
      return
    }

    await fnResult
  }
}

export async function tryEachModule(
  fn:
    | ((module: string) => Promise<void>)
    | ((module: string) => () => Promise<void>)
    | ((module: string) => void)
): Promise<boolean> {
  let failed = false

  for (const module of modules) {
    try {
      const fnResult = fn(module)

      if (!fnResult) continue

      if (typeof fnResult === 'function') {
        await fnResult()
        continue
      }

      await fnResult
    } catch (e) {
      failed = true
    }
  }

  return failed
}
