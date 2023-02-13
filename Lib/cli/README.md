# Unlock Music Project - CLI Edition

Original: Web Edition https://git.unlock-music.dev/um/web

- [![Build Status](https://ci.unlock-music.dev/api/badges/um/cli/status.svg)](https://ci.unlock-music.dev/um/cli)
- [Release Download](https://git.unlock-music.dev/um/cli/releases/latest)
- [Latest Build](https://git.unlock-music.dev/um/-/packages/generic/cli-build/)

## Features

- [x] All Algorithm Supported By `unlock-music/web`
- [x] Complete Metadata & Cover Image

## Hou to Build

- Requirements: **Golang 1.19**

1. run `go install unlock-music.dev/cli/cmd/um@master`

## How to use

- Drag the encrypted file to `um.exe` (Tested on Windows)
- Run: `./um [-o <output dir>] [-i] <input dir/file>`
- Use `./um -h` to show help menu
