# TATApps - All-in-One Triastek Management System

TATApps adalah aplikasi web all-in-one untuk mengelola berbagai aspek bisnis termasuk warehouse management, employee management, Support/Ticket, leads/CRM, project management, dan purchase orders dengan dukungan multi-location dan multi-role (Admin, Manager & Employee).

## 🚀 Tech Stack

### Backend
- **Go 1.25.x** - High-performance backend
- **Gin** - Web framework
- **GORM** - ORM for PostgreSQL
- **JWT** - Authentication & Authorization
- **PostgreSQL** - Database
- **robfig/cron** - Scheduler for background jobs (low stock checks)

### Frontend
- **Vue 3** - Progressive JavaScript framework
- **Vite** - Fast build tool
- **Pinia** - State management
- **Vue Router** - Routing
- **TailwindCSS** - Utility-first CSS
- **PrimeVue** - UI component library

### Notifications & Integrations
- **MPWA API** - WhatsApp notifications
- **SMTP** - Transactional email

## ✨ Features

### Access Control & Permissions
- Role-based access (Admin, Manager, Employee) with granular permission keys
- Dynamic menu visibility per role
- User-specific warehouse restrictions for inventory access

### Operational Modules
- **Warehouse Management** - Multi-location warehouse data with assigned managers
- **Inventory Management** - Items, transactions, bulk import/export (CSV & PDF), low stock monitoring
- **Purchase Orders** - Drafting, approval/rejection workflow, supplier & cost breakdown
- **Employee Directory** - Employee data, divisions, positions, and batch operations
- **Category Management** - Shared taxonomy for classifying inventory items

### Settings & Automations
- Company profile and branding assets (logo, favicon)
- Notification preferences for email/WhatsApp and low stock scheduler
- Database backup & restore utilities exposed via API
- WhatsApp and SMTP configuration stored in site settings

## 📋 Prerequisites

- **Go 1.25.3**
- **Node.js 18+** & npm/yarn
- **PostgreSQL 14+**
- **Git**

## 🛠️ Installation & Setup

### 1. Clone Repository

```bash
cd /var/www/html/tatapps
```

### 2. Setup Database

```bash
# Login ke PostgreSQL
sudo -u postgres psql

# Buat database
CREATE DATABASE tatapps;
CREATE USER tatapps_user WITH PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE tatapps TO tatapps_user;
\q
```

### 3. Setup Backend

```bash
cd backend

# Copy environment file
cp .env.example .env

# Edit .env dengan konfigurasi Anda
nano .env

# Install dependencies
go mod download

# Run migrations (akan otomatis saat pertama kali run)
# Jalankan server (migrasi & seeding dev berjalan otomatis)
go run cmd/api/main.go
```

Backend akan berjalan di: `http://localhost:8080`
Saat `APP_ENV=development`, aplikasi akan melakukan:
- Auto-migrate seluruh tabel
- Seeding role, permission, admin user default (`admin@tatapps.com / admin123`), warehouse, dan sample data

### 4. Setup Frontend

```bash
cd ../frontend

# Install dependencies
npm install
# atau
yarn install

# Copy environment file
cp .env.example .env

# Edit .env jika perlu
nano .env

# Run development server
npm run dev
# atau
yarn dev
```

Frontend akan berjalan di: `http://localhost:5173`

## 🔧 Configuration

### Backend (.env)

```env
# Application
APP_NAME=TatApps
APP_ENV=development
APP_PORT=8080
APP_URL=http://localhost:8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=tatapps
DB_USER=tatapps_user
DB_PASSWORD=your_secure_password
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-super-secret-key-change-this
JWT_EXPIRATION=24

# WhatsApp API
WA_API_URL=https://wa.drpnet.my.id/send-message
WA_API_KEY=your_mpwa_api_key
WA_SENDER=your_whatsapp_number

# SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=your-email@gmail.com
SMTP_FROM_NAME=TatApps

# Frontend URL
FRONTEND_URL=http://localhost:5173
```

### Frontend (.env)

```env
VITE_API_URL=http://localhost:8080/api/v1
```

## 📁 Project Structure

```
tatapps/
├── API_DOCUMENTATION.md         # Detail API reference
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go           # API entry point
│   ├── internal/
│   │   ├── config/               # Environment & config loader
│   │   ├── database/             # DB init, migration, seeder
│   │   ├── handlers/             # REST handlers (auth, inventory, settings, etc.)
│   │   ├── middleware/           # JWT, RBAC, and CORS middleware
│   │   ├── models/               # GORM models
│   │   ├── routes/               # Route registration
│   │   └── services/
│   │       └── notification/     # Email & WhatsApp delivery + scheduler
│   ├── go.mod
│   ├── go.sum
│   ├── README.md                 # Backend-specific guide
│   ├── start.sh|stop.sh|restart.sh
│   ├── uploads/                  # Stored site assets (logo, favicon, exports)
│   └── .env
├── database/
│   └── setup.sql                 # Optional helper for manual DB prep
└── frontend/
    ├── src/
    │   ├── api/                  # Axios clients
    │   ├── layouts/              # MainLayout & auth layouts
    │   ├── stores/               # Pinia stores (auth, hr, inventory, etc.)
    │   ├── views/                # Feature pages (inventory, users, settings, warehouse, ...)
    │   ├── App.vue
    │   └── main.js
    ├── package.json
    ├── vite.config.js
    ├── .env
    └── README.md (optional)
```

