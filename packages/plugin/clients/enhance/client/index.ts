import { Config, Context, Schema, useConfig } from '@koishijs/client'
import { RemovableRef } from '@vueuse/core'
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
*:not(
    input,
    textarea,
    .monaco-mouse-cursor-text,
    .monaco-mouse-cursor-text *,
    .k-text-selectable,
    .k-text-selectable *
  ) {
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

const enhanceTransCSS = `
.layout-container,
.main-container,
.layout-container .layout-aside,
.layout-status {
  background: transparent !important;
}

`

let themeObserver: MutationObserver
let styleSheet: HTMLStyleElement

const getComputedColorHex = (s: string) => {
  const r = colorString.get(
    window
      .getComputedStyle(window.document.documentElement)
      .getPropertyValue(s),
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

const syncStyleSheet = (config: RemovableRef<Config>) => {
  if (!styleSheet) return

  switch (config.value.desktop.enhance) {
    case 'enhanceColor':
      styleSheet.innerHTML = baseCSS + enhanceColorCSS
      break
    case 'enhance':
      styleSheet.innerHTML = baseCSS + enhanceCSS
      break
    case 'enhanceTrans':
      styleSheet.innerHTML = baseCSS + enhanceCSS + enhanceTransCSS
      break
  }
}

const syncTheme = (config: RemovableRef<Config>) => {
  switch (config.value.desktop.enhance) {
    case 'enhanceColor':
      send(
        `T${
          window.document.documentElement.classList.contains('dark') ? 'D' : 'L'
        }C${getComputedColorHex('--k-color-border')}${getComputedColorHex(
          '--bg1',
        )}${getComputedColorHex('--fg1')}`,
      )
      break
    case 'enhance':
    case 'enhanceTrans':
      send(
        window.document.documentElement.classList.contains('dark')
          ? 'TD'
          : 'TL',
      )
      break
  }
}

const resetTheme = () => send('TR')

const sync = (config: RemovableRef<Config>) => {
  syncStyleSheet(config)
  syncTheme(config)
}

const reset = () => {
  resetTheme()
}

const supports = (f: 'enhance' | 'enhanceColor') =>
  Array.isArray(window.__KOI_SHELL__?.supports) &&
  window.__KOI_SHELL__.supports.includes(f)

const enhance = (config: RemovableRef<Config>) => {
  if (!supports('enhance')) return

  if (!styleSheet) {
    styleSheet = window.document.getElementById(
      styleSheetId,
    ) as HTMLStyleElement
    styleSheet = document.createElement('style')
    styleSheet.id = styleSheetId
    document.head.appendChild(styleSheet)
  }

  if (!themeObserver) {
    themeObserver = new MutationObserver(() => sync(config))
    themeObserver.observe(window.document.documentElement, { attributes: true })
  }

  sync(config)
}

const disposeEnhance = () => {
  if (!supports('enhance')) return

  if (styleSheet) window.document.head.removeChild(styleSheet)
  if (themeObserver) themeObserver.disconnect()

  reset()
}

declare module '@koishijs/client' {
  interface Config {
    desktop: {
      enhance: 'off' | 'enhance' | 'enhanceColor' | 'enhanceTrans'
    }
  }
}

export default (ctx: Context) => {
  ctx.settings({
    id: 'desktop-enhance',
    schema: Schema.object({
      desktop: Schema.object({
        enhance: Schema.union(
          [
            Schema.const('off').description('增强关闭'),
            Schema.const('enhance')
              .description('增强')
              // @ts-expect-error 【管理员】孤梦星影 1:35:20 看了源码，实际上是可用的  只是没有类型  你可以 @ts-ignore
              .disabled(!supports('enhance')),

            Schema.const('enhanceColor')
              .description('增强色彩')
              // @ts-expect-error 【管理员】孤梦星影 1:35:20 看了源码，实际上是可用的  只是没有类型  你可以 @ts-ignore
              .disabled(!supports('enhanceColor')),

            Schema.const('enhanceTrans')
              .description('增强透视')
              // @ts-expect-error 【管理员】孤梦星影 1:35:20 看了源码，实际上是可用的  只是没有类型  你可以 @ts-ignore
              .disabled(!supports('enhance')),
          ].filter(Boolean),
        )
          .default(supports('enhance') ? 'enhance' : 'off')
          .description('Koishi 桌面增强模式。'),
      }).description('Koishi 桌面设置'),
    }),
  })

  ctx.on('ready', () => {
    const config = useConfig()

    if (config?.value?.desktop?.enhance !== 'off') {
      enhance(config)
      ctx.on('dispose', disposeEnhance)
    }
  })
}
