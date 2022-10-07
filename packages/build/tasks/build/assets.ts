import { Icns, IcnsImage } from '@fiahfy/icns'
import { parallel, series } from 'gulp'
import { info } from 'gulplog'
import * as fs from 'node:fs'
import sharp from 'sharp'
import { exists } from '../../utils/fs'
import { dir } from '../../utils/path'

const builtFilename = dir('buildAssets', '.built')

const icnsTypes = [
  {
    size: 16,
    ostypes: ['icp4'],
  },
  {
    size: 32,
    ostypes: ['icp5', 'ic11'],
  },
  {
    size: 64,
    ostypes: ['icp6', 'ic12'],
  },
  {
    size: 128,
    ostypes: ['ic07'],
  },
  {
    size: 256,
    ostypes: ['ic08', 'ic13'],
  },
  {
    size: 512,
    ostypes: ['ic09', 'ic14'],
  },
  {
    size: 1024,
    ostypes: ['ic10'],
  },
] as const

const buildIcons = () =>
  fs
    .readdirSync(dir('srcAssets'))
    .filter((x) => x.endsWith('.svg'))
    .map((x) => x.slice(0, x.length - 4))
    .filter((x) => !iconsBlacklist.includes(x))

const iconsBlacklist = ['koishi', 'koishi-template']
const icons = buildIcons()
const iconsWithKoishi = [...icons, ...iconsBlacklist]

const sizes = [256, 64, 48, 40, 32, 24, 20, 16]

const getRowStride = (width: number): number => {
  if (width % 32 === 0) {
    return width / 8
  } else {
    return 4 * (Math.floor(width / 32) + 1)
  }
}

const getPixelColor = (
  size: number,
  image: Buffer,
  x: number,
  y: number
): number => {
  let xi = x < 0 ? 0 : x
  let yi = y < 0 ? 0 : y

  if (x >= size) xi = size - 1
  if (y >= size) yi = size - 1

  const i =
    xi < 0 || xi >= size || yi < 0 || yi >= size ? -1 : (size * yi + xi) << 2
  return image.readUInt32BE(i)
}

const generateBmp = async (size: number, image: sharp.Sharp) => {
  const data = await image.raw().toBuffer()

  const headerbuf = Buffer.alloc(40)
  headerbuf.writeUInt32LE(40, 0) // Size of this header
  headerbuf.writeInt32LE(size, 4) // Width
  headerbuf.writeInt32LE(size * 2, 8) // Height (XOR + AND)
  headerbuf.writeUInt16LE(1, 12) // Color plane
  headerbuf.writeUInt16LE(32, 14) // Bits per pixel (bpp)
  headerbuf.writeUInt32LE(0, 16) // Compression method
  headerbuf.writeUInt32LE(0, 20) // Image size
  headerbuf.writeInt32LE(0, 24) // Horizontal resolution
  headerbuf.writeInt32LE(0, 28) // Vertical resolution
  headerbuf.writeUInt32LE(0, 32) // No color palette
  headerbuf.writeUInt32LE(0, 36) // No important colors

  const len = data.length
  const buf = Buffer.alloc(len + getRowStride(size) * size)

  // XOR map
  for (let y = 0; y < size; y++) {
    for (let x = 0; x < size; x++) {
      const pxColor = getPixelColor(size, data, x, y)

      const r = (pxColor >> 24) & 255
      const g = (pxColor >> 16) & 255
      const b = (pxColor >> 8) & 255
      const a = pxColor & 255
      const newColor = b | (g << 8) | (r << 16) | (a << 24)

      const pos = ((size - y - 1) * size + x) * 4
      buf.writeInt32LE(newColor, pos)
    }
  }

  // AND map (32 bits per line)
  for (let y = 0; y < size; y++) {
    for (let x = 0; x < size; x++) {
      const pxColor = getPixelColor(size, data, x, y)
      const alpha = (pxColor & 255) > 0 ? 0 : 1
      const bitNum = (size - y - 1) * size + x
      // width per line in multiples of 32 bits
      const width32 =
        size % 32 === 0 ? Math.floor(size / 32) : Math.floor(size / 32) + 1

      const line = Math.floor(bitNum / size)
      const offset = Math.floor(bitNum % size)
      const bitVal = alpha & 0x00000001

      const pos = len + line * width32 * 4 + Math.floor(offset / 8)
      const newVal = buf.readUInt8(pos) | (bitVal << (7 - (offset % 8)))
      buf.writeUInt8(newVal, pos)
    }
  }

  return Buffer.concat([headerbuf, buf])
}

