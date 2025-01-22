<h2 align="center">
    ChoccOps/auth
</h2>

## Development

Generate pems

```bash
openssl genrsa -out keys/private.pem 2048
openssl rsa -in keys/private.pem -pubout -out keys/public.pem
```

.env

```bash
VERSION="0.0.0"
POSTGRES_URI="postgresql://auth:serious-password@postgres:5432/auth"
PRIVATE_KEY="./keys/private.pem"
PUBLIC_KEY="./keys/public.pem"

GOOSE_DBSTRING="postgresql://auth:serious-password@postgres:5432/auth"
GOOSE_DRIVER="postgres"
GOOSE_MIGRATION_DIR="./migrations"
```

Start the service

```
docker-compose up
```
