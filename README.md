# Chatty

A modern full-stack chat application built with React TypeScript frontend and Go backend.

## Features

- Real-time messaging interface
- Modern React 18 with TypeScript
- Go backend with Gin framework
- CORS support for cross-origin requests
- Docker containerization
- Responsive design
- Centralized configuration system

## Tech Stack

### Frontend
- React 18
- TypeScript 5
- CSS3 with responsive design
- Vite (fast build tool)

### Backend
- Go 1.21
- Gin web framework
- CORS middleware
- Environment variable support

## Quick Start

### Prerequisites
- Node.js 18+ 
- Go 1.21+
- Docker (optional)

### Development Setup

1. **Clone and setup environment files:**
   ```bash
   git clone <your-repo-url>
   cd chatty
   cp backend/.env.example backend/.env
   cp frontend/.env.example frontend/.env
   ```

2. **Start the backend:**
   ```bash
   cd backend
   go mod tidy
   go run main.go
   ```
   Backend runs on `http://localhost:8080`

3. **Start the frontend:**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```
   Frontend runs on `http://localhost:3000`

### Docker Setup

Run the entire application with Docker Compose:

```bash
docker-compose up --build
```

- Frontend: `http://localhost:3000`
- Backend API: `http://localhost:8080`

## API Endpoints

- `GET /api/health` - Health check
- `GET /api/messages` - Get all messages
- `POST /api/messages` - Send a new message

## Project Structure

```
chatty/
├── backend/
│   ├── main.go          # Go server entry point
│   ├── go.mod           # Go dependencies
│   └── .env.example     # Environment variables template
├── frontend/
│   ├── src/
│   │   ├── App.tsx      # Main React component
│   │   ├── App.css      # Styles
│   │   ├── index.tsx    # React entry point
│   │   └── index.css    # Global styles
│   ├── public/
│   │   └── index.html   # HTML template
│   ├── package.json     # Node dependencies
│   ├── tsconfig.json    # TypeScript configuration
│   ├── vite.config.ts   # Vite configuration
│   └── .env.example     # Frontend environment variables
├── docker-compose.yml   # Multi-container Docker setup
├── Dockerfile.backend   # Backend container
├── Dockerfile.frontend  # Frontend container
├── nginx.conf          # Nginx configuration for frontend
└── README.md           # This file
```

## Development

### Backend Development
```bash
cd backend
go run main.go
```

### Frontend Development
```bash
cd frontend
npm run dev
```

### Type Checking
```bash
cd frontend
npm run type-check
```

### Linting
```bash
cd frontend
npm run lint
```

## Building for Production

### Backend
```bash
cd backend
go build -o chatty-server main.go
```

### Frontend
```bash
cd frontend
npm run build
```

## Environment Variables

### Backend (.env)
```
PORT=8080
GIN_MODE=debug
```

### Frontend (.env)
```
VITE_API_URL=http://localhost:8080
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request