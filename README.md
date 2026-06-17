# divvy

## Setup

### Prerequisites

- [Node 25+](https://nodejs.org/en)
- [Go](https://go.dev/)
- [Docker](https://www.docker.com/)

### .env

Create a .env file in the root of the repository based on `.env.example`

```
POSTGRES_USER=                  # Database username
POSTGRES_PASSWORD=              # Database password
POSTGRES_DB=                    # Database name
POSTGRES_HOST=                  # Database host (optional, defaults to localhost)
POSTGRES_PORT=                  # Database port (optional, defaults to 5432)
```

### Database

1. In the root of the repository, create the docker container for the PostgreSQL database

```bash
docker compose -f docker-compose.db.yaml up -d
```

### Root Scripts

For a streamlined developer experience, root-level scripts are available to initialize and run both the backend and frontend concurrently.

1. **Initialize both projects:**

   ```bash
   npm run init
   ```

   _(This runs `go mod tidy` in `/backend` and `npm install` in `/frontend`)_

2. **Start both development servers:**
   ```bash
   npm run dev
   ```
   _(This starts the backend Go API server and the frontend dev server concurrently)_

---

### Manual Setup & Running

If you prefer to manage the services individually:

#### Backend

1. Ensure the PostgreSQL container is running
2. Move to the backend directory
3. Start the Gin API server
   ```bash
   cd backend
   go run ./cmd
   ```
   > Pass `--skip-migrations` to skip database migrations on startup.

#### Frontend

1. Move to the frontend directory
2. Install dependencies with `npm`
3. Start the development server
   ```bash
   cd frontend
   npm install
   npm run dev
   ```
