# Subdomain Monitor Script

A Python script for monitoring subdomains using the crt.sh certificate search.

## Features

- Fetches subdomains for a given list of domains from crt.sh
- Monitors subdomains and sends updates to Discord (optional)
- Logs errors to a file for debugging

## Prerequisites

- Python 3.x
- Install required Python packages:

  ```bash
  pip install requests beautifulsoup4 discord-webhook
  ```

## Usage

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/submonit88r.git
   ```

2. Navigate to the project directory:

   ```bash
   cd submonit88r
   ```

3. Create a virtual environment (optional but recommended):

   ```bash
   python -m venv venv
   source venv/bin/activate
   ```

4. Install dependencies:

   ```bash
   pip install -r requirements.txt
   ```

5. Run the script:

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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
