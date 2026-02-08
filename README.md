# Movie Platform

Project status: functional prototype — backend REST API + simple frontend SPA.
This project is a learning / demo movie platform: it stores movie metadata and reviews, integrates with TMDB for movie metadata and trailers, and provides a small React-free frontend (or plain HTML pages) for browsing, searching, viewing trailers and managing user accounts.

Note: the project does not host full movies — only metadata and trailers (YouTube). Do not attempt to stream copyrighted full-length movies.


*Idea / Overview*

The Movie Platform is a lightweight web application for exploring movies, reading/writing reviews and watching official trailers. It demonstrates a clean Go backend (Gin + layered services + Postgres) integrated with an external provider (TMDB), plus a small client UI that can be used for demos or a course project.


*Goals:*

Simple, well-structured Go backend (handlers → services → repository).

Secure endpoints using JWT authentication for protected actions (create, update, delete, reviews, profile).

Integration with TheMovieDB (TMDB) for additional metadata and trailer URLs.

Minimal frontend that shows movies, search, login/register and a movie detail page with a trailer iframe.


*Functionality*

Public / unauthenticated

List movies (GET /api/movies)

Search movies (GET /api/movies/search?title=...&year=...)

Get movie details (GET /api/movies/:id)

Get movie reviews (GET /api/movies/:id/reviews)

Get TMDB movie metadata + trailer (GET /api/tmdb/movies/:id)

Get popular movies from TMDB and (optionally) import to local DB (GET /api/movies/tmdb/popular)


*Authentication*

Register (POST /api/auth/register)

Login (POST /api/auth/login) → returns JWT token

Authenticated (requires Bearer token)

Create / update / delete movies (POST /api/movies, PUT /api/movies/:id, DELETE /api/movies/:id)

Add / update / delete own reviews (POST /api/movies/:id/reviews, PUT /api/movies/:id/reviews, DELETE /api/movies/:id/reviews)

Admin endpoints to delete any review / user (protected by role)

Profile endpoints (GET /api/me, PUT /api/me, PUT /api/me/password, DELETE /api/me)


*Frontend*

Homepage: grid of movies with search

Movie detail page: metadata, server JSON, TMDB JSON and YouTube trailer iframe

Simple login/register/profile flows (uses localStorage token)


*Architecture*

Project layout (high-level):

```
cmd/server/main.go            # entrypoint, router & route grouping
configs/                      # config loader (configs/config.yaml)
pkg/db                        # postgres connection wrapper
internal/
  ginhandler/                 # HTTP handlers (Gin)
  service/                    # business logic
  postgres/                   # repositories (DB access)
  tmdb/                       # TMDB client
  middleware/                 # JWT auth middleware
model/                        # domain models (Movie, Review, User)
web/                          # static frontend (index.html, movie.html, /static)
```

*Layers:*

Handlers receive HTTP requests, validate and call services.

Services contain business rules and orchestrate calls to repositories or external clients.

Repositories (postgres) talk to Postgres and expose simple DB methods.

External clients (tmdb client) encapsulate calls to external APIs.


*Security:*

JWT for auth, middleware extracts user_id and role.

Passwords must be hashed at registration (server code should use bcrypt — check implementation).

TMDB API requests use a TMDB v4 Bearer Read Access Token (not v3 API key) — see configuration.

Config (example)


*File: configs/config.yaml*

```
database:
  url: "postgresql://user:password@host:5432/dbname"

tmdb:
  # TMDB **Read Access Token (v4)**, looks like: eyJhbGciOiJIUzI1NiJ9...
  api_key: "YOUR_TMDB_V4_READ_ACCESS_TOKEN"

auth:
  jwt_secret: "super-secret-key-123"

```

database.url — Postgres connection string (pgxpool compatible).

tmdb.api_key — must be TMDB v4 Read Access Token (Bearer).

auth.jwt_secret — secret used to sign JWT tokens.


*Database schema (SQL)*

Run these statements in your Postgres to create minimal tables (example):

```
CREATE TABLE movies (
  id SERIAL PRIMARY KEY,
  tmdb_id INT,
  title TEXT NOT NULL,
  year INT,
  description TEXT,
  rating DOUBLE PRECISION DEFAULT 0
);

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT UNIQUE,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  role TEXT DEFAULT 'user',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE reviews (
  id SERIAL PRIMARY KEY,
  movie_id INT REFERENCES movies(id) ON DELETE CASCADE,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  score INT NOT NULL CHECK (score >= 1 AND score <= 5),
  text TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
```

