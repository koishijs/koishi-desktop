import { Context } from '@koishijs/client'
import './index.css'

declare global {
  interface Window {
    __KOI_SHELL__: {
      agent: string
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

const shellThemeMap = {
  light: 'TL',
  dark: 'TD',
  reset: 'TR',
} as const

const sendTheme = (theme: keyof typeof shellThemeMap) => {
  switch (window.__KOI_SHELL__?.agent) {
    case 'shellwin':
      window.chrome?.webview?.postMessage?.(shellThemeMap[theme])
      return

    case 'shellmac':
      window.webkit?.messageHandlers?.shellmacHandler?.postMessage?.(
        shellThemeMap[theme]
      )
      return

    case 'shelllinux':
      return

    default:
      return
  }
}

let themeObserver: MutationObserver

const enhance = () => {
  const agent = window.__KOI_SHELL__?.agent

  if (agent === 'shellwin' || agent === 'shellmac') {
    sendTheme(
      window.document.documentElement.classList.contains('dark')
        ? 'dark'
        : 'light'
    )

    themeObserver = new MutationObserver((mutations) => {
      for (const mutation of mutations) {
        if (mutation.attributeName === 'class')
          sendTheme(
            (mutation.target as HTMLElement).classList.contains('dark')
              ? 'dark'
              : 'light'
          )
      }
    })
    themeObserver.observe(window.document.documentElement, { attributes: true })
  }
}

const disposeEnhance = () => {
  const agent = window.__KOI_SHELL__?.agent

  if (agent === 'shellwin' || agent === 'shellmac') {
    sendTheme('reset')

    themeObserver.disconnect()
  }
}

export default (ctx: Context) => {
  enhance()
  const timer = setInterval(enhance, 4000)
  ctx.on('dispose', () => {
    clearInterval(timer)
    disposeEnhance()
  })
}
