import json
import os
import time

import requests
import icmplib


def load_config():
    router_host = os.environ.get("TP_ROUTER_HOST")
    auth_password = os.environ.get("TP_ROUTER_PASSWORD")
    ipc_host = os.environ.get("TP_IPC_HOST")
    if router_host and auth_password and ipc_host:
        print(f"use os env config: router_host={router_host}, auth={auth_password}, ipc_host={ipc_host}")
        return {
            "router_host": router_host,
            "auth_password": auth_password,
            "ipc_host": ipc_host
        }
    else:
        from config import Config
        local_conf = Config
        print(
            f"use file config: router_host={local_conf['router_host']}, auth={local_conf['auth_password']}, ipc_host={ipc_host}")
        return local_conf


def login(host, auth_password):
    login_url = f"http://{host}/"
    user_agent = ("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) "
                  "Chrome/86.0.4240.75 Safari/537.36")
    payload = {
        "method": "do",
        "login": {
            "password": auth_password
        }
    }
    try:
        resp = requests.post(login_url, headers={
            'referer': login_url,
            'origin': login_url,
            'user-agent': user_agent,
            'Content-Type': 'application/json; charset=UTF-8',
        }, data=json.dumps(payload))
    except Exception as e:
        print(f"login request failed: {e}")
        return None
    try:
        error_code, stok = resp.json()['error_code'], resp.json()['stok']
        if error_code != 0:
            print(f'login resp throws error: code={error_code}')
            return None
        return stok
    except Exception as e:
        print(f"parse login exception error : {e}")
        return None


def modify_24g_wifi_channel(host, stok, channel):
    modify_url = f"http://{host}/stok={stok}/ds"
    referer = f"http://{host}/"
    user_agent = ("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) "
                  "Chrome/86.0.4240.75 Safari/537.36")
    payload = {
        "method": "set",
        "wireless": {
            "wlan_host_2g": {
                "enable": 1,
                "ssid": "XDR5480",
                "key": "ymbn680351",
                "auth": "0",
                "ssidbrd": 1,
                "encryption": 1,
                "channel": int(channel),
                "mode": 9,
                "bandwidth": 0,
                "twt": 0,
                "ofdma": 1
            },
            "wlan_bs": {
                "bs_enable": "0"
            }
        }
    }

    try:
        resp = requests.post(modify_url, headers={
            'referer': referer,
            'origin': referer,
            'user-agent': user_agent,
            'Content-Type': 'application/json; charset=UTF-8',
        }, data=json.dumps(payload))
    except Exception as e:
        print(f"modify wifi channel request failed: {e}")
        return None
    try:
        error_code, wait_time = resp.json()['error_code'], resp.json()['wait_time']
        if error_code != 0:
            print(f'modify wifi channel resp throws error: code={error_code}')
            return None
        return int(wait_time)
    except Exception as e:
        print(f"parse modify wifi channel exception error : {e}")
        return None


def is_ipc_alive(ipc_host):
    host = icmplib.ping(ipc_host, count=5, interval=1)
    return host.is_alive


if __name__ == "__main__":
    conf = load_config()
    if not is_ipc_alive(conf['ipc_host']):
        print(f'ipc server not alive, start modify wifi channel')
        stok = login(conf['router_host'], conf['auth_password'])
        print(f'get login stock success: stok={stok}')
        if stok:
            wait_time = modify_24g_wifi_channel(conf['router_host'], stok, 1)
            print(f'modify wifi channel to 1, wai_time={wait_time}')
            if wait_time:
                time.sleep(wait_time)
                modify_24g_wifi_channel(conf['router_host'], stok, 6)
                print(f'modify wifi channel to 6, finished')
