# screenshare

## what

screenshare images to a web browser.

## why

ttyd isn't always enough.

## implementations

two identical implementations are maintained:

- [python](https://github.com/nathants/screenshare/tree/python): pypy3 + [maim](https://github.com/naelstrof/maim)
- [go](https://github.com/nathants/screenshare/tree/go): go

## how

```bash
>> pypy3 screenshare.py -h
usage: screenshare [-h] [-a AUTH] [-m MILLIS] [-d DIMENSIONS] [-p PORT]

optional arguments:
  -h, --help            show this help message and exit
  -a AUTH, --auth AUTH  shared secret: http://localhost:8080?auth=AUTH (default: -)
  -m MILLIS, --millis MILLIS
                        millis per frame in browser (default: 30)
  -d DIMENSIONS, --dimensions DIMENSIONS
                        '1920x1080'
  -p PORT, --port PORT  8080
```

on ec2 or any server with a public ip and dns pointing at it, run a reverse proxy to terminate ssl and cache all but a single request to your laptop per millis interval:

```bash
>> cat Caddyfile
{
    order http_cache before reverse_proxy
}

screenshare.example.com  {

    reverse_proxy {
        to localhost:8080
    }

    http_cache {
        cache_type in_memory
        match_path /
        default_max_age 30ms
    }
}

>> go get -u github.com/caddyserver/xcaddy/cmd/xcaddy

>> xcaddy build v2.2.1 --with github.com/sillygod/cdp-cache

>> sudo ./caddy run
```

forward local traffic on 8080 to that server:

```bash
>> ssh -R 8080:0.0.0.0:8080 $user@$server_ip
```

send people the link:

```
https://screenshare.example.com?auth=LOTS_OF_RANDOM_CHARACTERS
```
