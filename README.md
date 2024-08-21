# subfalcon

subfalcon is a subdomain enumeration tool that allows you to discover and monitor subdomains for a given list of domains or a single domain. It fetches subdomains from various sources, saves them to a SQLite database, and can notify updates via Discord.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Options](#options)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)
- [Disclaimer](#disclaimer)
- [ToDo](#todo)

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
- Option to process a single domain with the `-d` flag.

## Installation
You can install subfalcon using the following command: 
```bash
go install github.com/cyinnove/subfalcon/cmd/subfalcon@latest
```

## Usage

```bash
subfalcon -l domains.txt -w "YOUR_DISCORD_WEBHOOK_URL" -m
```

or

```bash
subfalcon -d example.com -w "YOUR_DISCORD_WEBHOOK_URL"
```

## Options

- `-l` or `--domain_list`: Specify a file containing a list of domains.
- `-m` or `--monitor`: Monitor subdomains and send updates to Discord.
- `-w` or `--webhook`: Specify the Discord webhook URL.
- `-d` or `--domain`: Specify a single domain for processing.

## Examples

- Basic usage with a list of domains:

  ```bash
  subfalcon -l domains.txt
  ```

- Process a single domain:

  ```bash
  subfalcon -d example.com
  ```

- Monitor a single domain and send updates to Discord:

  ```bash
  subfalcon -d example.com -m -w "YOUR_DISCORD_WEBHOOK_URL"
  ```

- Monitor and send updates to Discord with a list of domains:

  ```bash
  subfalcon -l domains.txt -m -w "YOUR_DISCORD_WEBHOOK_URL"
  ```

- Monitor new subdomains with a list of domains:
  
  ```bash
  subfalcon -l domains.txt -m
  ```

## Contributing

Feel free to contribute by opening issues or submitting pull requests.

## License

This project is licensed under the [MIT License](LICENSE).

## Disclaimer

Use this tool responsibly and only on systems you have permission to scan. The authors are not responsible for any misuse or damage caused by this tool.

## ToDo

- Add monitoring using Telegram.
- Add more resources to monitor subdomains.
- Add flags to customize monitoring time intervals.
- Add concurrency for domain processing.
- Improve the logging system.

> If you enjoy what we do, please support us:
> Buy Me Ko-fi! https://ko-fi.com/h0tak88r
