import { dir } from '../utils/path'
import { exec2 } from '../utils/spawn'

export const startApp = () => exec2('koi', [], dir('buildPortable'))
