# Coach Calendar

A simple meeting booking application built with Go.

## Project Structure

```
.
├── main.go              # Application entry point and routing
├── database.go          # Database initialization and slot generation
├── handlers/
│   ├── api.go          # API handlers for bookings and admin operations
│   └── pages.go        # Page handlers (home, admin, health)
├── go.mod
└── bookings.db         # SQLite database (generated at runtime)
```

## Features

- **Public Booking Interface** (`/`)
  - Calendar view with Monday-first weeks
  - Available time slot selection
  - Booking form with name and email

- **Admin Panel** (`/admin`)
  - View all slots (available, booked, blocked)
  - Block/unblock time slots
  - View booking details
  - Filter slots by status

## API Endpoints

### Public API
- `GET /api/slots` - Get available time slots
- `POST /api/bookings` - Create a new booking

### Admin API
- `GET /api/admin/slots` - Get all slots with status and booking info
- `POST /api/admin/block` - Block a time slot
- `POST /api/admin/unblock` - Unblock a time slot

## Running the Application

```bash
# Build
go build -o coach-calendar

# Run
./coach-calendar

# Or directly with go run
go run .
```

The server will start on port 8080 by default (configurable via `PORT` environment variable).

## Database

The application uses SQLite with two main tables:
- `bookings` - Stores booking information
- `blocked_slots` - Stores administratively blocked time slots

## Development

The codebase is organized into:
- **main.go** - Minimal entry point with routing
- **database.go** - Database logic and slot generation
- **handlers/** - HTTP request handlers split by concern:
  - API handlers for JSON endpoints
  - Page handlers for HTML pages
