<h1 align="center">Reaper 💀</h1>
<p>
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

> A fast and simple dataset generator, using the power of search engines.

## Install

```sh
go get https://github.com/sno6/reaper
```

## Usage

```sh
Generate a dataset for a given topic

Usage:
  reaper [flags]

Flags:
  -h, --help               help for reaper
  -o, --out string         Output directory (default "./")
  -t, --term stringArray   Search terms
```

## Example

```sh
reaper -t "cats" -t "dogs" -o "./animals"

# After this completes you should have around ~2k images in your local animals/ folder
```

## Author

👤 **sno6**

* Website: https://sno6.github.io
* Github: [@sno6](https://github.com/sno6)