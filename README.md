# anxdns-go

``
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

Run "anxdns-go <command> --help" for more information on a command.
``

``
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
  ``