import { Context } from '@koishijs/client'

const enhance = () => {
  //
}

const disposeEnhance = () => {
  //
}

export default (ctx: Context) => {
  const timer = setInterval(enhance, 4000)
  ctx.on('dispose', () => {
    clearInterval(timer)
    disposeEnhance()
  })
}
