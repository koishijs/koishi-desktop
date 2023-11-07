import type {} from '@koishijs/plugin-console'
import { Context } from 'koishi'
import { resolve } from 'path'

export const name = 'desktop-enhance'

export const inject = ['console']

export function apply(ctx: Context) {
  ctx.console.addEntry({
    dev: resolve(__dirname, '../../clients/enhance/client/index.ts'),
    prod: resolve(__dirname, '../../clients/enhance/dist'),
  })
}
