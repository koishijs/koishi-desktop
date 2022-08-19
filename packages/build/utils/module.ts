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
