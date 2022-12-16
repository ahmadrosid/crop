# crop

```bash
crop - a simple CLI to crop images

Usage:
  crop [flags]

Flags:
  -h, --help                help for crop
      --out-folder string   Output folder name. (default "output")
      --padding ints        Set padding left and top (default [0,0])
      --size ints           Set image size: width x height (default [0,0])
```

## Installation

```bash
go install github.com/ahmadrosid/crop@latest
```

## Example

```bash
crop original.png --size=512,512 --padding=80,280
```
