import { parallel } from 'gulp'
import { Exceptions } from '../../utils/exceptions'
import { packMacApp } from './mac'
import { packMsi } from './msi'
import { packPortable } from './portable'

export * from './mac'
export * from './msi'
export * from './portable'

const buildPack = () => {
  switch (process.platform) {
    case 'win32':
      return parallel(packPortable, packMsi)
    case 'darwin':
      return parallel(packPortable, packMacApp)
    case 'linux':
      return parallel(packPortable)
    default:
      throw Exceptions.platformNotSupported()
  }
}

export const pack = buildPack()
