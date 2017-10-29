import json
import logging
import os
import requests
import shlex
import socket
import subprocess
import time
import urlparse
import base64


def handler(event, context):
    logger = logging.getLogger()
    config = read_config()
    port = pick_port()
    p = start_command(config["command"], port)

    event = json.loads(event)

    logger.info("receive event: {} context: {}".format(event, context))

    headers = event.get("headers", {})

    for k, v in headers.items():
        headers[k] = ";".join(v)

    if "Content-Length" in headers:
        del headers["Content-Length"]

    event["headers"] = headers

    method = event["httpMethod"]
    data = event.get("body", None)
    if data and event["isBase64Encoded"]:
        data = base64.b64decode(data)
    path = event.get("path", "/")
    params = event.get("queryParameters", None)

    wait_for_listen(port)
    response = send_request(port=port, headers=headers, method=method, data=data, path=path, params=params)

    data = response.content
    if data:
        data = base64.b64encode(data)

    p.send_signal(subprocess.signal.SIGTERM)
    while p.poll() is None:
        time.sleep(0.1)

    response = {
        "isBase64Encoded": bool(data),
        "statusCode": response.status_code,
        "headers": dict(response.headers),
        "body": data
    }
    return response


def wait_for_listen(port):
    while not is_listen(port):
        time.sleep(0.02)


def read_config():
    with open('qi.json', 'r') as f:
        return json.loads(f.read())


def start_command(cmd, port):
    my_env = os.environ.copy()
    my_env["PORT"] = str(port)
    args = shlex.split(cmd)
    return subprocess.Popen(args, env=my_env, shell=True)


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


def send_request(port, headers, method, data, path, params):
    target = "http://127.0.0.1:{}".format(port)
    url = urlparse.urljoin(target, path)
    print "request url: {}".format(url)
    return getattr(requests, method.lower())(url, headers=headers, data=data, params=params)
