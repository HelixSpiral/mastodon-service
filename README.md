mastodon-service
---

This is a simple service that reads off an MQTT queue and sends things to Mastodon.
Expected to be run in a docker container.

# Setup instructions
1. checkout this repo to the machine you intend to run this on
2.
```docker build -t mqtt-to-mastodon .```
3. Make a copy of docker-compose.yml and change the "replaceme" values
(**DO NOT ADD YOUR COPY TO GIT**)
