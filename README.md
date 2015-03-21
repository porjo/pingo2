## Pingo2

Pingo2 is a utility for monitoring the availability of hosts and services. It is based on [Pingo](https://github.com/orcheus/pingo), with the following modifications:

- Target address is specified as a URL, with 'http', 'https', 'tcp' and 'ping' as possible schemes
- Email recipient and alert interval can be specified to receive alerts
- Check HTTP virtualhost independently of server hostname. This is useful where GeoDNS might send a request to the closest
  server, rather than the specific server you need to check (works for both http & https).
- Standoff interval: prevents brief flap (host down, then back up) from generating an alert


### Usage

```
./pingo2 -c config.json
```

An example config file is as follows:

```json
{
	"Timeout":10,
	"Standoff":60,
	"SMTP":{
		"Hostname":"localhost",
		"Port":25
	},
	"Alert":{
		"ToEmail":"hostmaster@foobar.org",
		"FromEmail":"noreply@foobar.org",
		"Interval": 900
	},
	"Targets":[
	{
		"Name":"basic HTTP example, default interval 30s",
		"Addr": "http://example.com",
	},
	{
		"Name":"virtualhost example, specific interval, keyword match",
		"Addr": "http://dogbert.example.com",
		"Host": "vhost1.example.com",
		"Interval":20,
		"Keyword":"Look for this phrase"
	},
	{
		"Name":"basic HTTPS example",
		"Addr": "https://secure.example.com",
	},
	{
		"Name":"ping example",
		"Addr": "ping://dogbert.example.com",
	},
	{
		"Name":"tcp example",
		"Addr": "tcp://dogbert.example.com:5432",
	}
	]
}
```
