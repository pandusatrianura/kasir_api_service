<div align="center">
<img src="./assets/chatBot.png" width="90" alt="Logo" />
<h2> Project Name : Kasir API </h2>

[![Postman](https://img.shields.io/badge/Postman-FF6C37?logo=postman&logoColor=white)](#)
[![Swagger](https://img.shields.io/badge/Swagger-85EA2D?logo=swagger&logoColor=173647)](#)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-6BA539?logo=openapiinitiative&logoColor=white)](#)
[![IntelliJ IDEA](https://img.shields.io/badge/IntelliJIDEA-000000.svg?logo=intellij-idea&logoColor=white)](#)
[![Supabase](https://img.shields.io/badge/Supabase-3FCF8E?logo=supabase&logoColor=fff)](#)
[![Go](https://img.shields.io/badge/Go-00ADD8?logo=Go&logoColor=white&style=for-the-badge)](#)

<img src ="./assets/web.webp" width="80%">

</div>

## 🎯 Material 

- [x] Session 1: https://docs.kodingworks.io/s/17137b9a-ed7a-4950-ba9e-eb11299531c2#h-%F0%9F%8E%AF-tugas
- [x] Session 2: https://docs.kodingworks.io/s/a378a9fe-c0e0-4fa1-a896-43ae347a7b61
- [x] Session 3: https://docs.kodingworks.io/s/820d006c-a994-4487-b993-bc3b4171a35d
- [x] Session 4: https://docs.kodingworks.io/s/69587984-2523-410f-8364-488bec624812

## 💡 Overview

### Task Session 1
1. Implementasikan CRUD Kategori pada project API kalian.

### Task Session 2
1. Pindah categories temen-temen ke layered architecture.
2. **Challange (Optional):** Explore Join, tambah category_id ke table products, setiap product mempunyai kategory, dan ketika Get Detail return category.name dari product.

```json
{
   "code": "1000",
   "message": "Product retrieved successfully",
   "data": {
      "id": 1,
      "name": "Bebelac",
      "price": 10000,
      "stock": 100,
      "category_id": 1,
      "category_name": "Susu",
      "created_at": "2026-01-29T04:27:44.690013+07:00",
      "updated_at": "2026-01-29T04:27:44.690013+07:00"
   }
}
```
### Task Session 3
1. Perbaiki repositories/transaction_repository.go ketika insert transaction details ke db. Perbaiki biar nggak satu-satu insert nya.
   ```bash
	  for i := range details {
		  details[i].TransactionID = transactionID
		  _, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)", transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		  if err != nil {
			   return nil, err
		  }
	  }
   ```
2. Sales Summary Hari Ini
   ```bash
   GET /api/report/hari-ini

   Response:
   {
      "total_revenue": 45000,
      "total_transaksi": 5,
      "produk_terlaris": { 
        "nama": "Indomie Goreng", 
        "qty_terjual": 12 
      }
   }
   
   Optional Challenge
   - Get api/report?start_date=2026-01-01&end_date=2026-02-01
   ```

### Task Session 4
1. To be added.

## ✨ Model

### Category
- **ID**
- **Name**
- **Description**
- **Created At**
- **Updated At**

### Product
- **ID**
- **Name**
- **Price**
- **Stock**
- **Category ID**
- **Created At**
- **Updated At**

### Transaction
- **ID**
- **Total Amount**
- **Created At**
- **Updated At**

### Transaction Detail
- **ID**
- **Transaction ID**
- **Product ID**
- **Quantity**
- **Subtotal**
- **Created At**
- **Updated At**

### Report
- **Total Revenue**
- **Total Transaction**
- **Product with Most Sales**

### Auth Login
- **Username**
- **Password**

## 📖 API Endpoints

The application provides several API endpoints for the functionalities mentioned above. Below are some key endpoints:

### General
- **Health Check API Endpoint**: `GET /api/health/service` 
- **Health Check Database Endpoint**: `GET /api/health/db`
- **Kasir API Docs Endpoint**: `GET /api/docs`

### Category
- **Health Check Category API Endpoint**: `GET /api/categories/health`
- **Ambil semua kategori**: `GET /api/categories`
- **Tambah satu kategori**: `POST /api/categories`
- **Update satu kategori**: `PUT /api/categories/{id}`
- **Ambil detail satu kategori**: `GET /api/categories/{id}`
- **Hapus satu kategori**: `DELETE /api/categories/{id}`

### Product
- **Health Check Product API Endpoint**: `GET /api/products/health`
- **Ambil semua produk**: `GET /api/products`
- **Cari produk berdasarkan nama produk**: `GET /api/products?name=bawang`
- **Tambah satu produk**: `POST /api/products`
- **Update satu produk**: `PUT /api/products/{id}`
- **Ambil detail satu produk**: `GET /api/products/{id}`
- **Hapus satu produk**: `DELETE /api/products/{id}`

### Transaction / Checkout
- **Health Check Transactions/Checkout API Endpoint**: `GET /api/transactions/health`
- **Checkout transaksi**: `POST /api/transactions/checkout`

### Report
- **Health Check Report API Endpoint**: `GET /api/reports/health`
- **Menampilkan laporan penjualan hari ini**: `GET /api/reports/hari-ini`
- **Menampilkan laporan penjualan dengan tanggal tertentu**: `GET /api/reports?start_date=2026-02-04&end_date=2026-02-05` **(Default today if start_date and end_date not provided)**

### Auth
- **Login API Endpoint**: `POST /api/auth/login`
- **Logout API Endpoint**: `POST /api/auth/logout`
- 
## 🛠️ Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/pandusatrianura/kasir_api_service
   ```

2. **Install Dependencies**:
   ```bash
   cd github.com/pandusatrianura/kasir_api_service
   go mod tidy
   ```

3. **Add .env file**
   ```bash
   PORT=8080
   DATABASE_HOST="aws.supabase.com"
   DATABASE_PORT=1234
   DATABASE_USER="postgres"
   DATABASE_PASSWORD="xxxx"
   DATABASE_NAME="kasir-db"
   DATABASE_SSL_MODE="disable"
   DATABASE_MAX_LIFETIME_CONNECTION=30
   DATABASE_MAX_IDLE_CONNECTION= 50
   DATABASE_MAX_OPEN_CONNECTION= 100
   API_KEY= "xxxx"
   JWT_SECRET_KEY= "xxxx"
   JWT_ISSUER= "xxxx"
   JWT_DURATION= "24h"
   ```

4. **Run the Application**:
   ```bash
   go run main.go 
   ```

## 📦 Access the API Documentation
    
1. Use tools like Postman or cURL to interact with the API endpoints.
    
   ```bash
   postman collection: docs/kasir-api.postman_collection.json
   postman environment: docs/kasir-api.postman_environment.json
   ```
2. Use API Documentation to interact with the API endpoints.
    
    ```bash
    swagger: http://{{url}}/api/docs
   ```
   Note: 
   - Replace `{{url}}` with the URL of your deployed API (see 📖 Hosted API).
   - This docs are generated using [swaggo](https://github.com/swaggo/swag).
   - You can also generate the docs locally using `swag init`.

## 📃 List of API Endpoints

### General

1. Kasir API Docs Endpoint:
   ```bash
   curl --location '{{url}}/api/docs'
   ```
   
2. Health Check API Endpoint:
    ```bash
   curl --location '{{url}}/api/health/service'
   ```

3. Health Check Database Endpoint:
    ```bash
   curl --location '{{url}}/api/health/db'
   ``` 

### Category

1. Health Check Endpoint:
   ```bash
   curl --location '{{url}}/api/categories/health'
   ```
2. Display All Categories Endpoint:
   ```bash
   curl --location '{{url}}/api/categories'
   ```
3. Display Category By ID Endpoint:
   ```bash
   curl --location '{{url}}/api/categories/6' \
   --header 'Authorization: Bearer xxx' \
   --header 'X-API-Key: your-secret-api-key-here'
   ```
4. Create New Category Endpoint:
   ```bash
   curl --location '{{url}}/api/categories' \
   --header 'Authorization: Bearer xxx' \
   --header 'X-API-Key: your-secret-api-key-here' \
   --header 'Content-Type: application/json' \
   --data '{
   "name": "Susu",
   "description": "Kategori Susu"
   }'
   ```
5. Update Existing Category Endpoint:
   ```bash
   curl --location --request PUT '{{url}}/api/categories/9' \
   --header 'Authorization: Bearer xxx' \
   --header 'X-API-Key: your-secret-api-key-here' \
   --header 'Content-Type: application/json' \
   --data '{
   "name": "Minuman",
   "description": "Kategori Minuman"
   }'
   ```
6. Delete Existing Category Endpoint:
   ```bash
   curl --location --request DELETE '{{url}}/api/categories/9' \
   --header 'Authorization: Bearer xxx' \
   --header 'X-API-Key: your-secret-api-key-here'
   ```
### Product

1. Health Check Endpoint:
   ```bash
   curl --location '{{url}}/api/products/health'
   ```
2. Display All Products Endpoint:
   ```bash
   curl --location '{{url}}/api/products'
   ```
   
3. Search Products by Product Name Endpoint:
   ```bash
   curl --location '{{url}}/api/products?name=bawang'
   ```
   
4. Display Product By ID Endpoint:
   ```bash
   curl --location '{{url}}/api/products/6' \
   --header 'Authorization: Bearer xxx' \
   --header 'X-API-Key: your-secret-api-key-here'
   ```
5. Create New Product Endpoint:
   ```bash
   curl --location '{{url}}/api/v1/products' \
   --header 'Authorization: Bearer xxx' \
   --header 'X-API-Key: your-secret-api-key-here' \
   --header 'Content-Type: application/json' \
   --data '{
    "name": "Bebelac",
    "price": 10000,
    "stock": 100,
    "category_id": 2
   }'
   ```
6. Update Existing Product Endpoint:
   ```bash
   curl --location --request PUT '{{url}}/api/products/9' \
   --header 'Authorization: Bearer xxx' \
   --header 'X-API-Key: your-secret-api-key-here' \
   --header 'Content-Type: application/json' \
   --data '{
    "name": "Bebelac",
    "price": 10000,
    "stock": 10,
    "category_id": 2
   }'
   ```
7. Delete Existing Product Endpoint:
   ```bash
   curl --location --request DELETE '{{url}}/api/products/9' \
   --header 'Authorization: Bearer xxx' \
   --header 'X-API-Key: your-secret-api-key-here'
   ```

### Transactions

1. Health Check Endpoint:
   ```bash
   curl --location '{{url}}/api/transactions/health'
   ```

2. Checkout Endpoint:
   ```bash
   curl --location '{{url}}/api/transactions/checkout' \
   --header 'Authorization: Bearer xxx' \
   --header 'X-API-Key: your-secret-api-key-here' \
   --header 'Content-Type: application/json' \
   --data '{
    "checkout" : [
        {
            "product_id" : 1,
            "quantity": 5
        },
        {
            "product_id" : 4,
            "quantity": 2
        },
        {
            "product_id" : 5,
            "quantity": 5
        }
    ]
   }'
   ```
   
### Reports

1. Health Check Endpoint:
   ```bash
   curl --location '{{url}}/api/reports/health'
   ```

2. Show Today Sales Report Endpoint:
   ```bash
   curl --location '{{url}}/api/reports/hari-ini' \
   --header 'Authorization: Bearer xxx' \
   --header 'X-API-Key: your-secret-api-key-here'
   ```
3. Show Sales Report With Date Range Endpoint (Default today if start_date and end_date not provided):
   ```bash
   curl --location '{{url}}/api/reports?start_date=2026-02-04&end_date=2026-02-05' \
   --header 'Authorization: Bearer xxx' \
   --header 'X-API-Key: your-secret-api-key-here'
   ```

### Auth
1. Login Endpoint:
   ```bash
   curl --location '{{url}}/api/auth/login' \
   --header 'Content-Type: application/json' \
   --data-raw '{
   "email": "manager@mail.com",
   "password": "manager123"
   }'
   ```
2. logout Endpoint:
   ```bash
   curl --location '{{url}}/api/auth/logout' \
   --header 'Authorization: Bearer xxx' \
   --header 'X-API-Key: your-secret-api-key-here' \
   --data-raw '{}'
   ```

   
**Note:** Replace `{{url}}` with the URL of your deployed API (see 📖 Hosted API).

## 📖 Hosted API

Please change the `{{url}}` to the URL of your deployed API.

   Localhost:
   ```bash
   http://localhost:8080
   ```

   Public URL:
   ```bash
   http://loki-kasir-api-service-a9hfrd-455a6b-203-194-115-248.traefik.me
   ```

   Public Domain:
   ```bash
   https://kasir-api-service.pandusatrianura.cloud
   ```