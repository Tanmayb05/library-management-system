# Library Management System - Setup Manual

This is a complete library management system built with Go (backend), React (frontend), and PostgreSQL (database). Everything runs in Docker containers for easy setup.

## Prerequisites

Before you start, make sure you have the following installed on your computer:

### 1. Install Docker Desktop

**For Windows:**
1. Go to [Docker Desktop for Windows](https://www.docker.com/products/docker-desktop/)
2. Download and install Docker Desktop
3. Start Docker Desktop after installation
4. Make sure Docker is running (you should see the Docker icon in your system tray)

**For Mac:**
1. Go to [Docker Desktop for Mac](https://www.docker.com/products/docker-desktop/)
2. Download and install Docker Desktop
3. Start Docker Desktop after installation
4. Make sure Docker is running (you should see the Docker icon in your menu bar)

**For Linux (Ubuntu/Debian):**
```bash
# Update package index
sudo apt update

# Install Docker
sudo apt install docker.io

# Start Docker service
sudo systemctl start docker
sudo systemctl enable docker

# Install Docker Compose
sudo apt install docker-compose

# Add your user to docker group (to avoid using sudo)
sudo usermod -aG docker $USER
# Log out and log back in after this command
```

### 2. Install Git (if not already installed)

**Windows:**
- Download from [git-scm.com](https://git-scm.com/download/win)

**Mac:**
```bash
# If you have Homebrew
brew install git

# Or download from git-scm.com
```

**Linux:**
```bash
sudo apt install git
```

## Setup Instructions

### Step 1: Clone the Repository

Open your terminal/command prompt and run:

```bash
# Clone the repository
git clone [YOUR_GITHUB_REPO_URL]

# Navigate to the project directory
cd library-management
```

### Step 2: Verify Project Structure

Make sure your project structure looks like this:

```
library-management/
├── backend/
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── database/
│   │   ├── handlers/
│   │   └── models/
│   ├── Dockerfile
│   ├── docker-compose.yml
│   ├── go.mod
│   └── init.sql
├── frontend/
│   ├── public/
│   │   └── index.html
│   ├── src/
│   │   ├── App.js
│   │   ├── App.css
│   │   └── index.js
│   ├── Dockerfile
│   ├── nginx.conf
│   ├── package.json
│   └── .env
├── docker-compose.yml
└── README.md
```

### Step 3: Build and Run the Application

From the root directory (`library-management/`), run:

```bash
# Build and start all services
docker-compose up --build
```

This command will:
1. Download all required Docker images
2. Build the backend Go application
3. Build the frontend React application
4. Start PostgreSQL database
5. Start all services

**Note:** The first time you run this, it will take several minutes to download and build everything.

### Step 4: Wait for Services to Start

You should see logs from all three services. Wait until you see messages like:

```
library_postgres  | database system is ready to accept connections
library_backend   | Server starting on port 8080...
library_frontend  | /docker-entrypoint.sh: Configuration complete; ready for start up
```

### Step 5: Access the Application

Once all services are running, open your web browser and go to:

- **Frontend (Main App)**: http://localhost:3000
- **Backend API**: http://localhost:8080/health (for testing)

## Using the Application

### Main Features

1. **View Books**: The homepage shows all books in a grid layout
2. **Add New Book**: Click the "Add Book" button to create a new book entry
3. **Edit Book**: Click the edit icon (pencil) on any book card
4. **Delete Book**: Click the trash icon on any book card (with confirmation)
5. **Search**: Use the search box to find books by title, author, or ISBN

### Sample Data

The application comes with 3 sample books pre-loaded:
- The Go Programming Language
- Clean Code
- Design Patterns

## Troubleshooting

### Problem: Port Already in Use

If you get an error saying ports are already in use:

```bash
# Stop any running containers
docker-compose down

# Check what's using the ports
# On Windows/Mac: Check if you have other applications using ports 3000, 8080, or 5432
# On Linux: 
sudo netstat -tulpn | grep :3000
sudo netstat -tulpn | grep :8080
sudo netstat -tulpn | grep :5432

# Try running again
docker-compose up --build
```

### Problem: Docker Not Running

Make sure Docker Desktop is running on your computer. You should see the Docker icon in your system tray (Windows) or menu bar (Mac).

### Problem: Permission Denied (Linux)

If you get permission errors on Linux:

```bash
# Make sure your user is in the docker group
sudo usermod -aG docker $USER
# Then log out and log back in

# Or run with sudo (not recommended for regular use)
sudo docker-compose up --build
```

### Problem: Build Fails

If the build fails:

```bash
# Clean up everything and try again
docker-compose down -v
docker system prune -f
docker-compose up --build
```

### Problem: Application Not Loading

1. Check that all containers are running:
   ```bash
   docker-compose ps
   ```

2. Check the logs for errors:
   ```bash
   docker-compose logs frontend
   docker-compose logs backend
   docker-compose logs postgres
   ```

3. Make sure you're accessing the correct URL: http://localhost:3000

## Stopping the Application

To stop all services:

```bash
# Stop services (keeps data)
docker-compose down

# Stop services and remove all data
docker-compose down -v
```

## Development Mode (Optional)

If you want to make changes to the code and see them immediately:

### Backend Development
```bash
# Start only database
docker-compose up postgres -d

# Navigate to backend directory
cd backend

# Install Go dependencies (requires Go installed)
go mod download

# Run backend locally
go run cmd/main.go
```

### Frontend Development
```bash
# Start backend services
docker-compose up postgres backend -d

# Navigate to frontend directory
cd frontend

# Install Node.js dependencies (requires Node.js installed)
npm install

# Start development server
npm start
```

## API Endpoints (For Testing)

You can test the API directly using curl or a tool like Postman:

```bash
# Health check
curl http://localhost:8080/health

# Get all books
curl http://localhost:8080/api/v1/books

# Create a new book
curl -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "New Book",
    "author": "Author Name",
    "isbn": "1234567890123",
    "publication_year": 2024
  }'
```

## Support

If you encounter any issues:

1. Check the troubleshooting section above
2. Make sure Docker is running and updated
3. Try restarting Docker Desktop
4. Check the GitHub repository for updates

## System Requirements

- **RAM**: At least 4GB recommended
- **Disk Space**: At least 2GB free space
- **OS**: Windows 10/11, macOS 10.14+, or Linux
- **Docker**: Latest version recommended

---

**Note**: The first startup will take longer as Docker needs to download base images. Subsequent startups will be much faster!