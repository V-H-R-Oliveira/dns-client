# DNS Lookup

A simple client software for the DNS protocol.

## Supported Features

- Reverse DNS Queries.
- Standard DNS Queries.
- Support both Ipv4 and Ipv6 protocol.
- Multiple inputs in a domain form or in an ip form.

## How it works

Given an input domain or an input ip, the program will query the google DNS server and it will output all the IPs and domains that belongs to the requested input in JSON. Please, see this [section](#output-example).

It does all this stuff using only the golang standard library and using the DNS protocol specification available in [rfc 1035](https://datatracker.ietf.org/doc/html/rfc1035) and [microsoft documentation](https://docs.microsoft.com/en-us/previous-versions/windows/it-pro/windows-server-2008-R2-and-2008/dd197470(v=ws.10)?redirectedfrom=MSDN#dns-query-message-header)

## Usage example

Basic execution: `go run ./main.go 2001:4860:4860::8888 8.8.4.4 google.com microsoft.com cloudflare.com facebook.com`

You can build a binary using the follow command: `go build -o <some name> main.go` and execute with `./<binary-name> 2001:4860:4860::8888 8.8.4.4 google.com`

## Output Example

Standard DNS Query:

```
{
  "header": {
    "id": 59981,
    "questions": 1,
    "answers": 1,
  },
  "answers": [
    {
      "metadata": {
        "domain": "google.com",
        "type": 1,
        "class": 1,
        "ttl": 20,
        "length": 4
      },
      "resource": "142.250.218.238"
    }
  ],
  "name_servers": null
}
```

Reverse DNS Query:

```
{
  "header": {
    "id": 44802,
    "questions": 1,
    "answers": 1,
  },
  "answers": [
    {
      "metadata": {
        "domain": "4.4.8.8.in-addr.arpa",
        "type": 12,
        "class": 1,
        "ttl": 21189,
        "length": 12
      },
      "resource": "dns.google"
    }
  ],
  "name_servers": null
}
```