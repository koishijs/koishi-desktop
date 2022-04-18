declare module 'gulplog' {
  /**
   * Highest log level. Typically used for debugging purposes.
   *
   * If the first argument is a string, all arguments are passed to node's util.format() before being emitted.
   * @param msg Message to log
   * @param args Arguments to format message with via util.format()
   */
  export function debug(msg: string, ...args: unknown[]): void
  export function debug(msg: unknown): void
  /**
   * Standard log level. Typically used for user information.
   *
   * If the first argument is a string, all arguments are passed to node's util.format() before being emitted.
   * @param msg Message to log
   * @param args Arguments to format message with via util.format()
   */
  export function info(msg: string, ...args: unknown[]): void
  export function info(msg: unknown): void
  /**
   * Warning log level. Typically used for warnings.
   *
   * If the first argument is a string, all arguments are passed to node's util.format() before being emitted.
   * @param msg Message to log
   * @param args Arguments to format message with via util.format()
   */
  export function warn(msg: string, ...args: unknown[]): void
  export function warn(msg: unknown): void
  /**
   * Error log level. Typically used when things went horribly wrong.
   *
   * If the first argument is a string, all arguments are passed to node's util.format() before being emitted.
   * @param msg Message to log
   * @param args Arguments to format message with via util.format()
   */
  export function error(msg: string, ...args: unknown[]): void
  export function error(msg: unknown): void
}
