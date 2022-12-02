const sleepIntl = (ms: number): Promise<void> =>
  new Promise((resolve) => setTimeout(resolve, ms))

export const sleep = () => sleepIntl(2000)
