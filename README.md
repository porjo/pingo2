## Pingo2
=========

Pingo2 is an open source tool written in Golang that allows you to monitor the avilability of TCP server applications.

It is based on [Pingo](https://github.com/orcheus/pingo), with the following modifications:

- TOML configuration file
- Target address is specified as a URL, and if set to 'http' or 'https' then a HTTP GET is attempted
- Email recipient can be specified to receive alerts, together with an alert interval

