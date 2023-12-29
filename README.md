# Subdomain Monitor Script

This Python script allows you to monitor subdomains using the various passive resources. It fetches subdomains for a given list of domains from various sources, including crt.sh, hackertarget, anubis, Alienvault, rapiddns, and urlscan.io. Optionally, it can monitor subdomains and send updates to Discord.

## Features

- Fetches subdomains for a given list of domains from:
    1. crt.sh
    2. hackertarget
    3. anubis
    4. Alienvault
    5. rapiddns
    6. urlscan.io
- Monitors subdomains and sends updates to Discord (optional)

## Prerequisites

- Python 3.x
- Install required Python packages:

  ```bash
  pip install requests beautifulsoup4 discord-webhook
  ```

## Usage

1. **Clone the repository:**

   ```bash
   git clone https://github.com/yourusername/submonit88r.git
   ```

2. **Navigate to the project directory:**

   ```bash
   cd submonit88r
   ```

3. **Create a virtual environment (optional but recommended):**

   ```bash
   python -m venv venv
   source venv/bin/activate
   ```

4. **Install dependencies:**

   ```bash
   pip install -r requirements.txt
   ```

5. **Run the script:**

   - To fetch subdomains:

     ```bash
     python submonit88r.py -l domains.txt
     ```

   - To monitor and send updates to Discord, provide the webhook URL:

     ```bash
     python submonit88r.py -l domains.txt -m -w YOUR_DISCORD_WEBHOOK_URL
     ```

## Options

- `-l` or `--domain_list`: Specify a file containing a list of domains.
- `-m` or `--monitor`: Monitor subdomains and send updates to Discord (optional).
- `-w` or `--webhook_url`: Discord Webhook URL for sending updates (required when using `-m`).

For any issues or improvements, please [create an issue](https://github.com/h0tak88r/submonit88r/issues), or open pull request.
