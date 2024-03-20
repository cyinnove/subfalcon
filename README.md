# subMonit88r

subMonit88r is a subdomain enumeration tool that allows you to discover and monitor subdomains for a given list of domains. It fetches subdomains from various sources, saves them to a SQLite database, and can notify updates via Discord.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Options](#options)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)
- [Disclaimer](#disclaimer)

## Features

- Subdomain enumeration from multiple sources.
    1. crt.sh
    2. hackertarget
    3. anubis
    4. Alienvault
    5. rapiddns
    6. urlscan.io
- SQLite database to store discovered subdomains.
- Discord integration for monitoring updates.
- Easy-to-use command-line interface.

## Installation
You can install subMonit88r using the following command: 
```bash
go install github.com/h0tak88r/subMonit88r@latest
```
## Usage

```bash
subMonit88r -l domains.txt -w "YOUR_DISCORD_WEBHOOK_URL" -m
```

## Options

- `-l` or `--domain_list`: Specify a file containing a list of domains.
- `-m` or `--monitor`: Monitor subdomains and send updates to Discord.
- `-w` or `--webhook`: Specify the Discord webhook URL.

## Examples

- Basic usage:

  ```bash
  subMonit88r -l domains.txt
  ```

- Monitor and send updates to Discord:
  ```bash
  subMonit88r -l domains.txt -m -w "YOUR_DISCORD_WEBHOOK_URL"
  ```

- Monitor New Subdomains 
```bash
subMonit88r -l domains.txt -m
```

## Contributing

Feel free to contribute by opening issues or submitting pull requests.

## License

This project is licensed under the [MIT License](LICENSE).

## Disclaimer

Use this tool responsibly and only on systems you have permission to scan. The authors are not responsible for any misuse or damage caused by this tool.
