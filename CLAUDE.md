# SummVideo Project Commands and Style Guide

## Build & Run Commands
- Run application: `docker compose up -d`
- Stop application: `docker compose down`
- View logs: `docker compose logs -f`
- Rebuild container: `docker compose build`
- Install dependencies: `pip install -r requirements.txt`

## Code Style Guidelines
- **Python Version**: 3.10
- **Formatting**: PEP 8 compliant
- **Documentation**: Use docstrings for functions and classes
- **Imports**: Group standard library, third-party, and local imports
- **Naming**:
  - Functions: snake_case
  - Variables: snake_case
  - Constants: UPPER_CASE
- **Error Handling**: Use try/except blocks with specific exceptions
- **Language**: Russian comments, English variable names

## FastAPI Specifics
- Use type hints for request and response models
- Organize routes logically with appropriate HTTP methods
- Handle errors with proper status codes and meaningful messages