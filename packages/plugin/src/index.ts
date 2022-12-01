import { Context, Schema } from 'koishi'
import { Enhance } from './plugins'

export const name = 'desktop'

export interface Config {
  enhance: {
    enabled: boolean
  }
}

export const Config: Schema<Config> = Schema.intersect([
  Schema.object({
    enhance: Schema.object({
      enabled: Schema.boolean()
        .default(true)
        .description('启用 Koishi 桌面增强。'),
    }),
  }).description('增强'),
])

export function apply(ctx: Context, config: Config) {
  if (config.enhance.enabled) ctx.plugin(Enhance)
}
