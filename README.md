# Cloudflare DDNS client
This service periodically checks if the IP address of the environment it's currently running in was changed (by calling `https://api.ipify.org`).
If it changed, It will to update the configured Cloudflare CNAME DNS records.

## Configuration
Example config:
```yaml
cfKey: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
cfEmail: your.registered@email.com
zone: zone.net
names:
  - cname.record.name
  - cname2.record.name
```
## Running
Example docker command (config is at `~/ddns-config/config.yaml`):
```sh
docker run -d --restart always \
        -v "~/ddns-config/config.yaml:/otp/cfddns/config.yaml" \
        ghcr.io/balazshorvath/cloudflare-ddns:latest
```
## Kubernetes
`kubernetes/deploy.yaml`: I can't remember if it was working or not, still needs to be tested.
