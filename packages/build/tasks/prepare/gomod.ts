import { eachModule } from '../../utils/module'
import { dir } from '../../utils/path'
import { exec2 } from '../../utils/spawn'

export const buildPrepareGoMod = (pkg: string) => async () => {
  await exec2('go', ['mod', 'download'], dir('root', `packages/${pkg}`))
}

/**
 * Prepare go modules.
 *
 * Use series instead of parallel to maximize caches.
 */
export const prepareGoMod = () => eachModule(buildPrepareGoMod)
