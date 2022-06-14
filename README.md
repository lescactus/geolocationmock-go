# geolocationmock-go

This repository contains a simple quick 'n dirty mock server of the [https://ip-api.com/](https://ip-api.com/) and [https://ipbase.com/](https://ipbase.com/) IP geolocation API.

It is mostly used for the development and testing of [geolocation-go](https://github.com/lescactus/geolocation-go).

## Configuration and usage

`geolocationmock-go` is configurable with flags:

```
Usage of ./geolocationmock-go:
  -failure int
    	Failure response rate
  -latency duration
    	Response request latency
  -provider string
    	Provier name. Available: 'ipapi', 'ipbase' (default "ipapi")
```

### Available flags

* `-provider` (default `ipapi`). Start a mock server of the given IP geolocation API. Available values are `ipapi` and `ipbase`.

* `-latency` (default `0s`). Inject given latency into the server response. Examples: `150ms`, `2s`, `1m`, ...

* `-failure` (default `0`). Inject given failure rate into the server response. Examples: `10` for 10%, `50` for 50%.