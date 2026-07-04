# Infrastructure

This directory contains all infrastructure configuration for the Food Delivery Platform.

## Docker Compose (Local Dev)

Starts all backend services needed for local development:

```bash
# Start everything
docker compose up -d

# Check status
docker compose ps

# View logs
docker compose logs -f

# Stop everything
docker compose down

# Stop and remove volumes (fresh start)
docker compose down -v
```

### Services

| Service | Port | Purpose |
|---------|------|---------|
| PostgreSQL (PostGIS) | 5432 | Primary database (12 logical DBs) |
| Redis | 6379 | Cache + sessions + geospatial |
| Zookeeper | 2181 | Kafka coordination |
| Kafka | 9092 | Event bus |
| Schema Registry | 8081 | Avro schema registry |
| Kafka UI | 9000 | Kafka web UI |
| ElasticSearch | 9200 | Restaurant/menu search |
| ClickHouse | 8123, 9000 | Analytics |
| MinIO | 9001, 9002 | S3-compatible object storage |
| Mailhog | 1025, 8025 | Dev email testing |

### Default Credentials

- **PostgreSQL**: `postgres` / `postgres`
- **Redis**: no password
- **MinIO**: `minio` / `minio123`
- **Kafka**: no auth (dev only)

## Terraform (AWS Production)

TODO: Add Terraform configuration for AWS production environment.

```
terraform/
├── modules/         # Reusable modules (vpc, eks, rds, etc.)
└── environments/    # dev/staging/prod configurations
```

## Kubernetes

TODO: Add K8s manifests and Helm charts.

```
k8s/                 # Raw K8s manifests
helm/                # Helm charts per service
```
