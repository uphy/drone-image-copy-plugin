# Drone Image Copy Plugin

This is a Drone plugin to copy Docker images from a registry to the another one.

drone.yml

```yaml
  copy-images:
    image: uphy/drone-image-copy
    registry: localhost:5000
    images:
      - "hello-world"
      - "bash:4.4"
      - "nginx:latest"
```

