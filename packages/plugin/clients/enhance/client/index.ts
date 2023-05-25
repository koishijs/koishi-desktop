import { Context } from '@koishijs/client'
import * as colorString from 'color-string'

declare global {
  interface Window {
    __KOI_SHELL__?: {
      agent?: string
      supports?: string[]
    }

    chrome: {
      webview: {
        postMessage: (message: string) => void
      }
    }

    webkit: {
      messageHandlers: {
        shellmacHandler: {
          postMessage: (message: string) => void
        }
      }
    }
  }
}

const styleSheetId = 'koishell-enhance-stylesheet'

const baseCSS = `
input,
textarea {
  -webkit-touch-callout: auto !important;
  user-select: auto !important;
  -webkit-user-select: auto !important;
  cursor: auto !important;
}

*:not(input, textarea, .monaco-mouse-cursor-text, .monaco-mouse-cursor-text *) {
  -webkit-touch-callout: none !important;
  user-select: none !important;
  -webkit-user-select: none !important;
  cursor: default !important;
}

`

const enhanceCSS = `
body,
nav.layout-activity {
  background: transparent !important;
}

@media screen and (min-width: 768px) {
  div.layout-container {
    clip-path: inset(0 0 round 12px 0 0 0) !important;
  }
}

nav.layout-activity {
  border: 0 !important;
}

`

const enhanceColorCSS = `
.layout-activity,
.layout-container {
  border-top: var(--k-color-divider-dark) 1px solid;
}

`

let themeObserver: MutationObserver
let styleSheet: HTMLStyleElement

const getComputedColorHex = (s: string) => {
  const r = colorString.get(
    window.getComputedStyle(window.document.documentElement).getPropertyValue(s)
  )
  if (!r || !r.value) return '000000'
  return colorString.to.hex(r.value).slice(1, 7)
}

const send = (message: string) => {
  switch (window.__KOI_SHELL__?.agent) {
    case 'shellwin':
      window.chrome?.webview?.postMessage?.(message)
      return

    case 'shellmac':
      window.webkit?.messageHandlers?.shellmacHandler?.postMessage?.(message)
      return

    case 'shelllinux':
      return

    default:
      return
  }
}

const syncStyleSheet = () => {
  if (!styleSheet) return

  // TODO: Get config
  switch ('enhanceColor' as 'enhanceColor' | 'enhance') {
    case 'enhanceColor':
      styleSheet.innerHTML = baseCSS + enhanceColorCSS
      break
    case 'enhance':
      styleSheet.innerHTML = baseCSS + enhanceCSS
      break
  }
}

const syncTheme = () => {
  // TODO: Get config
  switch ('enhanceColor' as 'enhanceColor' | 'enhance') {
    case 'enhanceColor':
      send(
        `T${
          window.document.documentElement.classList.contains('dark') ? 'D' : 'L'
        }C${getComputedColorHex('--k-color-border')}${getComputedColorHex(
          '--bg1'
        )}${getComputedColorHex('--fg1')}`
      )
      break
    case 'enhance':
      send(
        window.document.documentElement.classList.contains('dark') ? 'TD' : 'TL'
      )
      break
  }
}

const resetTheme = () => send('TR')

const sync = () => {
  syncStyleSheet()
  syncTheme()
}

const reset = () => {
  resetTheme()
}

const supportsEnhance = () =>
  Array.isArray(window.__KOI_SHELL__?.supports) &&
  window.__KOI_SHELL__.supports.includes('enhance')

const enhance = () => {
  if (!supportsEnhance()) return

  if (!styleSheet) {
    styleSheet = window.document.getElementById(
      styleSheetId
    ) as HTMLStyleElement
    styleSheet = document.createElement('style')
    styleSheet.id = styleSheetId
    document.head.appendChild(styleSheet)
  }

  if (!themeObserver) {
    themeObserver = new MutationObserver(sync)
    themeObserver.observe(window.document.documentElement, { attributes: true })
  }

  sync()
}

const disposeEnhance = () => {
  if (!supportsEnhance()) return

  if (styleSheet) window.document.head.removeChild(styleSheet)
  if (themeObserver) themeObserver.disconnect()

  reset()
}

export default (ctx: Context) => {
  enhance()
  ctx.on('dispose', disposeEnhance)
}
