# Deployment Checklist

## Pre-Deployment

- [ ] All tests passing (`go test ./...`)
- [ ] Update `.env` with production values
- [ ] Change `JWT_SECRET` to secure random string
- [ ] Set `DB_SSLMODE=require` for production
- [ ] Review and set appropriate `PORT`

## Database Setup
```bash
# Create production database
psql -U postgres -h your-db-host
CREATE DATABASE book_api_prod;
\q

# Run migrations
go run cmd/server/main.go
```

## Build for Production
```bash
# Build binary
CGO_ENABLED=0 GOOS=linux go build -o book-api cmd/server/main.go

# Run
./book-api
```

## Environment Variables (Production)

Required:
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `JWT_SECRET` (minimum 32 characters)
- `PORT`

## Health Check
```bash
curl http://your-domain/health
```

Should return: `OK`

## Monitoring

Monitor these metrics:
- Response time
- Error rate
- Database connection pool
- Memory usage
- CPU usage

## Backup

Schedule daily database backups:
```bash
pg_dump -U postgres -d book_api_prod > backup_$(date +%Y%m%d).sql
```