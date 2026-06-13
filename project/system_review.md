# 🔍 Honest System Review & Portfolio Rating

Review menyeluruh terhadap **2 proyek** backend yang kamu bangun.

---

## 📊 Scorecard Ringkas

| Aspek | Go Microservice | Express.js API |
|-------|:-:|:-:|
| **Arsitektur** | ⭐⭐⭐⭐ (8/10) | ⭐⭐⭐ (6/10) |
| **Kompleksitas** | ⭐⭐⭐⭐⭐ (9/10) | ⭐⭐⭐ (5/10) |
| **Code Quality** | ⭐⭐⭐⭐ (7/10) | ⭐⭐⭐ (7/10) |
| **Portfolio Value** | ⭐⭐⭐⭐ (8.5/10) | ⭐⭐⭐ (5/10) |
| **Production Readiness** | ⭐⭐⭐ (6/10) | ⭐⭐ (4/10) |

---

## 🏗️ Proyek 1: Go Microservice — Distributed Transaction System

### ✅ Kelebihan (Yang Sudah Bagus)

#### 1. Arsitektur — **8/10**
- **Service Decomposition yang masuk akal**: User, Product, Order, Payment, Broker (API Gateway) — masing-masing punya bounded context yang jelas.
- **API Gateway Pattern**: Broker service sebagai single entry point dengan HTTP→gRPC translation. Ini pattern production-grade yang benar.
- **Event-Driven Architecture**: Kafka untuk `order.created` event yang dikonsumsi Payment service secara async. Ini Saga Pattern yang legitimate.
- **gRPC inter-service communication**: Pilihan yang tepat untuk internal service communication (lebih cepat dari REST, strongly typed via protobuf).
- **Clean layer separation**: `cmd/ → transport/ → service/ → repository/` di setiap service. Consistent across all services.

#### 2. Kompleksitas — **9/10**
Ini yang bikin proyek ini **sangat kuat untuk portfolio**:

| Teknologi | Implementasi | Level |
|-----------|-------------|-------|
| gRPC + Protobuf | 4 service definitions, bidirectional calls | 🟢 Advanced |
| Kafka | Producer (Order) + Consumer (Payment) | 🟢 Advanced |
| Redis Caching | Cache-aside pattern, TTL differentiation | 🟢 Advanced |
| Token Bucket Rate Limiter | Custom Lua script di Redis (atomic!) | 🟢 Advanced |
| Payment Gateway | Midtrans Snap API + Webhook SHA512 verification | 🟢 Advanced |
| OAuth 2.0 | Google, Facebook, GitHub dengan PKCE | 🟢 Advanced |
| Docker + K8s | Multi-stage Dockerfile, K8s deployments, secrets, health checks | 🟢 Advanced |
| CI/CD | Jenkins pipeline (test, build, push, deploy) | 🟡 Intermediate |
| Database Migrations | golang-migrate with up/down | 🟢 Good |
| Frontend | Next.js + React Query + Zustand + TypeScript | 🟢 Good |

> [!TIP]
> Dari segi **technology breadth**, ini sudah sangat impressive. Recruiter/interviewer pasti akan tertarik karena kamu menyentuh hampir semua buzzword yang dicari: microservices, event-driven, gRPC, Kafka, K8s, payment gateway.

#### 3. Code Quality — **7/10**
- Lua script untuk rate limiter rapi dan atomic — ini menunjukkan pemahaman mendalam tentang Redis.
- Transaction management di `CreateOrder` benar: `Begin → operations → Commit`, dengan proper rollback defer.
- Cache invalidation strategy yang thoughtful: invalidate list + single order cache saat status update.
- Idempotency di payment initiation (reuse existing gateway token jika masih pending).

---

### ❌ Kelemahan (Yang Perlu Diperbaiki)

#### 1. **Tidak ada unit test sama sekali** ⚠️
```
broker/  → 0 test files
user/    → 0 test files  
order/   → 0 test files
payment/ → 0 test files
product/ → 0 test files
```

> [!CAUTION]
> **Ini kelemahan terbesar.** Proyek sekompleks ini TANPA unit test akan sangat menurunkan nilai di mata senior engineer. Bahkan Jenkinsfile-nya mengarah ke `./internal/service` (path yang bahkan tidak ada di project ini — tampaknya copy-paste dari proyek lain).

#### 2. **Tidak ada Interface/Contract untuk service layer**
```go
// Sekarang: concrete dependency
type OrderService struct {
    orderRepo repository.OrderRepository  // ← Ini sudah interface? Perlu dicek
}
```
Kalau `OrderRepository` sudah interface, bagus. Tapi service layer sendiri tidak expose interface, jadi sulit di-mock untuk testing.

#### 3. **Error Handling terlalu generik**
```go
return nil, errors.New("stock is not enough")  // ← Tidak ada error code/type
return nil, err  // ← Raw error propagation tanpa wrapping
```
Tidak ada custom error types, error codes, atau structured error responses yang konsisten.

