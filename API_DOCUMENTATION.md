# QMS Engine API Documentation

## Overview

The QMS Engine provides a RESTful API for managing projects with full CRUD operations.

## Base URL

```
http://localhost:{PORT}/api/v1
```

## Endpoints

### 1. Create Project

**Endpoint:** `POST /api/v1/projects`

**Request Body:**
```json
{
  "name": "Project Name",
  "description": "Project Description"
}
```

**Response (201 Created):**
```json
{
  "message": "Project created successfully",
  "data": {
    "id": 1,
    "name": "Project Name",
    "description": "Project Description",
    "created_at": "2025-11-07T10:00:00Z",
    "updated_at": "2025-11-07T10:00:00Z",
    "deleted_at": null
  }
}
```

### 2. Get Project by ID

**Endpoint:** `GET /api/v1/projects/:id`

**Response (200 OK):**
```json
{
  "data": {
    "id": 1,
    "name": "Project Name",
    "description": "Project Description",
    "created_at": "2025-11-07T10:00:00Z",
    "updated_at": "2025-11-07T10:00:00Z",
    "deleted_at": null
  }
}
```

### 3. List Projects (with Pagination)

**Endpoint:** `GET /api/v1/projects?limit=20&offset=0`

**Query Parameters:**
- `limit` (optional, default: 20) - Number of projects to return
- `offset` (optional, default: 0) - Number of projects to skip

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Project Name",
      "description": "Project Description",
      "created_at": "2025-11-07T10:00:00Z",
      "updated_at": "2025-11-07T10:00:00Z",
      "deleted_at": null
    }
  ],
  "meta": {
    "limit": 20,
    "offset": 0,
    "count": 1
  }
}
```

### 4. Update Project

**Endpoint:** `PUT /api/v1/projects/:id`

**Request Body:**
```json
{
  "name": "Updated Project Name",
  "description": "Updated Description"
}
```

**Response (200 OK):**
```json
{
  "message": "Project updated successfully",
  "data": {
    "id": 1,
    "name": "Updated Project Name",
    "description": "Updated Description",
    "created_at": "2025-11-07T10:00:00Z",
    "updated_at": "2025-11-07T10:05:00Z",
    "deleted_at": null
  }
}
```

### 5. Delete Project (Soft Delete)

**Endpoint:** `DELETE /api/v1/projects/:id`

**Response (200 OK):**
```json
{
  "message": "Project deleted successfully"
}
```

## Health Check Endpoints

### Health Check
```
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "qms-engine",
  "time": "2025-11-07T10:00:00Z"
}
```

### Ping
```
GET /ping
```

**Response:**
```json
{
  "message": "pong"
}
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid request payload",
  "details": "validation error details"
}
```

### 404 Not Found
```json
{
  "error": "Project not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to create project"
}
```

## Example Usage with cURL

### Create Project
```bash
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Project",
    "description": "This is a test project"
  }'
```

### Get Project
```bash
curl http://localhost:8080/api/v1/projects/1
```

### List Projects
```bash
curl "http://localhost:8080/api/v1/projects?limit=10&offset=0"
```

### Update Project
```bash
curl -X PUT http://localhost:8080/api/v1/projects/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Project",
    "description": "Updated description"
  }'
```

### Delete Project
```bash
curl -X DELETE http://localhost:8080/api/v1/projects/1
```