If you use Supabase (as in sample config), make sure the DB user has appropriate privileges.

API Reference (selected endpoints)


*Some example requests and responses:*

GET /api/movies
Response: [{ id, tmdb_id, title, year, description, rating }, ...]

GET /api/movies/search?title=fight&year=1999
Response: filtered list

GET /api/movies/:id
Response: single movie JSON

GET /api/movies/:id/reviews
Response: [{ id, movieId, userId, score, text, createdAt }, ...]

GET /api/tmdb/movies/:tmdb_id
Response: { id, title, description (overview), release_date, trailer_url }

POST /api/auth/register
Body: { "username": "bob", "email": "bob@example.com", "password": "secret" }

POST /api/auth/login
Body: { "email": "bob@example.com", "password": "secret" }
Response: { "token": "<JWT>" }


*Protected endpoints require header: Authorization: Bearer <JWT>*

*Run locally*

*Prerequisites*

*Go 1.20+ (or whichever your project uses)*

*Postgres (local or remote) and configs/config.yaml updated*

*TMDB Read Access Token (v4)*


*Steps:*

1. Update configs/config.yaml with your Postgres URL, TMDB read token and jwt secret.

2. Create database schema (use the SQL above).



*Start the server from project root:*

```
go run ./cmd/server
```
*It will print:*
```
server running on http://localhost:8080
```


*Open the frontend:*

SPA: http://localhost:8080/ (if you use SPA index.html)

Movie page: http://localhost:8080/movie.html?id=<movieID> or ?tmdb=<tmdbID>

Example: Get trailer JSON directly:

```
curl http://localhost:8080/api/tmdb/movies/550
```


*Using Docker Compose (optional)*

*1. You can run Postgres locally with docker-compose. Example docker-compose.yml:*

```
version: "3.8"
services:
  db:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: moviesdb
    ports:
      - "5432:5432"
    volumes:
      - ./pgdata:/var/lib/postgresql/data
```

*2. Then set configs/config.yaml database.url to postgresql://postgres:password@localhost:5432/moviesdb.*


*Frontend (web)*

*Simple static files are stored in web/:*

web/index.html — main SPA (or plain list page)

web/movie.html — movie detail page (uses /api/tmdb/movies/:id to render trailer)

web/static/app.js — SPA JavaScript (if used)

web/static/style.css — CSS


*If you use the SPA, ensure main.go's router serves web/index.html on unknown routes:*

```
r.NoRoute(func(c *gin.Context) {
    c.File("./web/index.html")
})
r.Static("/static", "./web/static")
```


*Notes & Caveats*

1. TMDB token: use TMDB v4 Read Access Token (Bearer). v3 API keys will return 401 when you attempt to use Authorization: Bearer ....

2. Trailers: TMDB provides video metadata; trailers are often YouTube videos. The app extracts the YouTube key and embeds an <iframe> to show the trailer.

3. Full movies: project intentionally does not host or stream full-length copyrighted movies.

4. Auth: JWT secret in configs/config.yaml must be kept secret for production. Passwords should be hashed (bcrypt) — double-check your registration implementation stores hashed passwords, not plain text.

5. Worker: review service runs a background worker to recalculate movie rating after new reviews. Make sure StartRatingWorker() is called in main.


*Testing*

1. Use curl / Postman to test API endpoints.

2. Login, save token, then call protected endpoints with Authorization: Bearer <token>.

3. Check server logs (console) for TMDB response errors (status codes, messages) — helpful when troubleshooting trailer fetching.


*Creators*

Safaryan Artyom
Faizrakhman Alikhan
Ayazbaev Daniyar


*License*

NO


*Possible improvements and extensions:*

Frontend: migrate to a modern SPA toolchain (Vite + React) for maintainability and nicer UX.

Add pagination, filters (rating, genres if added), sorting.

Add admin role UI to manage movies (create/import from TMDB).

Add tests (unit tests for services, integration tests with a test DB).

Add docker-compose for full stack (app + db) and CI pipeline.
