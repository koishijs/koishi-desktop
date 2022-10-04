// Module "do" isn't here because it's treated as a vendor.
const modules = ['core', 'sdk', 'app', 'unfold']

export async function eachModule(
  fn:
    | ((module: string) => Promise<void>)
    | ((module: string) => () => Promise<void>)
    | ((module: string) => void)
): Promise<void> {
  for (const module of modules) {
    const fnResult = fn(module)

    if (!fnResult) continue

    if (typeof fnResult === 'function') {
      await fnResult()
      continue
    }

    await fnResult
  }
}

export async function tryEachModule(
  fn:
    | ((module: string) => Promise<void>)
    | ((module: string) => () => Promise<void>)
    | ((module: string) => void)
): Promise<void> {
  const errors: unknown[] = []

  for (const module of modules) {
    try {
      const fnResult = fn(module)

      if (!fnResult) continue

      if (typeof fnResult === 'function') {
        await fnResult()
        continue
      }

      await fnResult
    } catch (e: unknown) {
      errors.push(e)
    }
  }

  if (errors.length) throw errors
}
