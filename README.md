# Pixler

Pixel art tools

Current features:
- Integer scaling
- Automatic downscaling by elimination of repeated lines of pixels

## Usage

`pixler <command> <options>`

### Commands

`upscale <factor> <image input> [image output]`

`downscale <factor> <image input> [image output]`

`unrepeat <image input> [image output]` - Automatically downscale image to 1:1 format by removing repeated lines of pixels. Warning: any intentional repeated lines will also be eliminated. Works only with non-interpolated images.
