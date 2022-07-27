export const Exceptions = {
  platformNotSupported: () => new Error('Platform not supported.'),
  fileNotFound: (file: string) => new Error(`${file} not found.`),
}
