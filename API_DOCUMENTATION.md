# API Documentation - TatApps

Base URL: `http://localhost:8080/api/v1`

- All endpoints (kecuali yang disebut **Public**) memerlukan header `Authorization: Bearer <jwt_token>`.
- Response JSON selalu dibungkus dalam bentuk object `{ "data": ..., "message": ... }` kecuali disebutkan lain.
- Timestamp mengikuti format ISO 8601 (`RFC3339`) dan zona waktu default `Asia/Jakarta`.

---

## Authentication

### 1. Login
**POST** `/auth/login`  
Body:
```json
{
  "email": "admin@tatapps.com",
  "password": "admin123"
}
```
Response `200 OK`:
```json
{
  "token": "jwt-token-here",
  "user": {
    "id": 1,
    "email": "admin@tatapps.com",
    "full_name": "System Administrator",
    "role": { "id": 1, "name": "admin", ... },
    "warehouses": [...]
  }
}
```

### 2. Register
**POST** `/auth/register`  
Body:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "full_name": "John Doe",
  "phone": "081234567890",
  "role_id": 3
}
```
Response `201 Created`:
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 2,
    "email": "user@example.com",
    "full_name": "John Doe",
    "role_id": 3,
    "is_active": true
  }
}
```

### 3. Public Site Settings
**GET** `/settings/site` *(Public)*  
Mengambil nama aplikasi, logo, dan favicon untuk halaman login.  
Response `200 OK`:
```json
{
  "data": {
    "app_name": "TatApps",
    "logo": "uploads/site/logo_1700000000000.png",
    "favicon": "uploads/site/favicon_1700000001000.png"
  }
}
```

---

## Authenticated Routes

Tambahkan header berikut pada setiap request:
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

### Profile
- **GET** `/auth/profile` - Detail lengkap user yang sedang login (termasuk role, permission, warehouse yang di-assign).
- **PUT** `/users/profile` - Update profil & upload avatar.
  - Content-Type: `multipart/form-data`
  - Fields: `full_name`, `email`, `phone`, optional file `avatar`.
  - Response `200 OK` menampilkan user terbaru.
- **PUT** `/users/change-password`
  ```json
  {
    "current_password": "oldPass123",
    "new_password": "newPass456"
  }
  ```

---

## Warehouses

| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/warehouses` | List gudang (support query `search`, `is_active`). |
| GET | `/warehouses/:id` | Detail gudang. |
| POST | `/warehouses` | Membuat gudang baru (Admin/Manager). |
| PUT | `/warehouses/:id` | Update gudang (Admin/Manager). |
| DELETE | `/warehouses/:id` | Hapus gudang (Admin). |

Request contoh `POST /warehouses`:
```json
{
  "code": "WH002",
  "name": "Branch Warehouse",
  "address": "Jl. Cabang No. 456",
  "city": "Bandung",
  "province": "Jawa Barat",
  "postal_code": "40123",
  "phone": "022-12345678",
  "email": "branch@tatapps.com",
  "manager_id": 1,
  "color": "#22C55E",
  "is_active": true
}
```

---

## Purchase Orders

| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/purchase-orders` | Query opsional: `status`, `supplier`, `warehouse_id`. |
| GET | `/purchase-orders/:id` | Detail PO + item. |
| POST | `/purchase-orders` | Membuat draft PO. |
| PUT | `/purchase-orders/:id` | Update PO. |
| POST | `/purchase-orders/:id/approve` | Approve PO (Admin/Manager). |
| POST | `/purchase-orders/:id/reject` | Reject PO (Admin/Manager) - body: `{ "reason": "..." }`. |

Request contoh `POST /purchase-orders`:
```json
{
  "po_number": "PO-2024-002",
  "supplier_name": "PT Supplier XYZ",
  "supplier_email": "supplier@xyz.com",
  "supplier_phone": "021-11111111",
  "delivery_date": "2024-01-25T10:00:00Z",
  "warehouse_id": 1,
  "project_id": 1,
  "tax_percent": 11,
  "discount_amount": 0,
  "shipping_cost": 500000,
  "items": [
    {
      "item_name": "Product B",
      "item_code": "PROD-B",
      "description": "High quality product",
      "unit": "pcs",
      "quantity": 50,
      "unit_price": 200000
    }
  ],
  "notes": "Urgent order"
}
```

---

## Inventory

