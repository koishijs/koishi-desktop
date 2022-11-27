import { Context, Schema } from 'koishi'

export const name = 'desktop'

export interface Config {
  enhance: boolean
  telemetry: boolean
}

export const Config: Schema<Config> = Schema.intersect([
  Schema.object({
    enhance: Schema.boolean()
      .default(true)
      .description('启用 Koishi 桌面增强。'),
  }).description('增强'),
  Schema.object({
    telemetry: Schema.boolean()
      .default(true)
      .description(
        '启用此实例的 Koishi 桌面遥测。还需在 Koishi 桌面配置中启用遥测。'
      ),
  }).description('遥测'),
])

export function apply(ctx: Context) {
  // write your plugin here
}
