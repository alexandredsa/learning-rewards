## 🛠️ Development Setup – Catalog API

This is a simple GraphQL API built in Go to serve course and category data from a mock catalog.

---

### 📦 Requirements

- Go 1.21+
- Docker & Docker Compose
- `make` (optional, for convenience)
- `ENV=dev` (for development-only DB seeding)
- `ENV=DATABASE_DSN` (default: "DATABASE_DSN=postgres://user:pass@localhost:5432/catalog?sslmode=disable")

---

### 🚀 Running Locally

#### 1. Clone the repository

```bash
git clone https://github.com/<your-org>/learning-rewards.git
cd learning-rewards/catalog-api
```

#### 2. Run Postgres via Docker

```bash
docker compose up -d
```

This spins up a Postgres instance at:

- **Host:** `localhost`
- **Port:** `5432`
- **User:** `user`
- **Password:** `pass`
- **Database:** `catalog`

You can configure this in `main.go` if needed.

#### 3. Set environment for development (for DB seed)

```bash
export ENV=dev
```

You can also create a `.env` file:

```bash
echo "ENV=dev" > .env
```

#### 4. Install gqlgen (once)

```bash
go install github.com/99designs/gqlgen@v0.17.74
```

#### 5. Generate code (if schema changed)

```bash
gqlgen generate
```

#### 6. Run the API

```bash
go run main.go
```

GraphQL Playground will be available at:

```
http://localhost:8080/query
```

---

### 🧪 Example Queries

```graphql
query {
  categories {
    id
    name
  }

  courses(ids: ["<some-uuid>"]) {
    id
    name
    category {
      id
      name
    }
  }
}
```

---

### 🧹 Clean Up

```bash
docker compose down -v
```

This will stop and remove containers + volumes.