const generateIco = async (icon: string, image: sharp.Sharp) => {
  const images = await Promise.all(
    sizes.map(async (size) => {
      return {
        size,
        buffer:
          size === 256
            ? await image.resize(size).toBuffer() // use png in 256
            : await generateBmp(size, image.resize(size)), // and bmp in other
      }
    })
  )

  const len =
    6 + // File header
    16 * images.length + // Image header
    images.map((x) => x.buffer.length).reduce((x, y) => x + y)
  let pos = 0
  let imagePos = 6 + 16 * images.length

  // ICO file header (ICONDIR)
  const buf = Buffer.allocUnsafe(len)
  pos = buf.writeUint16LE(0, pos) // Reserved
  pos = buf.writeUInt16LE(1, pos) // Type: ICO
  pos = buf.writeUInt16LE(images.length, pos) // Number of images

  for (const im of images) {
    // ICO image item header (ICONDIRENTRY)
    pos = buf.writeUInt8(im.size === 256 ? 0 : im.size, pos)
    pos = buf.writeUInt8(im.size === 256 ? 0 : im.size, pos)
    pos = buf.writeUInt8(0, pos) // No color palette
    pos = buf.writeUInt8(0, pos) // Reserved
    pos = buf.writeUInt16LE(1, pos) // Color plane
    pos = buf.writeUInt16LE(32, pos) // Bits per pixel (bpp)
    pos = buf.writeUInt32LE(im.buffer.length, pos) // Image (bmp) file size
    pos = buf.writeUInt32LE(imagePos, pos) // Image file offset

    im.buffer.copy(buf, imagePos)
    imagePos += im.buffer.length
  }

  await fs.promises.writeFile(dir('buildAssets', `${icon}.ico`), buf)
}

const generateIcns = async (icon: string, image: sharp.Sharp) => {
  const icns = new Icns()
  for (const icnsType of icnsTypes) {
    const png = await image.clone().resize(icnsType.size).toBuffer()
    for (const ostype of icnsType.ostypes)
      icns.append(IcnsImage.fromPNG(png, ostype))
  }

  await fs.promises.writeFile(dir('buildAssets', `${icon}.icns`), icns.data)
}

export const generateAssetsImage = async () => {
  info('Checking assets cache.')
  if (await exists(builtFilename)) {
    info('Assets were successfully built last time, skip rebuilding.')
    return
  }

  await Promise.all(
    iconsWithKoishi.map((icon) => {
      const image = sharp(dir('srcAssets', `${icon}.svg`)).png()

      return Promise.all([
        generateIco(icon, image.clone()),
        generateIcns(icon, image.clone()),
        image.clone().toFile(dir('buildAssets', `${icon}.png`)),
      ])
    })
  )

  await fs.promises.writeFile(builtFilename, '')
}

const generateAssetsCodeIntl = async (
  icon: string,
  os: string,
  src: string
) => {
  const filename = `${icon}_${os}.go`
  const buf = await fs.promises.readFile(src)

  const name = icon
    .split('-')
    .map((x) => `${x[0].toUpperCase()}${x.slice(1)}`)
    .join('')
  const data = Array.from(buf)
    .map((x) => x.toString(16))
    .map((x) => (x.length === 1 ? `0x0${x}` : `0x${x}`))
    .reduce((a, b) => `${a}, ${b}`)

  await fs.promises.writeFile(
    dir('srcIcon', filename),
    `package icon\n\nvar ${name} = []byte{${data}}`
  )
}

export const generateAssetsCode = async () => {
  const osList = ['windows', 'darwin', 'linux'] as const

  for (const icon of icons) {
    for (const os of osList) {
      await generateAssetsCodeIntl(
        icon,
        os,
        dir('buildAssets', `${icon}${os === 'windows' ? '.ico' : '.png'}`)
      )
    }
  }

  await generateAssetsCodeIntl(
    'koishi',
    'windows',
    dir('buildAssets', 'koishi.ico')
  )
  await generateAssetsCodeIntl(
    'koishi',
    'darwin',
    dir('buildAssets', 'koishi-template.png')
  )
  await generateAssetsCodeIntl(
    'koishi',
    'linux',
    dir('buildAssets', 'koishi.png')
  )
}

export const generateAssetsCopySvg = () =>
  fs.promises.cp(dir('srcAssets'), dir('buildAssets'), { recursive: true })

export const generateAssets = parallel(
  generateAssetsCopySvg,
  series(generateAssetsImage, generateAssetsCode)
)
