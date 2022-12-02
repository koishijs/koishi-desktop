import { Context } from '@koishijs/client'

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

const styleSheetId = 'koishell-enhance-stylesheet'

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

input,
textarea {
  -webkit-touch-callout: auto !important;
  user-select: auto !important;
  -webkit-user-select: auto !important;
  cursor: auto !important;
}

* {
  -webkit-touch-callout: none !important;
  user-select: none !important;
  -webkit-user-select: none !important;
  cursor: default !important;
}
`

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

    let styleSheet = window.document.getElementById(
      styleSheetId
    ) as HTMLStyleElement
    if (!styleSheet) {
      styleSheet = document.createElement('style')
      styleSheet.id = styleSheetId
      styleSheet.innerHTML = enhanceCSS
      document.head.appendChild(styleSheet)
    }
  }
}

const disposeEnhance = () => {
  const agent = window.__KOI_SHELL__?.agent

  if (agent === 'shellwin' || agent === 'shellmac') {
    sendTheme('reset')

    themeObserver.disconnect()

    const styleSheet = window.document.getElementById(styleSheetId)
    if (styleSheet) window.document.head.removeChild(styleSheet)
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
