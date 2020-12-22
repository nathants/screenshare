#!/usr/bin/env pypy3
import os
import sys
import argh
import collections
import subprocess
import time
import pool.thread
import tornado.ioloop
import util.misc
import web
import ssl

state = {}

async def handler(request: web.Request) -> web.Response:
    auth = state['auth']
    if auth and [auth] != request['query'].get('auth', []):
        return {'code': 401,
                'body': '401'}
    else:
        with open('index.html') as f:
            return {'code': 200,
                    'body': (f.read()
                             .replace('AUTH', auth or 'AUTH')
                             .replace('MILLIS_PER_FRAME', str(state['millis'])))}

async def js_handler(request: web.Request) -> web.Response:
    auth = state['auth']
    if auth and [auth] != request['query'].get('auth', []):
        return {'code': 401,
                'body': '401'}
    else:
        with open('axios.min.js') as f:
            return {'code': 200,
                    'body': f.read()}

async def img_handler(request: web.Request) -> web.Response:
    auth = state['auth']
    if auth and [auth] != request['query'].get('auth', []):
        return {'code': 401,
                'body': '401'}
    else:
        with open('/tmp/screen.jpg', 'rb') as f:
            return {'code': 200,
                    'body': f.read()}

routes = [
    ('/',             {'get': handler}),
    ('/axios.min.js', {'get': js_handler}),
    ('/img.jpg',      {'get': img_handler}),
]

images = collections.deque()

cmd = """
    path=/tmp/screen.jpg
    while true; do
        maim -g {dimensions} $id -m 7 -f jpg > $path.tmp
        mv -f $path.tmp $path
        echo
    done
"""

@util.misc.exceptions_kill_pid
def screenshotter(dimensions):
    last_print = time.monotonic()
    last = time.monotonic()
    for path in subprocess.Popen(cmd.format(dimensions=dimensions), shell=True, stdout=subprocess.PIPE).stdout:
        now = time.monotonic()
        if now - last_print > 1:
            last_print = now
            print('millis per frame:', int((now - last) * 1000))
        last = now

def main(auth: 'shared secret: http://localhost:8080?auth=AUTH' = None, # type: ignore
         millis: 'millis per frame in browser' = 30, # type: ignore
         dimensions: 'maim dimensions for -g' = '1920x1080', # type: ignore
         port=8080):
    try:
        subprocess.check_output('which maim', shell=True)
    except:
        print('fatal: maim not found: github.com/naelstrof/maim')
        sys.exit(1)
    state['auth'] = auth
    state['millis'] = millis
    web.app(routes).listen(port)
    pool.thread.new(screenshotter, dimensions)
    print('starting screenshare on port:', port)
    tornado.ioloop.IOLoop.current().start()

if __name__ == '__main__':
    os.chdir(os.path.dirname(os.path.abspath(__file__)))
    argh.dispatch_command(main)
