# Chapter 9: Deployment

Deploy your Student Management System to production.

## Overview

Steps to deploy:
1. Build plugins for production
2. Configure production database
3. Set up environment variables
4. Deploy using Docker
5. Configure reverse proxy

## Step 1: Production Build

```bash
# Build server plugin
laatoo plugin build studentmgmt --buildmode prod

# Build UI plugin
laatoo plugin build studentmgmt-ui --buildmode prod --getbuildpackages
```

## Step 2: Database Setup

Create production database:

```sql
CREATE DATABASE school_production;
CREATE USER school_app WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE school_production TO school_app;
```

## Step 3: Environment Configuration

Edit `applications/education/isolations/production/config/isolation.yml`:

```yaml
name: production
description: Production environment

database:
  type: postgres
  host: ${DB_HOST}
  port: ${DB_PORT}
  database: ${DB_NAME}
  user: ${DB_USER}
  password: ${DB_PASSWORD}
  
security:
  publickey: keys/production.pub
  pvtkey: keys/production.pem
```

## Step 4: Docker Deployment

```yaml
# docker-compose.yml
version: '3.8'
services:
  app:
    image: laatoo/server:latest
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
      - DB_NAME=school_production
    volumes:
      - ./applications:/app/applications
      
  db:
    image: postgres:14
    environment:
      - POSTGRES_DB=school_production
    volumes:
      - pgdata:/var/lib/postgresql/data
```

Run:
```bash
docker-compose up -d
```

## Step 5: Nginx Reverse Proxy

```nginx
server {
    listen 80;
    server_name school.example.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
    }
}
```

## Step 6: SSL with Let's Encrypt

```bash
certbot --nginx -d school.example.com
```

## Monitoring

Set up monitoring with:
- Prometheus for metrics
- Grafana for dashboards
- ELK stack for logs

## Summary

You've successfully:
- ✅ Built a complete Student Management System
- ✅ Created entities, services, and workflows
- ✅ Built comprehensive UI with forms and views
- ✅ Deployed to production

Congratulations! You've completed the tutorial.
