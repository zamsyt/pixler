# Pixler

Pixel art tools

Current features:
- Integer scaling
- Automatic downscaling by elimination of repeated lines of pixels
- Check images against a palette

## Usage

`pixler <command> <options>`

### Commands

`upscale <factor> <image input> [image output]`

`downscale <factor> <image input> [image output]`

`unrepeat <image input> [image output]` - Automatically downscale image to 1:1 format by removing repeated lines of pixels. Warning: any intentional repeated lines will also be eliminated. Works only on images scaled without interpolation.

`palette <image input> [diff output] [palette]` - Gives a count of pixels that aren't in the palette, and produces an image of them
