# 🛡️ GoAuth - Authentication Service with Docker + PostgreSQL

This project is a secure and scalable authentication microservice built with **Go**, **PostgreSQL**, **JWT**, and **Docker Compose**. It supports login, registration, token refresh, and email-based actions.

---

## 📦 Features

- JWT-based Authentication (Access & Refresh Tokens)
- Secure Email Verification using SMTP
- PostgreSQL Integration
- Containerized with Docker
- Healthcheck & Environment-based configuration
- RESTful API with clean structure

---

## ⚙️ Environment Variables

Create a `.env.production` file in the root of your project with the following contents:

```env
# App
APP_PORT=3000
APP_ENV=production

# Database
DB_HOST=db
DB_PORT=5432
DB_USER=talif
DB_PASSWORD=password
DB_NAME=goauth

# Email SMTP
EMAIL_FROM=pizzaaaa21@gmail.com
EMAIL_PASSWORD=ckozqicaqsqtmgkw
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

# JWT
JWT_SECRET=prod_jwt_secret
JWT_REFRESH_SECRET=prod_jwt_refresh
JWT_EXPIRE_MINUTES=180
JWT_REFRESH_EXPIRE_MINUTES=1440
```
