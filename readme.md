# screenshare

## what

screenshare images to a web browser.

## why

ttyd isn't always enough.

## how

```bash
>> pypy3 screenshare.py -h
usage: screenshare [-h] [-c CRT] [-k KEY] [-a AUTH] [-m MILLIS] [-d DIMENSIONS] [-p PORT]

    screenshare by using maim to stream jpgs to a web browser


optional arguments:
  -h, --help            show this help message and exit
  -c CRT, --crt CRT     ssl.crt (default: -)
  -k KEY, --key KEY     ssl.key (default: -)
  -a AUTH, --auth AUTH  shared secret: https://localhost:8080?auth=AUTH (default: -)
  -m MILLIS, --millis MILLIS
                        millis per frame in browser (default: 30)
  -d DIMENSIONS, --dimensions DIMENSIONS
                        '1920x1080'
  -p PORT, --port PORT  8080
```
