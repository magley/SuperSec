import sys
import argparse
import configparser
from loguru import logger

GLO = {
    "lesotho_url": "",
    "ip_address": "",
    "port": 0,
    "lesotho_api_client_name": "",
    "lesotho_api_key": "",
    "lesotho_https_cert": False,
}

def load_config():
    argparser = argparse.ArgumentParser()
    argparser.add_argument("--config", type=str, default="./config.ini", help="Config file path")
    args = argparser.parse_args()

    config = configparser.ConfigParser()
    config.read(args.config)

    global GLO

    lesotho_url_prefix = 'http://'
    if 'lesotho_https_cert' in config['MAIN']:
        GLO['lesotho_https_cert'] = config['MAIN']['lesotho_https_cert']
        lesotho_url_prefix = 'https://' 

    GLO['lesotho_url'] = f"{lesotho_url_prefix}{config['MAIN']['lesotho']}"
    GLO['ip_address'] = config['MAIN']['ip']
    GLO['port'] = config['MAIN']['port']
    GLO['lesotho_api_client_name'] = config['MAIN']['api_key_client_name']

    try:
        with open("apikey.secret") as f:
            GLO['lesotho_api_key'] = f.read()   
    except FileNotFoundError:
        logger.error("File 'apikey.secret' not found, please create the file and add a lesotho API key for client 'demo2_api' inside the file")
        sys.exit(1)