import os
import requests

def load_config():
    host = os.environ.get("TP-LINK_HOST") 
    auth_password = os.environ.get("TP-LINK_PASSWORD")
    if host and auth_password:
        print("use os env config: host=%s, auth=%s".format(host, auth_password))
        return {
            "host": host,
            "auth_password": auth_password
        }
    else:
        from config import Config
        local_conf = Config
        print("use file config: host={0}, auth={1}".format(local_conf["host"], local_conf["auth_password"]))
        return local_conf

def login(host, auth_password):
    login_url = "http://{0}/".format(host)
    useragent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36"
    payload = {
        "method": "do",
        "login": {
            "password": auth_password
        }
    }


if __name__ == "__main__":
    conf = load_config()
    login(conf["host"], conf["auth_password"])