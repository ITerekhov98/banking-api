# 🏦 Banking API (Go)

A RESTful banking service written in Go with PostgreSQL, featuring:

- User registration & JWT-based authentication
- Account and card management
- Funds transfers & balance operations
- Credit system with annuity schedule
- Financial analytics
- External integration with:
  - Central Bank (SOAP key rate)
- Security:
  - PGP encryption (card data)
  - HMAC signatures
  - bcrypt password hashing
- Fully containerized via Docker Compose
- Swagger documentation (enabled via `DEBUG=true`)

---

## 🚀 Quick Start (Dev)

### 1. Clone the project

```bash
git clone https://github.com/ITerekhov98/banking-api.git
cd banking-api
```
### 2. Set up environment
```bash
cp .env.example .env
```
You can customize DB credentials, JWT secret, and debug flags.

### 3. 🔐 Generate PGP key pair (required)
The application uses PGP to encrypt card data.
You must generate and place the following files in:

```
internal/security/keys/
├── pubkey.asc   # Public key (for encryption)
└── privkey.asc  # Private key (for decryption)
```
You can generate a test key pair using GnuPG:

```bash
gpg --quick-gen-key "Test User <test@example.com>" rsa2048 sign,encrypt
gpg --export --armor > internal/security/keys/pubkey.asc
gpg --export-secret-keys --armor > internal/security/keys/privkey.asc
```
Make sure the PGP_PASSPHRASE in .env matches the one used for the key.

### 3. Run with Docker Compose
```bash
docker-compose up --build
```
This will:

Start PostgreSQL

Build and run the Go server

Automatically apply DB migrations

## 📘 API Documentation (Swagger)
Available at: http://localhost:8080/swagger/index.html

Enabled only if DEBUG=true in .env

Add your JWT token via the Authorize button

## 🔐 Authentication
Register: POST /register

Login: POST /login → returns JWT token

Use JWT in Authorization: Bearer <token> for protected endpoints

## 🧾 Database
Uses PostgreSQL 16

Auto-migrations run on app start (in debug mode)

Schema includes:

users, accounts, cards, transactions, credits, payment_schedules

## 📤 External Integrations
CBR SOAP (Key Rate): /keyrate endpoint fetches daily key rate via SOAP

## PGP: 
Card data is encrypted using pubkey.asc/privkey.asc located in internal/security/keys/

