# SummVideo Project Commands and Style Guide

## Build & Run Commands
- Run application: `docker compose up -d`
- Stop application: `docker compose down`
- View logs: `docker compose logs -f`
- Rebuild container: `docker compose build`

## Frontend Development
- Start frontend dev server: `cd frontend && npm run serve`
- Build frontend: `cd frontend && npm run build`
- Test frontend: `cd frontend && npm run test`
- Lint frontend: `cd frontend && npm run lint`

## Go Code Style Guidelines
- Use `gofmt` for code formatting
- Follow standard Go idioms and conventions
- Error handling: always check errors and return them
- Comment exported functions and types
- Use descriptive variable names
- Keep functions small and focused

## Vue.js Guidelines
- Use composition API
- Create reusable components
- Centralize state in Vuex store
- Use TypeScript for type safety

## Project Structure
- Backend: handlers/, services/, models/, repositories/, middleware/
- Frontend: src/components/, src/views/, src/store/

## Common Tasks
- Add a new API endpoint: Create handler in backend-go/handlers/
- Add a new business logic: Create service in backend-go/services/ 
- Update models: Add to backend-go/models/models.go