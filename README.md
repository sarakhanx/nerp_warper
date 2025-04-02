# NERP Wrapper API

NERP Wrapper is a REST API service that provides an interface to Odoo ERP system, focusing on sales data management.

## Authentication Endpoints

### Login

```bash
POST /auth/login
POST /auth/logout
GET /auth/user-info/:id
```

## Sales Endpoints

### Get All Sale Orders

Retrieves a paginated list of all sale orders.

```bash
GET /sales
```

Query Parameters:

- `page` (optional): Page number (default: 1)
- `page_size` (optional): Number of items per page (default: 500, max: 500)

Example:

```bash
GET /sales?page=1&page_size=100
```

### Get Daily Sales Summary

Retrieves a paginated summary of sales grouped by day.

```bash
GET /sales/daily-summary
```

Query Parameters:

- `page` (optional): Page number (default: 1)
- `page_size` (optional): Number of items per page (default: 500, max: 500)

Example:

```bash
GET /sales/daily-summary?page=1&page_size=30
```

### Get Period Sales Summary

Retrieves sales summary for a specific period or date range.

```bash
GET /sales/period-summary
```

Query Parameters:

- `period_type` (optional): Type of period to summarize. Available options:
  - `1D`: Last 24 hours
  - `7D`: Last 7 days
  - `30D`: Last 30 days (default)
  - `90D`: Last 90 days
  - `MONTHLY`: Current month
  - `YEARLY`: Current year
- `start_date` (optional): Start date for custom range (format: YYYY-MM-DD)
- `end_date` (optional): End date for custom range (format: YYYY-MM-DD)

Examples:

```bash
# Predefined periods
GET /sales/period-summary?period_type=1D
GET /sales/period-summary?period_type=7D
GET /sales/period-summary?period_type=30D
GET /sales/period-summary?period_type=90D
GET /sales/period-summary?period_type=MONTHLY
GET /sales/period-summary?period_type=YEARLY

# Custom date range
GET /sales/period-summary?start_date=2024-01-01&end_date=2024-03-31
```

Notes:

- When using custom date range, both `start_date` and `end_date` must be provided
- For custom date range, `period_type` is optional
- All responses are in JSON format
- All amounts are in the company's default currency
- Only confirmed sales orders (state = 'sale') are included in summaries

## Invoice Endpoints

### Get All Invoices

Retrieves a paginated list of all invoices.

```bash
GET /invoices
```

Query Parameters:

- `page` (optional): Page number (default: 1)
- `page_size` (optional): Number of items per page (default: 500, max: 500)

Example:

```bash
GET /invoices?page=1&page_size=100
```

### Get Daily Invoice Summary

Retrieves a paginated summary of invoices grouped by day.

```bash
GET /invoices/daily-summary
```

Query Parameters:

- `page` (optional): Page number (default: 1)
- `page_size` (optional): Number of items per page (default: 500, max: 500)

Example:

```bash
GET /invoices/daily-summary?page=1&page_size=30
```

### Get Period Invoice Summary

Retrieves invoice summary for a specific period or date range.

```bash
GET /invoices/period-summary
```

Query Parameters:

- `period_type` (optional): Type of period to summarize. Available options:
  - `1D`: Last 24 hours
  - `7D`: Last 7 days
  - `30D`: Last 30 days (default)
  - `90D`: Last 90 days
  - `MONTHLY`: Current month
  - `YEARLY`: Current year
- `start_date` (optional): Start date for custom range (format: YYYY-MM-DD)
- `end_date` (optional): End date for custom range (format: YYYY-MM-DD)

Examples:

```bash
# Predefined periods
GET /invoices/period-summary?period_type=1D
GET /invoices/period-summary?period_type=7D
GET /invoices/period-summary?period_type=30D
GET /invoices/period-summary?period_type=90D
GET /invoices/period-summary?period_type=MONTHLY
GET /invoices/period-summary?period_type=YEARLY

# Custom date range
GET /invoices/period-summary?start_date=2024-01-01&end_date=2024-03-31
```

Notes:

- When using custom date range, both `start_date` and `end_date` must be provided
- For custom date range, `period_type` is optional
- All responses are in JSON format
- All amounts are in the company's default currency
- Only posted invoices (state = 'posted') are included in summaries

## 🚀 Features

- Authentication กับ Odoo
- ระบบ Login/Logout
- ดูข้อมูลผู้ใช้
- ติดตาม Last Login Time

## 🛠️ Tech Stack

- Go 1.21
- Fiber (Web Framework)
- go-odoo (Odoo API Client)
- Hexagonal Architecture

## 📁 Project Structure

```bash
nerp_wrapper/
├── domain/
│   ├── entity/          # Domain entities
│   └── repository/      # Repository interfaces
├── application/
│   ├── dto/            # Data Transfer Objects
│   └── service/        # Business logic
├── infrastructure/
│   └── odoo/           # Odoo implementation
└── interfaces/
    └── http/           # HTTP handlers and routers
```

## 📝 API Endpoints

### Authentication

##### Login

```http
POST /api/auth/login
Content-Type: application/json

{
    "username": "user@example.com",
    "password": "user_password"
}
```

Response:

```json
{
    "success": true,
    "message": "Login successful",
    "user": {
        "id": 1,
        "username": "user@example.com",
        "email": "user@example.com",
        "active": true,
        "last_login": "2024-03-14T15:30:00Z"
    }
}
```

#### Logout

```http
POST /api/auth/logout
```

Response:

```json
{
    "success": true,
    "message": "Logout successful"
}
```

#### Get User Info

```http
GET /api/auth/user/:id
```

Response:

```json
{
    "success": true,
    "message": "User info retrieved successfully",
    "user": {
        "id": 1,
        "username": "user@example.com",
        "email": "user@example.com",
        "active": true,
        "last_login": "2024-03-14T15:30:00Z"
    }
}
```

## 🙏 Acknowledgments

- [Fiber](https://github.com/gofiber/fiber)
- [go-odoo](https://github.com/skilld-labs/go-odoo)

# nerp_warper
