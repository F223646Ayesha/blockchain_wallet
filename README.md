
# Blockchain Wallet – Full Stack Go + React Application

A secure cryptocurrency wallet system built using:

* **Go + Gin** backend
* **Firestore (Google Firebase)** serverless database
* **React + Tailwind CSS** frontend
* **SHA-256 hashing**, **Proof-of-Work mining**, **UTXO-style logic**
* **Passwordless Email OTP Authentication (SMTP)**
* **Full transaction + mining + analytics + reports + zakat system**

---

## 🔥 Features

### ✔ Secure Authentication

* Passwordless email OTP login using Gmail SMTP
* JWT-based session authentication

### ✔ Wallet System

* Wallet ID generation
* Encrypted PKCS8 private keys
* Public keys stored for signature verification

### ✔ Transactions

* Signed using private key
* Sender → Receiver verification
* Blockchain block linking
* Firestore-backed ledger

### ✔ Blockchain Layer

* SHA-256 block hashing
* Proof-of-Work (difficulty-based nonce mining)
* Merkle Root
* Full-chain validation

### ✔ Beneficiaries

* Add/remove beneficiaries
* Prevents sending to random wallets

### ✔ Zakat System

* Auto-deduction of **2.5% monthly**

### ✔ Reports & Analytics

* Sent/received totals
* Monthly breakdown
* Zakat history
* Mining rewards
* System-wide activity dashboard

---

## 🏗 System Architecture

```
                               ┌─────────────────────────┐
                               │         Frontend        │
                               │   React + Tailwind CSS  │
                               └─────────────┬───────────┘
                                             │ REST API
                                             ▼
                      ┌──────────────────────────────────────────┐
                      │               API Gateway                │
                      │                Go + Gin                  │
                      └───────────────┬──────────────────────────┘
                                      │
         ┌────────────────────────────┼──────────────────────────────┐
         ▼                            ▼                              ▼
 ┌───────────────┐          ┌───────────────────┐          ┌────────────────────┐
 │ Auth Service  │          │ Blockchain Engine │          │ Wallet Service      │
 │ OTP + JWT     │          │ Mining + SHA256   │          │ Balances + Keys     │
 └──────┬────────┘          └──────────┬────────┘          └──────────┬─────────┘
        │                               │                                │
        ▼                               ▼                                ▼
 ┌───────────────┐       ┌────────────────────────┐       ┌────────────────────────┐
 │ Email Service │       │ Transactions Collection │       │ Wallets Collection     │
 │ SMTP/Gmail    │       │ Firestore (serverless)  │       │ Firestore              │
 └───────────────┘       └────────────────────────┘       └────────────────────────┘
```

---

## 🧱 Firestore Database Schema

### `users`

```
{
  name: "Ayena",
  email: "ayena@example.com",
  cnic: "42101-xxxxxxx-x",
  wallet_id: "WALLET123",
  created_at: 1706355220
}
```

### `wallets`

```
{
  wallet_id: "WALLET123",
  public_key: "04ab...",
  private_key_enc: "<AES256 encrypted key>",
  balance: 150,
  beneficiaries: ["WALLET999"],
  created_at: 1706355220
}
```

### `transactions`

```
{
  sender: "WALLET123",
  receiver: "WALLET999",
  amount: 50,
  type: "sent",
  timestamp: 1706355333,
  note: "Payment",
  blockHash: "00000a9f..."
}
```

### `blocks`

```
{
  index: 5,
  previous_hash: "0000x…",
  merkle_root: "2ab1…",
  nonce: 239452,
  hash: "00007c9f…",
  timestamp: 1706355000,
  transactions: [...]
}
```

### `otps`

```
{
  email: "user@example.com",
  otp: "843920",
  createdAt: FirestoreTimestamp
}
```

---

## 📡 API Endpoints

### 🔐 Authentication

| Method | Endpoint             | Description     |
| ------ | -------------------- | --------------- |
| POST   | `/api/auth/send-otp` | Send OTP email  |
| POST   | `/api/login`         | Login using OTP |

---

### 👤 User & Profile

| Method | Endpoint                   | Description       |
| ------ | -------------------------- | ----------------- |
| GET    | `/api/user/profile/:id`    | Get user + wallet |
| POST   | `/api/user/profile/update` | Update profile    |

---

### 💼 Wallet

| Method | Endpoint             | Description   |
| ------ | -------------------- | ------------- |
| POST   | `/api/wallet/create` | Create wallet |
| GET    | `/api/wallet/:id`    | Load wallet   |

#### Beneficiaries

| Endpoint                         | Description        |
| -------------------------------- | ------------------ |
| `/api/wallet/beneficiary/add`    | Add beneficiary    |
| `/api/wallet/beneficiary/remove` | Remove beneficiary |

---

### 💸 Transactions

| Method | Endpoint                       | Description      |
| ------ | ------------------------------ | ---------------- |
| POST   | `/api/transaction/send`        | Send transaction |
| GET    | `/api/transaction/history/:id` | History          |

---

### ⛏ Blockchain

| Method | Endpoint                   | Description    |
| ------ | -------------------------- | -------------- |
| POST   | `/api/mine?miner=W123`     | Mine block     |
| GET    | `/api/blockchain`          | Full chain     |
| GET    | `/api/blockchain/validate` | Validate chain |

---

### 📊 Analytics

| Method | Endpoint                | Description       |
| ------ | ----------------------- | ----------------- |
| GET    | `/api/analytics/system` | System statistics |

---

### 📑 Reports

| Method | Endpoint                  | Description                          |
| ------ | ------------------------- | ------------------------------------ |
| GET    | `/api/reports/wallet/:id` | Totals, zakat, mining, monthly stats |

---

## 🛠 Tech Stack

### Backend

* Go (Golang)
* Gin Framework
* Firestore
* SMTP Email OTP
* JWT Authentication
* Custom Blockchain Engine

### Frontend

* React.js
* TailwindCSS
* Axios
* React Router

---

## ⚙ Setup Instructions

### 🔧 Backend

Install dependencies:

```bash
go mod tidy
```

Create `.env`:

```
FIREBASE_PROJECT_ID=your_project_id
GOOGLE_APPLICATION_CREDENTIALS=backend/serviceAccountKey.json
SMTP_EMAIL=your@gmail.com
SMTP_PASSWORD=your_app_password
JWT_SECRET=supersecret
```

Run backend:

```bash
go run main.go
```

---

### 🎨 Frontend

```bash
cd frontend
npm install
npm run dev
```

---

## 🚀 Deployment (Render)

### Backend

* Build command:

  ```
  go build -o app
  ```
* Start command:

  ```
  ./app
  ```
* Add environment variables in Render dashboard.

### Frontend

```
npm run build
```

Deploy `/dist` folder.

---

## 🛡 Security Notes

* Private keys never leave browser
* AES-256 encrypted private keys
* JWT stored in localStorage
* OTP expires in 5 minutes
* Firestore rules must restrict access

---