### Items
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/inventory` | Query: `warehouse_id`, `category`, `search`, `low_stock=true`. Employee hanya melihat gudang yang di-assign. |
| GET | `/inventory/items/:id` | Detail item termasuk warehouse. |
| POST | `/inventory/items` | Membuat item baru. |
| PUT | `/inventory/items/:id` | Update sebagian field item. |
| DELETE | `/inventory/items/:id` | Hapus item tunggal. |
| DELETE | `/inventory/items` | Batch delete. Body: `{ "ids": [1,2,3] }`. |

Contoh create item:
```json
{
  "warehouse_id": 1,
  "sn": "SKU-001",
  "name": "Barcode Scanner",
  "description": "Scanner for warehouse operations",
  "category": "Equipment",
  "unit": "pcs",
  "quantity": 10,
  "min_stock": 2,
  "max_stock": 30,
  "unit_price": 1500000,
  "is_active": true
}
```

### Transactions
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/inventory/items/:id/transactions` | Riwayat transaksi per item. |
| POST | `/inventory/items/:id/transactions` | Membuat transaksi `in`, `out`, `transfer`, `adjustment`. Field body mengikuti `InventoryTransaction` (lihat di bawah). |
| GET | `/inventory/transactions` | List transaksi seluruh item. Query: `type`, `warehouse_id`, `start_date`, `end_date`, `search`. |
| DELETE | `/inventory/transactions/:id` | Menghapus transaksi (perlu izin `inventory.delete`). |

Contoh transaksi masuk:
```json
{
  "type": "in",
  "quantity": 5,
  "reference": "PO-2024-002",
  "notes": "Initial stock",
  "to_warehouse_id": 1
}
```

### Import / Export & Monitoring
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/inventory/import/template` | Download template CSV. |
| POST | `/inventory/import/csv` | Import CSV (Content-Type `multipart/form-data`, field `file`). |
| GET | `/inventory/export/csv` | Export inventory ke CSV. |
| GET | `/inventory/export/pdf` | Export inventory ke PDF. |
| GET | `/inventory/transactions/export/csv` | Export transaksi ke CSV. |
| GET | `/inventory/transactions/export/pdf` | Export transaksi ke PDF. |
| GET | `/inventory/low-stock` | Item dengan quantity <= min stock. |

---

## Categories

| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/categories` | Daftar kategori (dengan total item per kategori). |
| GET | `/categories/:id` | Detail kategori + `item_count`. |
| POST | `/categories` | Body: `{ "name": "...", "code": "...", "description": "...", "color": "#2563EB", "is_active": true }`. |
| PUT | `/categories/:id` | Update sebagian field. |
| DELETE | `/categories/:id` | Hanya berhasil bila tidak digunakan oleh item manapun. |

---

## Employees

### Employee Directory
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/employees` | Query: `search`, `division_id`, `position_id`, `status`, `employment_type`. |
| GET | `/employees/:id` | Detail karyawan. |
| POST | `/employees` | Membuat karyawan baru. Menggunakan field camelCase (`fullName`, `divisionId`, dst). |
| PUT | `/employees/:id` | Update data karyawan. |
| DELETE | `/employees/:id` | Hapus karyawan tunggal. |
| DELETE | `/employees` | Batch delete. Body: `{ "ids": [1,2,3] }`. |

Contoh payload create/update:
```json
{
  "employeeCode": "EMP002",
  "nik": "3173000000000002",
  "fullName": "Jane Doe",
  "divisionId": 1,
  "positionId": 3,
  "employmentType": "full-time",
  "status": "active",
  "joinDate": "2024-06-01",
  "phone": "081234567893",
  "email": "jane.doe@tatapps.com"
}
```

### Divisions
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/employees/divisions` | Daftar divisi. |
| POST | `/employees/divisions` | Body: `{ "name": "...", "description": "...", "headEmployeeId": 1 }`. |
| PUT | `/employees/divisions/:id` | Update divisi. |
| DELETE | `/employees/divisions/:id` | Hapus divisi. |

### Positions
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/employees/positions` | Daftar jabatan. |
| POST | `/employees/positions` | Body: `{ "title": "...", "divisionId": 1, "grade": "Senior" }`. |
| PUT | `/employees/positions/:id` | Update jabatan. |
| DELETE | `/employees/positions/:id` | Hapus jabatan. |

---

## User & Role Management (Settings)

Semua endpoint berada di prefix `/settings` dan memerlukan perizinan terkait (`employee.*` atau Admin).

### Users
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/settings/users` | List user internal. Query `role_id` tersedia. |
| GET | `/settings/users/:id` | Detail user + role, permission, warehouse access. |
| POST | `/settings/users` | Body: `{ "full_name": "...", "email": "...", "password": "...", "role_id": 2, "warehouse_ids": [1,2], "send_welcome": true }`. |
| PUT | `/settings/users/:id` | Update data (field sama dengan create). |
| PUT | `/settings/users/:id/status` | Enable/disable user. Body: `{ "is_active": true }`. |
| DELETE | `/settings/users/:id` | Hapus user (tidak bisa menghapus diri sendiri). |

