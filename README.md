
---

````markdown
# ☕ ANTEIKU — Hot Coffee Management System

**ANTEIKU** is an internal coffee shop management system inspired by a dark, atmospheric aesthetic and built around a robust backend core.  
It combines a RESTful API written in Go with a ready-to-use web interface for managing orders, menu items, inventory, and analytics.

> _“This network exists only for those who are allowed to enter.”_

---

## 🖥️ Overview

The project is designed as a **backend-first system**, where the server encapsulates all critical business logic, while the frontend acts as a visual control panel.

**Key principles:**
- Clear separation of concerns
- Predictable business rules
- Stateless REST API
- Simple and transparent data storage

---

## 🌐 Live Deployment (Temporary Access)

The backend is **temporarily deployed** and publicly available at:

👉 **https://hot-coffee-production.up.railway.app/**

### Notes:
- This deployment is **for demonstration and testing purposes**
- Data persistence depends on the hosting environment
- The server may be restarted or shut down at any time
- Intended for **API testing and frontend integration preview**

You can:
- Connect the web interface directly to this host
- Test all REST endpoints without local setup
- Verify server-side business logic in a live environment

> _Access is granted. For now._

---

## 🚀 Quick Start

```bash
# Build the server
go build -o hot-coffee .

# Run (default port: 8080)
./hot-coffee

# Run with custom settings
./hot-coffee --port 3000 --dir ./data

# Help
./hot-coffee --help
````

After startup:

* **API available at:** `http://localhost:8080`
* **Frontend can be connected immediately**
* **All endpoints are production-ready**

---

## 🎨 Web Interface

The server is fully compatible with the provided web interface.
All endpoints return JSON and are suitable for browser-based interaction.

### Available via UI:

* ✅ Menu management (CRUD)
* ✅ Inventory tracking
* ✅ Order creation and lifecycle control
* ✅ Automatic ingredient deduction
* ✅ Order cancellation with rollback
* ✅ Sales analytics
* ✅ Real-time validation of inventory availability

---

## 🧠 Server-Side Core Logic

The **server is the heart of the system**.
No business rules are duplicated on the frontend.

### Order lifecycle logic:

1. Validate all product IDs against the menu
2. Calculate required ingredients per order
3. Check inventory availability
4. Deduct ingredients atomically
5. Create an order with status `open`
6. Allow only valid state transitions (`open → closed / cancelled`)

### Inventory consistency:

* Ingredients are **deducted on order creation**
* Ingredients are **restored on order cancellation**
* Closed orders are immutable
* Invalid state transitions are rejected

This guarantees **data integrity even without a database**.

---

## 🛠 REST API

### 📋 Orders — `/orders`

| Method   | Endpoint              | Description        |
| -------- | --------------------- | ------------------ |
| `POST`   | `/orders`             | Create a new order |
| `GET`    | `/orders`             | Get all orders     |
| `GET`    | `/orders/{id}`        | Get order details  |
| `PUT`    | `/orders/{id}`        | Update order       |
| `DELETE` | `/orders/{id}`        | Delete order       |
| `POST`   | `/orders/{id}/close`  | Close order        |
| `POST`   | `/orders/{id}/cancel` | Cancel order       |

#### Create order

```json
POST /orders
{
  "customer_name": "Alice",
  "items": [
    { "product_id": "latte", "quantity": 2 },
    { "product_id": "croissant", "quantity": 1 }
  ]
}
```

---

### 🍽 Menu — `/menu`

| Method   | Endpoint     | Description      |
| -------- | ------------ | ---------------- |
| `POST`   | `/menu`      | Create menu item |
| `GET`    | `/menu`      | List menu        |
| `GET`    | `/menu/{id}` | Get menu item    |
| `PUT`    | `/menu/{id}` | Update menu item |
| `DELETE` | `/menu/{id}` | Delete menu item |

---

### 📦 Inventory — `/inventory`

| Method   | Endpoint          | Description       |
| -------- | ----------------- | ----------------- |
| `POST`   | `/inventory`      | Add ingredient    |
| `GET`    | `/inventory`      | List inventory    |
| `GET`    | `/inventory/{id}` | Get ingredient    |
| `PUT`    | `/inventory/{id}` | Update quantity   |
| `DELETE` | `/inventory/{id}` | Remove ingredient |

---

### 📊 Reports — `/reports`

| Method | Endpoint                 | Description       |
| ------ | ------------------------ | ----------------- |
| `GET`  | `/reports/total-sales`   | Total revenue     |
| `GET`  | `/reports/popular-items` | Top-selling items |

> Only **closed orders** are included in reports.

---

## 💾 Initial Data

### Menu (5 items)

* Caffe Latte — $3.50
* Espresso — $2.50
* Cappuccino — $3.80
* Blueberry Muffin — $2.00
* Butter Croissant — $2.50

### Inventory (6 ingredients)

* espresso_shot — 500 units
* milk — 5000 ml
* flour — 10000 g
* blueberries — 2000 g
* sugar — 5000 g
* butter — 3000 g

---

## 📡 HTTP Status Codes

| Code                        | Meaning                   |
| --------------------------- | ------------------------- |
| `200 OK`                    | Successful request        |
| `201 Created`               | Resource created          |
| `204 No Content`            | Resource deleted          |
| `400 Bad Request`           | Validation or logic error |
| `404 Not Found`             | Resource not found        |
| `405 Method Not Allowed`    | Invalid HTTP method       |
| `500 Internal Server Error` | Server failure            |

### Error examples

```json
{"error":"invalid product ID in order items"}
{"error":"insufficient inventory for ingredient 'Milk'"}
{"error":"cannot cancel a closed order"}
```

---

## 🧪 Manual Testing via UI

1. Start the server
2. Open the web interface
3. Connect it to `http://localhost:8080`
4. Test:

   * Order creation
   * Inventory deduction
   * Failed orders due to missing ingredients
   * Order cancellation with rollback
   * Order closing
   * Analytics updates

---

## 🧱 Architecture

```
hot-coffee/
├── main.go            # Entry point (minimal bootstrap)
├── utils/             # JSON helpers & HTTP responses
├── models/            # Domain models
├── internal/
│   ├── config/        # App configuration & DI
│   ├── router/        # HTTP routing
│   ├── dal/           # Data Access Layer
│   ├── service/       # Business logic (core rules)
│   └── handler/       # HTTP handlers
└── data/              # JSON-based storage
```

---

## 🛡 Technical Notes

* **Go 1.22+**
* **Standard library only**
* **No database** — JSON file persistence
* **Strict separation of layers**
* **Business logic isolated in services**
* **Clean, testable architecture**

---

## 👤 Authors

Developed by **mikozhaba&asakhmet**

---

**ANTEIKU is not just a coffee system.**
It is a controlled internal network — quiet, strict, and precise.
