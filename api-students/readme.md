# API Students - Dokumentasi API

API RESTful untuk manajemen data Mahasiswa yang dibangun menggunakan Go dan framework **Fiber v2**.

## 📌 Base URL

```text
http://localhost:3000/api/v1
```

---

## 🛠️ Persyaratan Header

Untuk method **`POST`**, **`PUT`**, dan **`PATCH`**, request header berikut **wajib** dikirim:

```http
Content-Type: application/json
```

---

## 📊 Format Respons Standar

Semua response dari API mengikuti struktur JSON konsisten:

```json
{
  "success": true,
  "message": "pesan deskriptif",
  "data": null,
  "meta": null,
  "errors": null
}
```

---

## 🚀 Endpoint API

### 1. Health Check

Memeriksa apakah server API berjalan dengan baik.

- **URL:** `/health`
- **Method:** `GET`
- **Response Success (200 OK):**
  ```json
  {
    "success": true,
    "message": "server berjalan",
    "data": {
      "timestamp": "2026-08-26T23:42:00Z"
    }
  }
  ```

---

### 2. Get All Students (Daftar Mahasiswa)

Mengambil daftar mahasiswa dengan dukungan filter, pencarian, pengurutan, dan paginasi.

- **URL:** `/students`
- **Method:** `GET`
- **Query Parameters:**

| Parameter   | Tipe Data | Default | Keterangan |
| :---------- | :-------- | :------ | :--------- |
| `page`      | `int`     | `1`     | Nomor halaman (minimal `1`) |
| `limit`     | `int`     | `10`    | Jumlah data per halaman (minimal `1`, maksimal `100`) |
| `search`    | `string`  | `""`    | Pencarian nama mahasiswa (case-insensitive) |
| `sort`      | `string`  | `"id"`  | Field pengurutan: `id`, `nim`, `name`, `grade` |
| `order`     | `string`  | `"asc"` | Arah pengurutan: `asc` atau `desc` |
| `is_active` | `bool`    | -       | Filter status aktif: `true` atau `false` |

- **Response Success (200 OK):**
  ```json
  {
    "success": true,
    "message": "daftar mahasiswa berhasil diambil",
    "data": [
      {
        "id": 1,
        "nim": "220101001",
        "name": "Budi Santoso",
        "grade": 85.5,
        "is_active": true
      }
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "total_pages": 1
    }
  }
  ```

---

### 3. Get Student By ID (Detail Mahasiswa)

Mengambil detail satu data mahasiswa berdasarkan ID.

- **URL:** `/students/:id`
- **Method:** `GET`
- **URL Parameters:**
  - `id` (integer, wajib, positif)
- **Response Success (200 OK):**
  ```json
  {
    "success": true,
    "message": "mahasiswa ditemukan",
    "data": {
      "id": 1,
      "nim": "220101001",
      "name": "Budi Santoso",
      "grade": 85.5,
      "is_active": true
    }
  }
  ```
- **Response Error (404 Not Found):**
  ```json
  {
    "success": false,
    "message": "mahasiswa tidak ditemukan"
  }
  ```

---

### 4. Create Student (Tambah Mahasiswa)

Menambahkan data mahasiswa baru.

- **URL:** `/students`
- **Method:** `POST`
- **Request Body (JSON):**
  ```json
  {
    "nim": "220101001",
    "name": "Budi Santoso",
    "grade": 85.5,
    "is_active": true
  }
  ```
- **Aturan Validasi:**
  - `nim`: Wajib diisi (string), harus unik.
  - `name`: Wajib diisi (string).
  - `grade`: Rentang angka `0` hingga `100` (`float64`).
  - `is_active`: Tipe data `boolean`.

- **Response Success (201 Created):**
  - **Header:** `Location: /api/v1/students/1`
  - **Body:**
    ```json
    {
      "success": true,
      "message": "mahasiswa berhasil dibuat",
      "data": {
        "id": 1,
        "nim": "220101001",
        "name": "Budi Santoso",
        "grade": 85.5,
        "is_active": true
      }
    }
    ```
- **Response Error Validasi (422 Unprocessable Entity):**
  ```json
  {
    "success": false,
    "message": "validasi gagal",
    "errors": {
      "nim": "wajib diisi",
      "name": "wajib diisi",
      "grade": "nilai harus berada di rentang 0-100"
    }
  }
  ```
- **Response Error Conflict (409 Conflict):**
  ```json
  {
    "success": false,
    "message": "NIM sudah terdaftar"
  }
  ```

---

### 5. Replace Student (Ganti Seluruh Data Mahasiswa - PUT)

Mengubah seluruh data mahasiswa berdasarkan ID.

- **URL:** `/students/:id`
- **Method:** `PUT`
- **Request Body (JSON):**
  ```json
  {
    "nim": "220101001",
    "name": "Budi Santoso Updated",
    "grade": 90.0,
    "is_active": true
  }
  ```
- **Response Success (200 OK):**
  ```json
  {
    "success": true,
    "message": "data mahasiswa berhasil diganti seluruhnya",
    "data": {
      "id": 1,
      "nim": "220101001",
      "name": "Budi Santoso Updated",
      "grade": 90.0,
      "is_active": true
    }
  }
  ```
- **Response Error (404 Not Found / 409 Conflict / 422 Unprocessable Entity)**

---

### 6. Patch Student (Perbarui Sebagian Data Mahasiswa - PATCH)

Memperbarui satu atau beberapa field data mahasiswa berdasarkan ID.

- **URL:** `/students/:id`
- **Method:** `PATCH`
- **Request Body (JSON - opsional per field):**
  ```json
  {
    "grade": 95.0
  }
  ```
- **Response Success (200 OK):**
  ```json
  {
    "success": true,
    "message": "data mahasiswa berhasil diperbarui sebagian",
    "data": {
      "id": 1,
      "nim": "220101001",
      "name": "Budi Santoso Updated",
      "grade": 95.0,
      "is_active": true
    }
  }
  ```
- **Response Error (400 Bad Request - Jika tidak ada field yang dikirim):**
  ```json
  {
    "success": false,
    "message": "tidak ada field yang dikirim untuk diubah"
  }
  ```

---

### 7. Delete Student (Hapus Mahasiswa)

Menghapus data mahasiswa berdasarkan ID.

- **URL:** `/students/:id`
- **Method:** `DELETE`
- **Response Success (204 No Content):** (Tidak mengembalikan body)
- **Response Error (404 Not Found):**
  ```json
  {
    "success": false,
    "message": "mahasiswa tidak ditemukan"
  }
  ```

---

## 🏃 Cara Menjalankan Server

1. Pastikan dependensi Go sudah terpasang:
   ```bash
   go mod tidy
   ```
2. Jalankan aplikasi:
   ```bash
   go run .
   ```
3. Server akan berjalan pada **`http://localhost:3000`**.
