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

install and run screenshare on your laptop:

```bash
>> go build

>> ./screenshare -h
Usage of ./screenshare:
  -a string
        auth: http://localhost:8080?auth=AUTH
  -d int
        display number (default 1)
  -m int
        millis per frame (default 30)
  -p int
        port (default 8080)

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
