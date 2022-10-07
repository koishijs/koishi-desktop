import { parallel, series } from 'gulp'
import { Exceptions } from '../../utils/exceptions'
import { packAppImage } from './appimage'
import { packMac } from './mac'
import { packMsi } from './msi'
import { packPortable } from './portable'
import { packUnfold } from './unfold'

export * from './mac'
export * from './msi'
export * from './appimage'
export * from './portable'
export * from './unfold'

const buildPack = () => {
  switch (process.platform) {
    case 'win32':
      return parallel(packPortable, series(packUnfold, packMsi))
    case 'darwin':
      return parallel(packPortable, series(packUnfold, packMac))
    case 'linux':
      return parallel(packPortable, series(packUnfold, packAppImage))
    default:
      throw Exceptions.platformNotSupported()
  }
}

export const pack = buildPack()
