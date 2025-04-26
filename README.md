# Distributed Numbers System

## Overview
This project demonstrates a distributed system using GoLang services deployed in containers. It showcases a simple Publisher-Consumer model, connected through Redis Streams, and backed by PostgreSQL for persistence.

- **Publishers**: 5 independent Go services that publish numbers to a Redis stream.
- **Consumer**: 1 Go service that reads numbers from the Redis stream and persists them to a PostgreSQL database.
- **Infrastructure**: Runs locally using Docker Compose.

---

## Architecture

```
[Publisher Services] ---> [Redis Stream] ---> [Consumer Service] ---> [PostgreSQL DB]
```

- Publishers push numbers into the `numbers-stream` Redis stream.
- Consumer service listens to the stream, consumes new messages, and stores them in the `published_numbers` table.

---

## Prerequisites

- Docker
- Docker Compose
- Go (for local testing scripts)
- Postico (optional, for DB inspection)

---

## Setup Instructions

1. **Clone the Repository**
```bash
git clone https://github.com/R-droid101/distributed-numbers-system
cd distributed-numbers-system
```

2. **Prepare `.env` File**
Create a `.env` file in the root:

```env
DB_USER=user
DB_PASS=password
DB_NAME=numbersdb
DB_HOST=localhost
DB_PORT=5432
AUTH_TOKEN=mysecuretoken
REDIS_ADDR=localhost:6379
```

*Note*: I have setup a static authentication token for the purpose of this project.

3. **Adjust Docker Compose (if needed)**
If you already have Postgres running locally, map Docker's DB to port 5433.

```yaml
ports:
  - "5433:5432"
```

4. **Build and Start Services**
```bash
docker-compose up --build -d
```

5. **Simulate Publisher Activity**
```bash
go test -run simulate_publishers.go
```

6. **Monitor Logs**
```bash
docker-compose logs -f consumer
```

7. **Inspect Database**
Connect using Postico:
- Host: `localhost` (you might need to use the docker IP if you have a local postgreSQL server running)
- Port: `5432` (this could be 5433 based on your configuration)
- User: `user`
- Password: `password`
- Database: `numbersdb`

Execute:
```sql
SELECT * FROM published_numbers;
```

---

## Project Structure

```
/distributed-numbers-system
├── publisher/         # Publisher Go service
├── consumer/          # Consumer Go service
├── docker-compose.yml # Docker Orchestration
├── .env               # Environment Variables
└── README.md          # Project Instructions
```
