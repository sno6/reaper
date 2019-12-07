<h1 align="center">Reaper 💀</h1>
<p>
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

> A fast and simple dataset generator, using the power of search engines.

## Install

```sh
go get github.com/sno6/reaper
```

## Usage

```sh
Usage:
  reaper [flags]

Flags:
  -h, --help                 help for reaper
      --mh int               Only allow images width a height >= to this value (default -1)
      --mw int               Only allow images width a width >= to this value (default -1)
  -o, --out string           Output directory (default "./")
  -s, --search stringArray   Search terms
  -t, --timeout duration     HTTP request timeout in seconds (default 15s)
  -v, --verbose              Log errors
```

## Example

```sh
reaper -s "snoop dogg" -o "./snoop"

# After this completes you should have around ~2k images in your local snoop/ folder
```

## Search Engines

Currently `reaper` only supports the Duck Duck Go search engine however more will be added in the near future.

## Author

👤 **sno6**

* Website: https://sno6.github.io
* Github: [@sno6](https://github.com/sno6)