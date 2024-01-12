# Submonit88r

Submonit88r is a subdomain enumeration tool that allows you to discover and monitor subdomains for a given list of domains. It fetches subdomains from various sources, saves them to a SQLite database, and can notify updates via Discord.

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

1. Clone the repository:

   ```bash
   git clone https://github.com/h0tak88r/submonit88r.git
   ```

2. Navigate to the project directory:

   ```bash
   cd submonit88r
   ```

3. Install dependencies:

   ```bash
   pip install -r requirements.txt
   ```

## Usage

```bash
python submonit88r.py -l domains.txt -w "YOUR_DISCORD_WEBHOOK_URL" -m
```

## Options

- `-l` or `--domain_list`: Specify a file containing a list of domains.
- `-m` or `--monitor`: Monitor subdomains and send updates to Discord.
- `-w` or `--webhook`: Specify the Discord webhook URL.

## Examples

- Basic usage:

  ```bash
  python submonit88r.py -l domains.txt
  ```

- Monitor and send updates to Discord:

  ```bash
  python submonit88r.py -l domains.txt -m -w "YOUR_DISCORD_WEBHOOK_URL"
  ```
- Run in virtual environment\
	```python
	python3 -m venv venv
	source venv/bin/activate
	```
## Contributing

Feel free to contribute by opening issues or submitting pull requests. Please follow the [Contributing Guidelines](CONTRIBUTING.md).

## License

This project is licensed under the [MIT License](LICENSE).

## Disclaimer

Use this tool responsibly and only on systems you have permission to scan. The authors are not responsible for any misuse or damage caused by this tool.
