#!/usr/bin/env python3
import sys
import argparse
import requests
from bs4 import BeautifulSoup
import json
import time
from discord_webhook import DiscordWebhook
import logging

# Set up logging
logging.basicConfig(filename='error.log', level=logging.DEBUG, format='%(asctime)s [%(levelname)s] %(message)s')

BASE_URL = "https://crt.sh/?q={}&output=json"
discord_webhook_url = "https://discord.com/api/webhooks/1052205480681951252/_hAFDr4MN8Z1iPsusHi0vFEb9Q_DtLAF-mnUGKO7ZSemNHP9OxjcN0i30gSVKZjdNPmb"
old_subdomains_file = "old_subdomains.txt"

def parser_error(errmsg):
    print("Usage: python3 " + sys.argv[0] + " [Options] use -h for help")
    print("Error: " + errmsg)
    sys.exit()

def parse_args():
    parser = argparse.ArgumentParser(epilog='\tExample: \r\npython3 ' + sys.argv[0] + " -l domains.txt")
    parser.error = parser_error
    parser._optionals.title = "OPTIONS"
    parser.add_argument('-l', '--domain_list', help='Specify a file containing a list of domains', required=True)
    parser.add_argument('-m', '--monitor', help='Monitor subdomains and send updates to Discord', action='store_true', required=False)
    return parser.parse_args()

def crtsh(domain):
    subdomains = set()
    wildcardsubdomains = set()

    try:
        response = requests.get(BASE_URL.format(domain), timeout=25)
        response.raise_for_status()  # Raise an HTTPError for bad responses
        content = response.content.decode('UTF-8')

        if content:
            soup = BeautifulSoup(content, 'html.parser')
            try:
                jsondata = json.loads(soup.text)
                for i in range(len(jsondata)):
                    name_value = jsondata[i].get('name_value', '')
                    if '\n' in name_value:
                        subname_value = name_value.split('\n')
                        for subname in subname_value:
                            if '*' in subname:
                                wildcardsubdomains.add(subname)
                            else:
                                subdomains.add(subname)
            except json.JSONDecodeError as e:
                logging.error(f"Error decoding JSON for domain {domain}: {e}")
                logging.error(f"Response content: {content}")

    except requests.exceptions.RequestException as e:
        logging.error(f"Error fetching subdomains for domain {domain}: {e}")

    return subdomains, wildcardsubdomains

def load_subdomains(file_path):
    try:
        with open(file_path, 'r') as file:
            subdomains = file.read().splitlines()
            return set(subdomains)
    except FileNotFoundError:
        return set()

def save_subdomains(file_path, subdomains):
    with open(file_path, 'w') as file:
        file.write('\n'.join(subdomains))

def send_to_discord(message):
    max_length = 2000
    chunks = [message[i:i+max_length] for i in range(0, len(message), max_length)]

    for chunk in chunks:
        webhook = DiscordWebhook(url=discord_webhook_url, content=chunk)
        webhook.execute()

if __name__ == "__main__":
    args = parse_args()

    if args.monitor:
        while True:
            try:
                with open(args.domain_list, 'r') as domains_file:
                    domains = domains_file.read().splitlines()

                all_subdomains = set()
                all_wildcardsubdomains = set()

                for domain in domains:
                    subdomains, wildcardsubdomains = crtsh(domain)
                    all_subdomains.update(subdomains)
                    all_wildcardsubdomains.update(wildcardsubdomains)

                # Load old subdomains
                old_subdomains = load_subdomains(old_subdomains_file)

                # Find new subdomains
                new_subdomains = all_subdomains - old_subdomains

                # Send new subdomains to Discord
                if new_subdomains:
                    message = f"New Subdomains found: {', '.join(new_subdomains)}"
                    send_to_discord(message)

                # Save current subdomains to file
                save_subdomains(old_subdomains_file, all_subdomains)

                # Wait for 10 hours before the next iteration
                time.sleep(10 * 60 * 60)

            except Exception as e:
                print(f"An error occurred: {e}")

    else:
        with open(args.domain_list, 'r') as domains_file:
            domains = domains_file.read().splitlines()

        for domain in domains:
            subdomains, wildcardsubdomains = crtsh(domain)
            print(f"Subdomains for {domain}:")
            for subdomain in subdomains:
                print(subdomain)