#### 4. **Hardcoded values dan environment coupling**
```go
func (u *OrderService) CreateOrder(...) {
    addr := []string{os.Getenv("KAFKA_BROKER_URL")}  // ← Di dalam business logic!
    topic := os.Getenv("KAFKA_ORDER_TOPIC")           // ← Harusnya di-inject
}
```
Environment variables dibaca langsung di service layer — melanggar dependency injection principle.

#### 5. **Tidak ada observability stack**
- Tidak ada distributed tracing (OpenTelemetry/Jaeger)
- Tidak ada metrics (Prometheus)
- Logging hanya pakai logrus tanpa structured fields yang konsisten
- Tidak ada correlation ID antar service

#### 6. **Database schema terlalu simpel**
- Semua pakai `SERIAL` (int auto-increment) — UUID lebih aman untuk distributed system
- Tidak ada soft delete
- `DOUBLE PRECISION` untuk harga — harusnya `DECIMAL`/`NUMERIC` untuk financial precision
- Tidak ada indexing strategy yang terlihat

#### 7. **Security gaps**
- Tidak ada input validation di gRPC level (hanya di HTTP handler broker)
- Webhook endpoint tidak rate-limited
- Password policies tidak di-enforce
- Refresh token rotation belum ada (token lama masih bisa dipakai setelah refresh)

#### 8. **Jenkinsfile untuk proyek lain**
```groovy
sh "go test ./internal/service -v"  // ← path ini tidak ada di project ini
dockerImage = docker.build("numpyh/currency-exchange:${env.BUILD_TAG}")  // ← Image name salah
```
Ini clearly copy-paste dari proyek `currency-exchange` dan belum disesuaikan.

---

## 🏗️ Proyek 2: Express.js API

### ✅ Kelebihan
- **Separation of concerns yang benar**: Controller → Service → Prisma (ORM). Ini idiomatic Express.
- **Dual-token JWT**: Access + Refresh token dengan rotation — lebih aman dari single token.
- **Zod validation** di route level — modern dan type-safe.
- **Custom AppError class** dengan operational vs programming error distinction.
- **Token Bucket Rate Limiter** pakai Redis Lua script — menunjukkan pemahaman yang konsisten dengan Go project.
- **Ownership-based access control**: User B tidak bisa edit/delete post milik User A.
- **Comprehensive integration test** (`test_api.js`) yang cover negative testing, auth flow, ownership security.

### ❌ Kelemahan
- **Ini cuma CRUD biasa** — User + Post, relasi one-to-many. Tidak ada yang "wow" dari sisi domain complexity.
- **Tidak ada caching layer** — Padahal Redis sudah ada untuk rate limiter.
- **SQLite (dev.db)** — Meskipun Prisma schema bilang PostgreSQL, ada file `dev.db` yang menunjukkan SQLite di dev.
- **Tidak ada Dockerfile atau deployment config** — Monolith tanpa containerization.
- **Tidak ada pagination yang proper** — Post list ada, tapi tidak ada cursor-based pagination.
- **Tidak ada logging framework** — Hanya `console.error`.
- **`catchAsync` didefinisikan 2x** — Di `authMiddleware.js` dan `postController.js`, copy-paste.

> [!WARNING]
> **Untuk portfolio, proyek ini KURANG.** CRUD User + Post sudah terlalu common dan tidak membedakan kamu dari kandidat lain. Proyek ini lebih cocok sebagai *learning exercise* daripada portfolio showcase.

---

## 🎯 Rating Portofolio

### Go Microservice: **8.5/10** — 🟢 SANGAT LAYAK untuk Portfolio

**Mengapa tinggi:**
- Menyentuh hampir semua teknologi yang dicari industri
- Menunjukkan kemampuan system design, bukan hanya coding
- Real payment gateway integration (bukan mock)
- Fullstack: backend + frontend + infra

**Mengapa bukan 10:**
- Zero test coverage (fatal flaw)
- Jenkinsfile copy-paste dari proyek lain
- Beberapa code smells (env reads di service layer, no interfaces)
- Tidak ada observability

### Express.js API: **5/10** — 🟡 KURANG untuk Portfolio

**Mengapa rendah:**
- Terlalu simpel (CRUD user + post)
- Tidak ada domain complexity
- Tidak ada deployment/infra
- Sudah terlalu banyak proyek serupa di GitHub

---

## 🚀 Rekomendasi: Apa yang Perlu Ditambahkan

### Go Microservice — Priority Tinggi

#### 1. **Unit Tests + Integration Tests** (WAJIB)
```
Effort: 2-3 hari
Impact: Portfolio value +2 poin
```
- Unit test service layer dengan mock repository
- Integration test gRPC endpoints
- Test coverage badge di README
- Fix Jenkinsfile agar mengarah ke path yang benar

#### 2. **Distributed Tracing (OpenTelemetry)**
```
Effort: 1-2 hari
Impact: Portfolio value +1 poin
```
- Tambahkan trace propagation antar service via gRPC metadata
- Export ke Jaeger (bisa jalan di Docker)
- Tampilkan screenshot trace di README