## 🔐 API Endpoints

Dokumentasi terperinci tersedia di [`API_DOCUMENTATION.md`](API_DOCUMENTATION.md). Ringkasan modul utama:
- **Authentication & Profile** - Login, register, profil, ubah profil, ganti password
- **Warehouses** - CRUD lokasi gudang, guard Admin/Manager
- **Purchase Orders** - Listing, detail, persetujuan/penolakan, update status
- **Inventory** - Item, transaksi, impor/ekspor CSV & PDF, low stock, batch delete
- **Categories** - CRUD kategori inventori
- **Employees** - Data karyawan, divisi, jabatan, batch delete
- **Settings & Notifications** - Site settings (logo, SMTP, WhatsApp), pengaturan notifikasi, backup/restore database, health check
- **User & Role Management** - CRUD user internal, status aktif, menu & permission assignment, daftar roles
- **Notifications** - Trigger test notification dan pengecekan stok rendah terjadwal

## 🎯 Initial Setup & Seeding

### Create Initial Admin User

Setelah backend berjalan, Anda perlu membuat role dan user admin pertama:

```sql
-- Connect to database
psql -U tatapps_user -d tatapps

-- Create roles
INSERT INTO roles (name, description, created_at, updated_at) VALUES
('admin', 'Administrator with full access', NOW(), NOW()),
('manager', 'Manager with limited access', NOW(), NOW()),
('employee', 'Employee with basic access', NOW(), NOW());

-- Create admin user (password: admin123)
-- Hash password using bcrypt dengan cost 10
INSERT INTO users (email, password, full_name, role_id, is_active, created_at, updated_at) VALUES
('admin@tatapps.com', '$2a$10$YourHashedPasswordHere', 'System Admin', 1, true, NOW(), NOW());
```

Atau gunakan API `/api/v1/auth/register` dengan role_id yang sesuai.

## 🚀 Production Deployment

### Backend (Golang)

```bash
# Build binary
cd backend
go build -o tatapps-api cmd/api/main.go

# Run with systemd
sudo nano /etc/systemd/system/tatapps-api.service
```

```ini
[Unit]
Description=TatApps API Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/html/tatapps/backend
ExecStart=/var/www/html/tatapps/backend/tatapps-api
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable tatapps-api
sudo systemctl start tatapps-api
```

### Frontend (Vue 3)

```bash
# Build for production
cd frontend
npm run build

# Serve with Nginx
sudo nano /etc/nginx/sites-available/tatapps
```

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    root /var/www/html/tatapps/frontend/dist;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

```bash
sudo ln -s /etc/nginx/sites-available/tatapps /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 📧 Notification Setup

### Email (SMTP)

Untuk Gmail, aktifkan "App Password":
1. Buka Google Account Settings
2. Security > 2-Step Verification
3. App passwords > Generate
4. Gunakan password tersebut di `.env`

### WhatsApp API

**API Provider:** wa.drpnet.my.id

Setup WhatsApp notifications:
1. Dapatkan API Key dari provider
2. Setup sender number (WhatsApp number)
3. Update di `.env`:
   ```env
   WA_API_URL=https://wa.drpnet.my.id/send-message
   WA_API_KEY=your_api_key
   WA_SENDER=62888xxxx
   ```

📖 **Dokumentasi lengkap:** Lihat [WA_API_CONFIG.md](WA_API_CONFIG.md)

## 🧪 Testing

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests (jika ada)
cd frontend
npm run test
```

## 📝 TODO / Roadmap

- [ ] Implementasi lengkap CRUD untuk semua modul
- [ ] Dashboard analytics & reporting
- [ ] Export to Excel/PDF
- [ ] File upload (avatar, documents)
- [ ] Audit log
- [ ] Email templates yang lebih menarik
- [ ] Notification queue dengan Redis
- [ ] Multi-language support
- [ ] Mobile responsive optimization

## 👥 Default Credentials (Development)

```
Email: admin@tatapps.com
Password: admin123
```

**⚠️ PENTING: Ganti password default saat production!**

## 🆘 Troubleshooting

### Database connection error
- Pastikan PostgreSQL running: `sudo systemctl status postgresql`
- Cek credentials di `.env`
- Cek firewall: `sudo ufw status`

### CORS error
- Pastikan `FRONTEND_URL` di backend `.env` sesuai dengan URL frontend
- Cek middleware CORS di `backend/internal/middleware/cors.go`

### JWT token expired
- Token expired setelah 24 jam (default)
- User harus login ulang
- Atau bisa extend expiration di `.env`

## 📄 License

MIT License - Feel free to use this project for your business needs.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📞 Support

Untuk bantuan lebih lanjut, hubungi tim development atau buat issue di repository.

---

**Built with ❤️ using Golang & Vue 3**
