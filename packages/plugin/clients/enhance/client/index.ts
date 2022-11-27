import { Context } from '@koishijs/client'

const styleSheetId = 'koishell-enhance-stylesheet'

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

const enhanceCSS = `
body, nav.layout-activity {
  background: transparent !important;
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

const enhance = () => {
  const agent = window.__KOI_SHELL__?.agent

  if (agent === 'shellwin' || agent === 'shellmac') {
    sendTheme(
      window.document.documentElement.classList.contains('dark')
        ? 'dark'
        : 'light'
    )

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
  sendTheme('reset')

  const styleSheet = window.document.getElementById(styleSheetId)
  if (styleSheet) window.document.head.removeChild(styleSheet)
}

export default (ctx: Context) => {
  enhance()
  const timer = setInterval(enhance, 4000)
  ctx.on('dispose', () => {
    clearInterval(timer)
    disposeEnhance()
  })
}
