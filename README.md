# NERP Wrapper API

NERP Wrapper API เป็น API Gateway สำหรับเชื่อมต่อกับ Odoo ERP โดยใช้ Hexagonal Architecture

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

#### Login

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

## 🔒 Security

- ใช้ HTTPS สำหรับการสื่อสาร
- เก็บ credentials ใน environment variables
- ใช้ admin credentials เฉพาะสำหรับการเชื่อมต่อกับ Odoo API

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- Your Name - Initial work

## 🙏 Acknowledgments

- [Fiber](https://github.com/gofiber/fiber)
- [go-odoo](https://github.com/skilld-labs/go-odoo)
# nerp_warper
