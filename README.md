# Mini EVV Logger

A comprehensive Electronic Visit Logger (EVV) system for caregiver shift tracking with geolocation support.

## Features

- **Schedule Management**: Create, read, update, delete caregiver schedules
- **Visit Tracking**: Track visits with geolocation coordinates and timing
- **Task Management**: Monitor task completion with detailed status tracking
- **Search & Pagination**: Advanced search with pagination support
- **Statistics & Analytics**: Schedule completion metrics and analytics
- **REST API**: Full REST API with proper HTTP methods and validation
- **Docker Support**: Complete Docker and Docker Compose configuration
- **Testing**: Comprehensive unit and integration test coverage

## Architecture

### Backend (Go + Echo Framework)
- **Framework**: Echo web framework
- **Database**: PostgreSQL with proper migrations
- **Architecture**: Clean architecture with separation of concerns
- **Validation**: Comprehensive input validation
- **Documentation**: OpenAPI/Swagger documentation

### Frontend (React + TypeScript)
- **Framework**: React with TypeScript
- **State Management**: React Query for data fetching
- **Routing**: React Router for navigation
- **UI Components**: Modern responsive design
- **Build Tool**: Vite

## Quick Start

### Prerequisites
- Docker and Docker Compose
- PostgreSQL (optional, if running locally)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-boilerplate
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start the application**
   ```bash
   docker-compose up -d
   ```

4. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - API Documentation: http://localhost:8080/swagger/index.html

### Manual Setup

#### Backend
```bash
cd apps/backend

# Install dependencies
go mod download

# Set up database
createdb evv_db
goose postgres "user=evv_user password=evv_password dbname=evv_db sslmode=disable" up

# Run the application
go run main.go
```

#### Frontend
```bash
cd apps/frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

## API Endpoints

### Schedules
- `GET /api/schedules` - Get all schedules with pagination
- `GET /api/schedules/today` - Get today's schedules
- `GET /api/schedules/:id` - Get schedule by ID
- `POST /api/schedules` - Create new schedule
- `PUT /api/schedules/:id` - Update schedule
- `PUT /api/schedules/:id/status` - Update schedule status
- `DELETE /api/schedules/:id` - Delete schedule
- `GET /api/schedules/search?q=query` - Search schedules

### Visits
- `GET /api/schedules/:id/visits` - Get visits for schedule
- `POST /api/schedules/:id/visits/start` - Start visit
- `PUT /api/visits/:id/end` - End visit
- `GET /api/schedules/:id/visit-exists` - Check if visit exists

### Tasks
- `GET /api/schedules/:id/tasks` - Get tasks for schedule
- `POST /api/schedules/:id/tasks` - Create task
- `PUT /api/tasks/:id/status` - Update task status

### Statistics
- `GET /api/schedules/statistics` - Get schedule statistics

## Database Schema

### Schedules Table
- `id` (UUID) - Primary key
- `client_name` (VARCHAR) - Client name
- `shift_time` (VARCHAR) - Shift time (HH:MM-HH:MM)
- `location` (VARCHAR) - Service location
- `status` (VARCHAR) - Schedule status (upcoming, in_progress, completed, missed)
- `visit_id` (UUID) - Reference to visit
- `created_at` (TIMESTAMP) - Creation timestamp
- `updated_at` (TIMESTAMP) - Last update timestamp

### Visits Table
- `id` (UUID) - Primary key
- `schedule_id` (UUID) - Foreign key to schedules
- `start_time` (TIMESTAMP) - Visit start time
- `end_time` (TIMESTAMP) - Visit end time
- `start_latitude` (DECIMAL) - Start latitude
- `start_longitude` (DECIMAL) - Start longitude
- `end_latitude` (DECIMAL) - End latitude
- `end_longitude` (DECIMAL) - End longitude
- `status` (VARCHAR) - Visit status
- `duration_minutes` (INTEGER) - Visit duration
- `created_at` (TIMESTAMP) - Creation timestamp
- `updated_at` (TIMESTAMP) - Last update timestamp

### Tasks Table
- `id` (UUID) - Primary key
- `schedule_id` (UUID) - Foreign key to schedules
- `name` (VARCHAR) - Task name
- `description` (TEXT) - Task description
- `status` (VARCHAR) - Task status (pending, completed, not_completed)
- `reason` (TEXT) - Reason for not completion
- `completed_at` (TIMESTAMP) - Completion timestamp
- `created_at` (TIMESTAMP) - Creation timestamp
- `updated_at` (TIMESTAMP) - Last update timestamp

## Development

### Backend Development
```bash
cd apps/backend

# Run tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...

# Run specific test
go test -v ./internal/service/...

# Run linting
go fmt ./...
go vet ./...

# Build the application
go build -o main .
```

### Frontend Development
```bash
cd apps/frontend

# Run tests
npm test

# Run tests with coverage
npm run test:coverage

# Run linting
npm run lint

# Build the application
npm run build
```

### Code Quality
- Follow Go conventions for backend code
- Use TypeScript for frontend development
- Write comprehensive tests
- Maintain proper documentation
- Follow REST API best practices

## Deployment

### Docker Deployment
```bash
# Build and start all services
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Production Deployment
1. Update environment variables for production
2. Use production database configuration
3. Configure SSL certificates
4. Set up proper logging and monitoring
5. Use orchestration tools like Kubernetes for scaling

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Backend server port | 8080 |
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 5432 |
| `DB_NAME` | Database name | evv_db |
| `DB_USER` | Database username | evv_user |
| `DB_PASSWORD` | Database password | evv_password123 |
| `REDIS_HOST` | Redis host | localhost |
| `REDIS_PORT` | Redis port | 6379 |
| `CORS_ORIGIN` | CORS allowed origin | http://localhost:3000 |
| `ENVIRONMENT` | Environment | development |

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For support and questions, please open an issue in the repository.