### Roles
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/settings/roles` | Daftar role beserta permission & menu. |
| GET | `/settings/roles/menu-options` | Opsi menu yang tersedia (untuk UI). |
| GET | `/settings/roles/permission-options` | Opsi permission terstandardisasi (sekalian memastikan data di DB). |
| POST | `/settings/roles` | Body: `{ "name": "...", "description": "...", "color": "#2563EB", "menu_keys": ["inventory"], "permission_keys": ["inventory.view"] }`. |
| PUT | `/settings/roles/:id` | Update role + assignment menu/permission. |
| DELETE | `/settings/roles/:id` | Hapus role. |

### Global Roles Endpoint
- **GET** `/roles` - Daftar role ringkas tanpa harus masuk ke `/settings`.

---

## Settings & Notifications

### Site Settings (Admin only)
- **GET** `/settings/site/admin` - Detail lengkap (termasuk konfigurasi WhatsApp & SMTP).
- **PUT** `/settings/site` - Update branding dan credential.
  - Content-Type: `multipart/form-data`
  - Field penting: `app_name`, `whatsapp_api_url`, `whatsapp_api_key`, `whatsapp_sender`, `smtp_host`, `smtp_port`, `smtp_username`, `smtp_password`, `smtp_from_email`, `smtp_from_name`, serta file opsional `logo`, `favicon`.

### Notification Settings (per user)
- **GET** `/settings/notifications` - Preferensi notifikasi low stock.
- **PUT** `/settings/notifications`
  ```json
  {
    "enabled": true,
    "threshold": 10,
    "check_frequency": "daily",
    "schedule_mode": "preset",
    "cron_expression": "0 9 * * *",
    "timezone": "Asia/Jakarta",
    "whatsapp_enabled": true,
    "whatsapp_number": "6281234567890,6289876543210",
    "email_enabled": true,
    "email_address": "ops-team@tatapps.com"
  }
  ```

### Database Maintenance (Admin only)
- **GET** `/settings/database/backup` - Menghasilkan file `*.sql` via `pg_dump`.
- **POST** `/settings/database/restore` - Restore dari file SQL.
  - Content-Type: `multipart/form-data`
  - Field: `backup` (file `.sql`).

### Notification Utilities
- **POST** `/notifications/test`
  ```json
  {
    "whatsapp_enabled": true,
    "whatsapp_number": "6281234567890",
    "email_enabled": true,
    "email_address": "ops-team@tatapps.com"
  }
  ```
- **POST** `/notifications/check-low-stock`
  ```json
  {
    "send_whatsapp": true,
    "send_email": true
  }
  ```
  Mengirimkan ringkasan stok rendah ke kontak user yang login.
- **GET** `/notifications/history?limit=10` - Riwayat notifikasi user.

---

## Health

- **GET** `/health` *(Public)* - Mengembalikan `{ "status": "ok" }`.

---

## Error Responses

| Status | Response |
|--------|----------|
| 400 | `{ "error": "Validation error message" }` |
| 401 | `{ "error": "Invalid or expired token" }` |
| 403 | `{ "error": "Insufficient permissions" }` |
| 404 | `{ "error": "Resource not found" }` |
| 409 | `{ "error": "Email already registered" }` atau konflik lain |
| 422 | `{ "error": "No item IDs provided" }` (contoh) |
| 500 | `{ "error": "Internal server error" }` |
| 206 | `{ "message": "Some notifications failed to send", "errors": [...] }` (partial success) |

---

## Testing with cURL

Login dan simpan token ke variabel:
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@tatapps.com","password":"admin123"}' | jq -r '.token')
```

Ambil profil:
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/auth/profile
```

List inventory (filter low stock):
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/inventory?low_stock=true"
```

Trigger pengecekan low stock manual:
```bash
curl -X POST http://localhost:8080/api/v1/notifications/check-low-stock \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"send_whatsapp": true, "send_email": true}'
```

--- 

Gunakan dokumentasi ini bersama README untuk memahami role/permission serta alur konfigurasi environment.
