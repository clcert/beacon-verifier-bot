# Beaconbot

Reports if a verifier API is validating correctly a pulse.

# Requirements

Go v1.14 or higher (With GOMODULES)

# Instructions

```bash
go mod tidy
go build
export TG_TOKEN=xxx:xxxx # Telegram bot token
export TG_GROUP_ID=-1 # Telegram group ID
export BEACON_VERIFIER_API=http://verifier.random.uchile.cl/chain/1/pulse/latest # Does not exist, is an example
./beaconbot
```


