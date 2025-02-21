SMTP Sender
===========

SMTP Sender is a frontend service to send test SMTP commands to a SMTP Proxy, this allows you to send emails from any domain you own without needing to further set these up with a provider.

It is developed in Go, and uses a SQLite database to store state. A docker image is available to build, and will be hosted at a later date.

## Motivation

I make use of a self hosted [anonaddy](https://addy.io/) instance in order to receive non-essential emails. I dont want to manage an SMTP server and the danger that comes with maintaining a good standing email IP, as such I looked at the use of available SMTP Proxies, these are email relays hosted by other companies that you pay in order to use their IP reputation.

The problem with that is that most smtp proxies dont allow forwarding emails as there is no control over what might be sent, if your application is forwarding any email it may forward malicious content and be flagged as spam or worse. As such I have my anonaddy forward all mail to a local test SMTP server, [mailpit](https://mailpit.axllent.org/), but need a way to send emails out. And thus I started working on SMTP Sender.

