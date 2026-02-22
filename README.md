# Honda Leasing API Backend

Backend API untuk sistem leasing motor berbasis Go, Gin, GORM, dan PostgreSQL.

## Ringkasan
Project ini mengimplementasikan:
- Arsitektur layered: `router -> handler -> service -> repository -> database`
- Domain modular sesuai ERD: `account`, `mst`, `dealer`, `leasing`, `payment`
- Endpoint CRUD generik untuk seluruh resource
- Endpoint workflow leasing end-to-end (`request` sampai `delivery`)

## Tech Stack
- Go `1.25.5` (mengikuti `go.mod`)
- Gin `v1.11.0`
- GORM `v1.31.1`
- PostgreSQL
- Viper (config management)

## Struktur Utama
```text
api/routers/           # registrasi route per domain
internal/handler/      # HTTP handler (CRUD + workflow)
internal/services/     # business logic
internal/repository/   # akses database
internal/domain/models # model GORM
internal/config/       # konfigurasi aplikasi
internal/response/     # format response standar
scripts/tests.sh       # script testing endpoint + coverage
```

## Konfigurasi
Konfigurasi development default ada di:
- `internal/config/configs.development.toml`

Nilai penting default:
- Address: `:8080`
- Base path: `/leasing/api`
- Database host: `localhost:5432`
- Database name: `leasing_db`
- Database user: `postgres`

## Menjalankan Aplikasi
1. Siapkan PostgreSQL dan database sesuai config.
2. Jalankan migration SQL di folder `db/migrations` secara berurutan.
3. Jalankan aplikasi:
```bash
go run .
```

Server akan aktif di:
- `http://localhost:8080/leasing/api`

## Format Response API
Semua endpoint menggunakan envelope response standar.

Success:
```json
{
  "success": true,
  "message": "...",
  "data": {},
  "meta": {}
}
```

Error:
```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "...",
    "details": "..."
  }
}
```

## Query Parameter untuk Endpoint List (CRUD)
Berlaku untuk seluruh endpoint `GET /<resource>`:
- `page` (default: `1`)
- `limit` (default: `20`, max: `100`)
- `sort_by`
- `sort_order` (`ASC|DESC`)
- `search`
- `preload` (pisahkan dengan koma)

## Base URL
Semua endpoint di bawah ini diasumsikan menggunakan prefix:
- `/leasing/api`

Contoh final URL:
- `GET /leasing/api/account/users`

## Daftar Endpoint

### 1) CRUD Endpoint Matrix (Method + Path)
Semua endpoint CRUD berikut aktif.

