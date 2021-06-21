# DNS Lookup

A simple software that finds IPs addresses that belongs to a given domain using the DNS protocol.

## How it works

Given an input domain, the program will query the google DNS server and it will output all the IPs and subdomains that belongs to the input domain.

It does all this stuff using only the golang standard library and using the DNS protocol specification available in [rfc 1035](https://datatracker.ietf.org/doc/html/rfc1035) and [microsoft documentation](https://docs.microsoft.com/en-us/previous-versions/windows/it-pro/windows-server-2008-R2-and-2008/dd197470(v=ws.10)?redirectedfrom=MSDN#dns-query-message-header)