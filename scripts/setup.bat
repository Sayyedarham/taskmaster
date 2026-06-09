@echo off
setlocal EnableDelayedExpansion

echo ==========================================
echo   TaskMaster Setup Script (Windows)
echo ==========================================

:: Check prerequisites
echo Checking prerequisites...
docker --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker is required but not installed.
    echo Please install Docker Desktop from https://www.docker.com/products/docker-desktop
    exit /b 1
)

echo [OK] Docker found

:: Determine compose command
docker compose version >nul 2>&1
if not errorlevel 1 (
    set COMPOSE=docker compose
    echo [OK] Docker Compose (plugin) found
) else (
    docker-compose --version >nul 2>&1
    if not errorlevel 1 (
        set COMPOSE=docker-compose
        echo [OK] Docker Compose (standalone) found
    ) else (
        echo [ERROR] Docker Compose not found
        exit /b 1
    )
)

:: Create .env if missing
if not exist .env (
    echo Creating .env from example...
    copy .env.example .env
    echo [OK] .env created
)

:: Stop existing containers
echo Stopping any existing containers...
%COMPOSE% down -v >nul 2>&1

:: Start infrastructure
echo Starting PostgreSQL and Redis...
%COMPOSE% up -d postgres redis

:: Wait for PostgreSQL
echo Waiting for PostgreSQL to be ready...
set RETRIES=30
:wait_loop
%COMPOSE% exec -T postgres pg_isready -U taskmaster >nul 2>&1
if not errorlevel 1 goto pg_ready
set /a RETRIES-=1
if !RETRIES! equ 0 (
    echo [ERROR] PostgreSQL failed to start after 30 attempts
    exit /b 1
)
echo   PostgreSQL not ready yet, retrying... (!RETRIES! left)
timeout /t 2 /nobreak >nul
goto wait_loop
:pg_ready
echo [OK] PostgreSQL is ready

:: Run migrations
echo Running database migrations...
%COMPOSE% up migrate

:: Seed data
echo Seeding demo data...
if exist scripts\seed.bat (
    call scripts\seed.bat
) else (
    echo [WARN] Seed script not found, skipping
)

:: Start API and Frontend
echo Starting API and Frontend...
%COMPOSE% up -d api frontend

echo.
echo ==========================================
echo   Setup Complete!
echo ==========================================
echo.
echo Access points:
echo   API:      http://localhost:8080
echo   Frontend: http://localhost:3000
echo   Health:   http://localhost:8080/health
echo.
echo Demo credentials:
echo   Admin:   admin@taskmaster.com / password123
echo   Manager: manager@taskmaster.com / password123
echo   Member:  member@taskmaster.com / password123
echo.
echo Next steps:
echo   - View logs:     %COMPOSE% logs -f api
echo   - Stop stack:    %COMPOSE% down
echo.

endlocal