| Domain | Resource | List | Detail | Create | Update | Delete |
|---|---|---|---|---|---|---|
| account | `/account/oauth_providers` | `GET /account/oauth_providers` | `GET /account/oauth_providers/:id` | `POST /account/oauth_providers` | `PUT /account/oauth_providers/:id` | `DELETE /account/oauth_providers/:id` |
| account | `/account/users` | `GET /account/users` | `GET /account/users/:id` | `POST /account/users` | `PUT /account/users/:id` | `DELETE /account/users/:id` |
| account | `/account/user_oauth_provider` | `GET /account/user_oauth_provider` | `GET /account/user_oauth_provider/:id` | `POST /account/user_oauth_provider` | `PUT /account/user_oauth_provider/:id` | `DELETE /account/user_oauth_provider/:id` |
| account | `/account/roles` | `GET /account/roles` | `GET /account/roles/:id` | `POST /account/roles` | `PUT /account/roles/:id` | `DELETE /account/roles/:id` |
| account | `/account/user_roles` | `GET /account/user_roles` | `GET /account/user_roles/:id` | `POST /account/user_roles` | `PUT /account/user_roles/:id` | `DELETE /account/user_roles/:id` |
| account | `/account/permissions` | `GET /account/permissions` | `GET /account/permissions/:id` | `POST /account/permissions` | `PUT /account/permissions/:id` | `DELETE /account/permissions/:id` |
| account | `/account/role_permission` | `GET /account/role_permission` | `GET /account/role_permission/:id` | `POST /account/role_permission` | `PUT /account/role_permission/:id` | `DELETE /account/role_permission/:id` |
| mst | `/mst/province` | `GET /mst/province` | `GET /mst/province/:id` | `POST /mst/province` | `PUT /mst/province/:id` | `DELETE /mst/province/:id` |
| mst | `/mst/kabupaten` | `GET /mst/kabupaten` | `GET /mst/kabupaten/:id` | `POST /mst/kabupaten` | `PUT /mst/kabupaten/:id` | `DELETE /mst/kabupaten/:id` |
| mst | `/mst/kecamatan` | `GET /mst/kecamatan` | `GET /mst/kecamatan/:id` | `POST /mst/kecamatan` | `PUT /mst/kecamatan/:id` | `DELETE /mst/kecamatan/:id` |
| mst | `/mst/kelurahan` | `GET /mst/kelurahan` | `GET /mst/kelurahan/:id` | `POST /mst/kelurahan` | `PUT /mst/kelurahan/:id` | `DELETE /mst/kelurahan/:id` |
| mst | `/mst/locations` | `GET /mst/locations` | `GET /mst/locations/:id` | `POST /mst/locations` | `PUT /mst/locations/:id` | `DELETE /mst/locations/:id` |
| mst | `/mst/template_tasks` | `GET /mst/template_tasks` | `GET /mst/template_tasks/:id` | `POST /mst/template_tasks` | `PUT /mst/template_tasks/:id` | `DELETE /mst/template_tasks/:id` |
| mst | `/mst/template_task_attributes` | `GET /mst/template_task_attributes` | `GET /mst/template_task_attributes/:id` | `POST /mst/template_task_attributes` | `PUT /mst/template_task_attributes/:id` | `DELETE /mst/template_task_attributes/:id` |
| dealer | `/dealer/motor_types` | `GET /dealer/motor_types` | `GET /dealer/motor_types/:id` | `POST /dealer/motor_types` | `PUT /dealer/motor_types/:id` | `DELETE /dealer/motor_types/:id` |
| dealer | `/dealer/motors` | `GET /dealer/motors` | `GET /dealer/motors/:id` | `POST /dealer/motors` | `PUT /dealer/motors/:id` | `DELETE /dealer/motors/:id` |
| dealer | `/dealer/motor_assets` | `GET /dealer/motor_assets` | `GET /dealer/motor_assets/:id` | `POST /dealer/motor_assets` | `PUT /dealer/motor_assets/:id` | `DELETE /dealer/motor_assets/:id` |
| dealer | `/dealer/customer` | `GET /dealer/customer` | `GET /dealer/customer/:id` | `POST /dealer/customer` | `PUT /dealer/customer/:id` | `DELETE /dealer/customer/:id` |
| leasing | `/leasing/leasing_product` | `GET /leasing/leasing_product` | `GET /leasing/leasing_product/:id` | `POST /leasing/leasing_product` | `PUT /leasing/leasing_product/:id` | `DELETE /leasing/leasing_product/:id` |
| leasing | `/leasing/leasing_contract` | `GET /leasing/leasing_contract` | `GET /leasing/leasing_contract/:id` | `POST /leasing/leasing_contract` | `PUT /leasing/leasing_contract/:id` | `DELETE /leasing/leasing_contract/:id` |
| leasing | `/leasing/leasing_tasks` | `GET /leasing/leasing_tasks` | `GET /leasing/leasing_tasks/:id` | `POST /leasing/leasing_tasks` | `PUT /leasing/leasing_tasks/:id` | `DELETE /leasing/leasing_tasks/:id` |
| leasing | `/leasing/leasing_tasks_attributes` | `GET /leasing/leasing_tasks_attributes` | `GET /leasing/leasing_tasks_attributes/:id` | `POST /leasing/leasing_tasks_attributes` | `PUT /leasing/leasing_tasks_attributes/:id` | `DELETE /leasing/leasing_tasks_attributes/:id` |
| leasing | `/leasing/leasing_contract_documents` | `GET /leasing/leasing_contract_documents` | `GET /leasing/leasing_contract_documents/:id` | `POST /leasing/leasing_contract_documents` | `PUT /leasing/leasing_contract_documents/:id` | `DELETE /leasing/leasing_contract_documents/:id` |
| payment | `/payment/payment_schedule` | `GET /payment/payment_schedule` | `GET /payment/payment_schedule/:id` | `POST /payment/payment_schedule` | `PUT /payment/payment_schedule/:id` | `DELETE /payment/payment_schedule/:id` |
| payment | `/payment/payments` | `GET /payment/payments` | `GET /payment/payments/:id` | `POST /payment/payments` | `PUT /payment/payments/:id` | `DELETE /payment/payments/:id` |

### 2) Leasing Workflow (Custom Endpoint)
| Method | Path | Deskripsi |
|---|---|---|
| `POST` | `/leasing/workflow/submit-application` | Submit pengajuan leasing baru |
| `POST` | `/leasing/workflow/auto-scoring` | Proses auto scoring / manual review |
| `POST` | `/leasing/workflow/survey` | Proses hasil survey |
| `POST` | `/leasing/workflow/final-approval` | Approval final |
| `POST` | `/leasing/workflow/akad` | Eksekusi akad |
| `POST` | `/leasing/workflow/initial-payment` | Catat pembayaran awal |
| `POST` | `/leasing/workflow/dealer-fulfillment` | Proses fulfillment dealer |
| `POST` | `/leasing/workflow/delivery` | Selesaikan delivery |

## Contoh Payload Workflow
Contoh `submit-application`:
```json
{
  "customer_id": 1,
  "motor_id": 1,
  "product_id": 1,
  "dp_dibayar": 7000000,
  "tenor_bulan": 24,
  "request_date": "2026-01-10T00:00:00Z",
  "documents": [
    {
      "file_name": "ktp.jpg",
      "file_size": 10.2,
      "file_type": "jpg",
      "file_url": "https://example.com/ktp.jpg"
    }
  ]
}
```

## Testing Semua Endpoint
Script test end-to-end + coverage endpoint tersedia di:
- `scripts/tests.sh`

Menjalankan test:
```bash
chmod +x scripts/tests.sh
BASE_URL=http://localhost:8080/leasing/api ./scripts/tests.sh
```

Script akan:
- Melakukan create/list/detail/update/delete untuk semua resource CRUD
- Menjalankan seluruh endpoint workflow leasing
- Cleanup data test
- Validasi coverage endpoint

## Catatan Penting
- Saat ini endpoint belum memakai middleware auth/authorization.
- Beberapa resource account memiliki field sensitif (mis. token/secret/password) yang perlu dibatasi di layer response jika dipakai untuk production public API.
