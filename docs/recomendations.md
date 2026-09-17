# Análisis Arquitectónico y Recomendaciones — Go-inventarie

> **Estado: Todas las recomendaciones de este documento fueron aplicadas.**
> Ver `docs/architecture.md` para la documentación de la arquitectura actual.

## 1. Estructura General del Proyecto

```
cmd/
  main.go              <- punto de entrada
  api/api.go           <- placeholder (pendiente de implementar)
internal/
  config/              <- carga de env y errores
  db/                  <- conexión a PostgreSQL (GORM)
  middleware/auth/     <- hashing + JWT (JWT pendiente de implementar)
  services/users/      <- handler, routes, store
  types/               <- structs compartidos
pkg/
  utils/utils.go       <- placeholder (pendiente de implementar)
  middleware/cache.go  <- placeholder (pendiente de implementar)
test/
  database_test.go     <- test de conexión
```

**Opinión general:** La estructura sigue el [Standard Go Project Layout](https://github.com/golang-standards/project-layout). El uso de `cmd/`, `internal/` y `pkg/` es correcto y demuestra planificación a futuro: los archivos placeholder (`cmd/api`, `pkg/utils`, `pkg/middleware/cache`, `internal/middleware/auth/jwt`) indican una hoja de ruta clara de funcionalidades por implementar. Los puntos a mejorar se centran en el código ya escrito: acoplamiento entre capas, convenciones de nomenclatura y seguridad.

---

## 2. Problemas Detectados

### 2.1 Falta de interfaces (acoplamiento fuerte)
- `Handler` depende directamente del tipo concreto `*Store`.
- `Store` depende directamente de `*gorm.DB`.

Esto impide:
- Unit testing con mocks.
- Reemplazar la capa de persistencia.
- Cumplir con el principio de inversión de dependencias (DIP).

### 2.2 Tipos definidos con nombres inconsistentes
- `types.Users` debería ser `User` (singular, convención Go).
- `ResponseData` no comunica qué representa; debería ser `UserResponse` o `UserDTO`.
- `UsersCreate` → `CreateUserRequest`, `UsersUpdate` → `UpdateUserRequest`.

### 2.3 Lógica de hashing dentro del Store
- `Store.CreateUser` llama a `auth.HashPassword`. El Store no debería conocer detalles de autenticación.
- El hashing debería ocurrir en una capa de servicio intermedia o en el handler, no en la capa de acceso a datos.

### 2.4 Manejo de errores
- Los errores se devuelven como texto plano en español e inglés mezclados ("error al traer usuarios", "invalid request body").
- No hay estructura de error estandarizada (código, mensaje, detalles).
- `GetUserByID` retorna `(error, *types.ResponseData)` — el error debería ir **siempre al final** según convención Go.

### 2.5 Ausencia de validación de input
- No se validan campos obligatorios, formato de email, longitud de password, etc.
- Se debería usar un validador como `go-playground/validator` o validación manual.

### 2.6 Configuración incompleta
- `Config` solo tiene `DatabaseURL`. Faltan: `Port`, `JWTSecret`, `LogLevel`, `Environment`, timeouts.
- El puerto `:8080` está hardcodeado en `main.go`.
- La variable de entorno se llama `URLDATABASE` pero el error dice `DATABASE_URL` — inconsistencia.

### 2.7 Sin migraciones de base de datos
- No hay auto-migrate de GORM ni sistema de migraciones (golang-migrate, goose, atlas).
- Los esquemas de DB no están versionados.

### 2.8 Sin middleware global
- No hay logging de requests, recovery de panics, CORS, ni rate limiting.
- Nota: los archivos `pkg/middleware/cache.go` y `internal/middleware/auth/jwt.go` son placeholders pendientes de implementar, lo cual es una buena planificación.

### 2.9 Test insuficientes
- Solo existe `TestDatabaseConnection`, que es un test de integración que requiere una DB real.
- No hay tests unitarios de handlers, store ni servicios.
- No hay uso de `httptest` para probar endpoints.

### 2.10 Seguridad
- El password se incluye en el JSON de respuesta de `CreateUser` (se retorna `types.UsersCreate` con el campo password).
- No hay HTTPS ni configuración de timeouts en el servidor HTTP.
- `bcrypt.DefaultCost` es aceptable pero debería ser configurable.
- Nota: la autenticación JWT está contemplada como placeholder en `internal/middleware/auth/jwt.go`, pendiente de implementar.

### 2.11 `main.go` no usa graceful shutdown
- `http.ListenAndServe` se ejecuta sin manejo de señales OS para cierre elegante.

### 2.12 `DeleteUser` no verifica si el usuario existe
- `gorm.Delete` no retorna error si el registro no existe; se debería verificar `RowsAffected`.

### 2.13 compose.yml
- El volumen usa `post_data` pero el path es `/var/lib/postgresql` en lugar de `/var/lib/postgresql/data` (el estándar de Postgres).
- No hay healthcheck definido para el servicio de DB.
- Falta la variable `URLDATABASE` que necesita la app.

### 2.14 .env.example incompleto
- No incluye `URLDATABASE`, que es la variable que realmente usa la aplicación.

---

## 3. Arquitectura Recomendada

Se propone una arquitectura en capas con inyección de dependencias:

```
cmd/
  api/
    server.go          <- configuración del servidor HTTP
  main.go              <- bootstrap, composición root

internal/
  config/
    config.go          <- struct Config completo
    errors.go
  db/
    postgres.go        <- conexión + migraciones
  platform/
    middleware/         <- logging, recovery, CORS, auth
  user/
    handler.go         <- handlers HTTP
    service.go         <- lógica de negocio
    store.go           <- acceso a datos
    repository.go      <- interfaz Repository
    dto.go             <- request/response types
    routes.go          <- registro de rutas
  auth/
    jwt.go
    hashing.go
    middleware.go       <- middleware de autenticación

pkg/
  httperr/             <- helper para errores HTTP estructurados
  validator/           <- validación de input
```

### 3.1 Inyección de dependencias

```go
// internal/user/repository.go
type Repository interface {
    GetAll(ctx context.Context) ([]User, error)
    GetByID(ctx context.Context, id int64) (*User, error)
    Create(ctx context.Context, user *CreateUserRequest) (*User, error)
    Update(ctx context.Context, id int64, user *UpdateUserRequest) error
    Delete(ctx context.Context, id int64) error
}
```

```go
// internal/user/service.go
type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}
```

```go
// internal/user/handler.go
type Handler struct {
    service *Service
}
```

### 3.2 Errores estructurados

```go
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}
```

### 3.3 Graceful shutdown

```go
srv := &http.Server{
    Addr:         ":" + cfg.Port,
    Handler:      mux,
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
    IdleTimeout:  60 * time.Second,
}

go func() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    srv.Shutdown(ctx)
}()
```

### 3.4 Validación con struct tags

```go
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```

### 3.5 Migraciones con GORM AutoMigrate

```go
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(&user.User{})
}
```

---

## 4. Resumen de Prioridades

| Prioridad | Acción | Impacto |
|-----------|--------|---------|
| **Alta** | Definir interfaces (Repository) para desacoplar capas | Testeabilidad, mantenibilidad |
| **Alta** | Agregar validación de input en todos los endpoints | Seguridad, robustez |
| **Alta** | No retornar passwords en respuestas (CreateUser) | Seguridad |
| **Alta** | Estandarizar manejo de errores (estructura + idioma) | Consistencia, debugging |
| **Alta** | Agregar variables de entorno faltantes (PORT, JWT_SECRET) | Configurabilidad |
| **Alta** | Corregir orden de retorno en `GetUserByID` → `(*ResponseData, error)` | Convención Go |
| **Media** | Implementar middleware de logging y recovery | Observabilidad, estabilidad |
| **Media** | Agregar graceful shutdown | Producción |
| **Media** | Implementar sistema de migraciones | Evolución del esquema |
| **Media** | Renombrar tipos a singular y con nombres descriptivos | Legibilidad |
| **Media** | Agregar tests unitarios con mocks e httptest | Cobertura |
| **Media** | Mover hashing fuera del Store | Separación de responsabilidades |
| **Media** | Implementar JWT auth (placeholder ya creado) | Seguridad, acceso |
| **Media** | Implementar cache middleware (placeholder ya creado) | Performance |
| **Baja** | Agregar healthcheck en compose.yml | DevOps |
| **Baja** | Completar `.env.example` con `URLDATABASE` | DX |
| **Baja** | Corregir volumen de Postgres a `/var/lib/postgresql/data` | Docker best practices |
| **Baja** | Agregar CORS middleware si hay frontend separado | Integración |

---

## 5. Conclusión

El proyecto demuestra una buena visión arquitectónica desde el inicio: la estructura de directorios sigue las convenciones estándar de Go y los archivos placeholder (`cmd/api`, `pkg/utils`, `pkg/middleware/cache`, `internal/middleware/auth/jwt`) evidencian una hoja de ruta planificada para funcionalidades futuras como autenticación JWT, caching y utilidades compartidas.

Los puntos críticos a resolver en el código actual son la **seguridad** (passwords en respuestas de `CreateUser`, sin validación de input, sin timeouts en el servidor) y la **testeabilidad** (acoplamiento fuerte entre Handler→Store→GORM sin interfaces que permitan mocks).

Recomendación de orden de implementación:
1. Resolver los problemas de alta prioridad en el código existente (interfaces, validación, errores, seguridad).
2. Implementar los placeholders planificados (JWT auth, cache middleware, utils).
3. Agregar la capa de tests unitarios.
4. Refinar la infraestructura (migraciones, graceful shutdown, compose.yml).

Aplicar este orden convertirá el proyecto en una base sólida, segura y escalable siguiendo las mejores prácticas del ecosistema Go.
