declare global {
  interface Window {
    __KOI_SHELL__: unknown
  }

  const DEFINE_AGENT: string
  const DEFINE_SUPPORTS: string[]
}

interface KoiShell {
  agent?: string
  supports?: string[]
}

function ensureKoiShell() {
  if (!window.__KOI_SHELL__ || typeof window.__KOI_SHELL__ !== 'object')
    window.__KOI_SHELL__ = {}

  const koiShell = window.__KOI_SHELL__ as KoiShell
  koiShell.agent = DEFINE_AGENT
  koiShell.supports = DEFINE_SUPPORTS
}

setInterval(ensureKoiShell, 10000)
ensureKoiShell()

export {}