#### 3. **Circuit Breaker Pattern**
```
Effort: 1 hari
Impact: Portfolio value +1 poin
```
- Implementasi circuit breaker di broker saat call ke downstream services
- Gunakan library seperti `sony/gobreaker`
- Ini menunjukkan pemahaman tentang fault tolerance di distributed systems

#### 4. **Proper Error Handling**
```
Effort: 1 hari
Impact: Code quality +1.5 poin
```
```go
// Buat custom error types
type AppError struct {
    Code    string
    Message string
    Status  int
}

var (
    ErrStockInsufficient = &AppError{Code: "ORDER_001", Message: "insufficient stock", Status: 400}
    ErrPaymentCompleted  = &AppError{Code: "PAY_001", Message: "payment already completed", Status: 409}
)
```

#### 5. **Fix Database Schema**
```
Effort: 0.5 hari
Impact: Credibility +1 poin
```
- Ganti `DOUBLE PRECISION` → `DECIMAL(12,2)` untuk harga
- Tambahkan indexes pada foreign keys dan frequently queried columns
- Tambahkan `deleted_at` untuk soft delete
- Pertimbangkan UUID untuk distributed-friendly IDs

#### 6. **Graceful Shutdown**
```
Effort: 0.5 hari
Impact: Production readiness +1 poin
```
- Handle `SIGTERM` dan `SIGINT`
- Drain Kafka consumer, close DB connections, stop gRPC server gracefully

### Go Microservice — Nice to Have

| Fitur | Effort | Impact |
|-------|--------|--------|
| Prometheus metrics + Grafana dashboard | 1-2 hari | 🟢 Tinggi |
| API versioning (`/v1/`, `/v2/`) | 0.5 hari | 🟡 Sedang |
| Health check endpoint yang proper (readiness vs liveness) | 0.5 hari | 🟡 Sedang |
| Helm Charts (replace raw YAML) | 1 hari | 🟡 Sedang |
| GitHub Actions (replace/complement Jenkins) | 1 hari | 🟡 Sedang |
| Rate limit per-endpoint (bukan hanya per-user) | 0.5 hari | 🟢 Tinggi |
| Dead letter queue untuk Kafka | 1 hari | 🟢 Tinggi |
| Stock rollback kalau Payment gagal (complete Saga) | 1 hari | 🟢 Tinggi |

### Express.js API — Apa yang Harus Dilakukan

> [!IMPORTANT]
> **Opsi 1: Upgrade drastis** — Tambahkan domain complexity (e.g., jadikan e-commerce mini dengan product, cart, order, caching, search). Tapi ini akan overlap dengan Go project.
>
> **Opsi 2 (Rekomendasi): Fokus di Go project saja.** Express project ini sudah cukup sebagai bukti bahwa kamu bisa Node.js, tapi jangan jadikan *highlight* portfolio. Cukup taruh di GitHub dengan README yang rapi dan mention "learning project" di deskripsi.

Jika tetap ingin keep Express project, minimal tambahkan:
1. **Dockerfile + docker-compose.yaml**
2. **Redis caching** di post listing (kamu sudah punya Redis connection)
3. **Pagination yang proper** (cursor-based)
4. **Structured logging** (winston/pino)
5. **Swagger/OpenAPI documentation**
6. **Test framework** (Jest + Supertest)

---

## 📋 Prioritas Action Items

```
🔴 CRITICAL (Lakukan minggu ini)
├── 1. Tulis unit tests untuk Go microservice (minimal service layer)
├── 2. Fix Jenkinsfile (hapus/perbaiki copy-paste currency-exchange)
└── 3. Fix schema: DOUBLE PRECISION → DECIMAL untuk harga

🟡 HIGH (Lakukan bulan ini)
├── 4. Tambahkan OpenTelemetry distributed tracing
├── 5. Implementasi circuit breaker di broker
├── 6. Complete Saga pattern (stock rollback on payment failure)
└── 7. Dead letter queue untuk Kafka failed messages

🟢 NICE TO HAVE
├── 8. Prometheus + Grafana monitoring
├── 9. Helm Charts
└── 10. API documentation (Swagger)
```

---

## 💡 Verdict Akhir

**Go Microservice** sudah di level yang **jauh di atas rata-rata** proyek portfolio mahasiswa/junior developer. Dengan beberapa perbaikan (terutama testing dan observability), proyek ini bisa jadi **centerpiece portfolio** yang sangat kuat.

**Express API** berfungsi sebagai bukti versatility (kamu bisa Go dan Node.js), tapi jangan taruh di posisi utama portfolio.

**Kalau saya interviewer dan melihat Go Microservice ini, pertanyaan pertama saya:**
1. *"Bagaimana kamu handle kalau Payment service down setelah order created?"* → Stock rollback / Saga compensation
2. *"Mana unit test-nya?"* → ❌ Ini yang harus kamu siapkan
3. *"Bagaimana kamu debug masalah di production lintas service?"* → OpenTelemetry/Jaeger
4. *"Apa yang terjadi kalau Kafka consumer gagal process message?"* → Dead letter queue

Siapkan jawaban teknis untuk pertanyaan-pertanyaan ini, dan kamu akan very well prepared. 🎯
