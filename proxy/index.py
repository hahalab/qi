import json
import logging
import os
import requests
import shlex
import socket
import subprocess
import sys
import time
import urlparse


def handler(event, context):
    logger = logging.getLogger()
    config = read_config()
    port = pick_port()
    p = start_command(config["command"], port)

    event = json.loads(event)

    logger.info("receive event: {} context: {}".format(event, context))

    headers = getattr(event, "headers", {})

    for k, v in headers:
        headers[k] = ";".join(v)

    event["headers"] = headers

    method = event["method"]
    data = getattr(event, "data", None)
    path = getattr(event, "path", "/")

    wait_for_listen(port)
    response = send_request(port=port, headers=headers, method=method, data=data, path=path)
    p.terminate()
    # logger.info(p.stdout)
    # logger.error(p.stderr)
    headers = response.headers
    for k, v in headers:
        headers[k] = v.split(';')
    return dict(headers=headers, data=response.content, code=response.status_code)


def wait_for_listen(port):
    while not is_listen(port):
        time.sleep(0.1)


def read_config():
    with open('ha.yml', 'r') as f:
        config = dict()
        for l in f.readlines():
            key, v = l.split(':')
            config[key] = v.strip()
        print "read config :", config
        return config


def start_command(cmd, port):
    my_env = os.environ.copy()
    my_env["PORT"] = str(port)
    args = shlex.split('sh -c "{}"'.format(cmd))
    return subprocess.Popen(args, stdout=sys.stdout, stderr=sys.stderr, env=my_env)


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
    print "request url: {}".format(url)
    return getattr(requests, method.lower())(url, headers=headers, data=data)
