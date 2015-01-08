## Pingo2
=========

Pingo2 is an open source tool written in Golang that allows you to monitor the avilability of TCP server applications.

It is based on [Pingo](https://github.com/orcheus/pingo), with the following modifications:

- Target address is specified as a URL, with 'http', 'https', and 'ping' as possible schemes
- Email recipient can be specified to receive alerts, together with an alert interval
- Check HTTP virtualhost independently of server hostname. This is useful where GeoDNS might send a request to the closest
  server, rather than the specific server you need to check (works for both http & https).


### Usage

```
./pingo2 -c config.json
```

An example config file is as follows:

```json
{
	"Timeout":10,
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
		"Interval":60,
		"Keyword":"Look for this phrase"
	},
	{
		"Name":"basic HTTPS example",
		"Addr": "https://secure.example.com",
	},
	{
		"Name":"ping example",
		"Addr": "ping://dogbert.example.com",
	}
	]
}
```
