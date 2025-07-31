# lgtm-lab

lgtm-lab is a small local lab to play with the LGTM stack:
Loki → Grafana → Tempo (i skiped Tempo here) → Mimir/Prometheus
plus MinIO for object storage, and a Go app to generate logs & metrics.

It’s not for production...just a playground to see how LGTM works locally.

---

## What's Inside

| Service        | Purpose                                          |
| -------------- | ------------------------------------------------ |
| **go-app**     | Simple Go API that generates logs and metrics    |
| **promtail**   | Collects container logs and sends them to Loki   |
| **loki**       | Stores logs (backed by MinIO)                    |
| **prometheus** | Stores and queries metrics                       |
| **grafana**    | UI to view logs & metrics                        |
| **minio**      | S3-compatible storage (used by Loki)             |
| **minio-init** | Creates the `loki` bucket in MinIO automatically |


---

## Requirements
    - Docker
    - Docker Compose


---

## Testing 

1. Clone the repo 

```bash
git clone https://github.com/BigBr41n/lgtm-stack-lab.git
cd lgtm-lab
```

2. Start the stack 

```bash
docker-compose up -d
```

3. Check Services

- Go App → http://localhost:8080
- Grafana → http://localhost:3000

    - Default login: admin / admin

- MinIO → http://localhost:9001
    - User: DUMMY_MINIO_ACCESS_KEY
    - Pass: DUMMY_MINIO_SECRET_KEY


---

# Notes

- Logs from go-app → collected by Promtail → stored in Loki (MinIO backend)
- Metrics from go-app → scraped by Prometheus
- View both in Grafana dashboards
- Tempo is not included (this lab focuses on logs + metrics only)
