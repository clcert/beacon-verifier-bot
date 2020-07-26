# Beaconbot

Reports if a verifier API is validating correctly a pulse.

# Requirements

Go v1.14 or higher (With GOMODULES)

# Instructions

```bash
go mod tidy
go build
export TG_TOKEN=xxx:xxxx # Telegram bot token. Ask for it with https://t.me/botfather
export TG_GROUP_ID=-1 # Telegram chat ID. It should be negative if it is a group.
export BEACON_VERIFIER_API=http://verifier.random.uchile.cl/chain/1/pulse/latest # Does not exist, is an example.
export IGNORED_SOURCES="ethereum radio" # Separated by spaces. This is useful if a source fails a lot but you already know that. 
export DEBUG=1 # To get a message even when everything is ok.
./beaconbot
```

# Or with `docker-compose`

* Create `.env` file with previously defined constants.
* `docker-compose up -d`