mastodon-service
---

This is a simple service that reads off an MQTT queue and sends things to Mastodon.

# Setup instructions for Kubernetes
1. Modify `mastodon-service.yml` as necessary
2. `kubectl apply -f mastodon-service`, add `-n namespace` if you want it in a specific namespace

# Setup instructions for docker
1. Clone the repo
2. Run `docker build -t mqtt-to-mastodon`
3. Copy `docker-compose.yml.example` to `docker-compose.yml` and edit the values
4. Run `docker compose up -d`
