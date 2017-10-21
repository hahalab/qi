import logging
import os
import shlex
import socket
import subprocess
import requests
import urlparse
import time
import yaml


def handler(event, context):
    logger = logging.getLogger()
    config = read_config()
    port = pick_port()
    p = start_command(config.command, port)

    headers = event.headers
    method = event.method
    data = event.data
    path = event.path

    logger.info(p.stdout)
    logger.error(p.stderr)
    wait_for_listen(port)
    response = send_request(port=port, headers=headers, method=method, data=data, path=path)
    return dict(headers=response.headers, data=response.data, code=response.code)


def wait_for_listen(port):
    while is_listen(port):
        time.sleep(0.1)


def read_config():
    return yaml.load(file('ha.yml', 'r'))


def start_command(cmd, port):
    my_env = os.environ.copy()
    my_env["PORT"] = port
    args = shlex.split('sh -c "{}"'.format(cmd))
    return subprocess.Popen(args, env=my_env)


def pick_port():
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.bind(('', 0))
    port = sock.getsockname()[1]
    sock.close()
    return port


def is_listen(port):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        sock.connect(("127.0.0.1", int(port)))
    except socket.error as msg:
        return False
    finally:
        sock.close()
    return True


def send_request(port, headers, method, data, path):
    target = "http://127.0.0.1:{}".format(port)
    url = urlparse.urljoin(target, path)
    return getattr(requests, method.lower())(url, headers=headers, data=data)
