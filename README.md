# ANX API Client - GO

This is an unofficial pythonic implementation of
[ANX's API](https://dyn.anx.se), described here [API documentation](http://dyn.anx.se/users/apidok.jsf)

Documentation
-------------
This is alpha state software, and I haven't bothered with documentation yet.

Command line client
-------------------
```
Usage: anxdns-go <command>

Flags:
  -h, --help              Show context-sensitive help.
  -k, --apikey            API key used in request header
  -v, --verbose           Verbose
  -b, --baseurl=STRING    Url to API

Commands:
  get [<name>]
    Get Records

  add
    Add Record

  update
    Update Record

  delete
    Delete Record
````
Run "anxdns-go <command> --help" for more information on a command.

````
Usage: anxdns-go get [<name>]

Get Records

Arguments:
  [<name>]    Name of the records to get

Flags:
  -h, --help              Show context-sensitive help.
  -k, --apikey            API key used in request header
  -v, --verbose           Verbose
  -b, --baseurl=STRING    Url to API

  -a, --all               Get all records
  -t, --txt=STRING        Text value of record
````

Client requires two parameters APIKEY and DOMAIN. These can be provided as ENV or arguments in the call.
````
export ANXDNS_APIKEY=keygoeshere
export ANXDNS_DOMAIN=domain.se
````
or
````
./anxdns-go -d domain.se --apikey keygoeshere
````

#### Examples
Get all records 
````
./bin/anxdnsclient -d domain.se --apikey keygoeshere get -a
````
Get records by name
````
./bin/anxdnsapi get -n www.domain.se -d domain.se --apikey keygoeshere <name>
````

Get TXT records by txt
````
./bin/anxdnsapi get -t txt -d domain.se --apikey keygoeshere <name>
````

TODO
-----
* Update of names
* Test cases

Changelog
---------
