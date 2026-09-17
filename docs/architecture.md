# Arquitectura del Proyecto — Go-inventarie

## 1. Estructura actual del proyecto

```
cmd/
  main.go                          Punto de entrada: carga config, conecta DB, lanza servidor
  api/
    server.go                      Configura el servidor HTTP, middlewares y graceful shutdown

internal/
  config/
    env.go                         Struct Config + función Load() que lee variables de entorno
    errors.go                      Errores sentinel del paquete config

  db/
    config.go                      Conexión a PostgreSQL con GORM + AutoMigrate

  types/
    types_users.go                 Modelos: Users, UsersCreate, UsersUpdate, ResponseData

  middleware/
    logging.go                     Middleware que loggea método, path, status y duración
    recovery.go                    Middleware que captura panics y retorna 500
    cors.go                        Middleware que agrega headers CORS
    auth/
      hashing.go                   HashPassword con bcrypt
      jwt.go                       GenerateToken, ValidateToken, Middleware JWT, GetUserID

  services/
    user/                          <-- Módulo de usuarios (lógica de negocio)
      errors.go                    Errores del dominio: ErrUserNotFound, ErrEmailAlreadyTaken
      repository.go                Interfaz Repository (contrato que Store implementa)
      store.go                     Implementación de Repository con GORM (capa de datos)
      logic.go                     Lógica de negocio: orquesta repo + hashing + mapeo a ResponseData
      handler.go                   Handlers HTTP: parsean request, llaman al service, escriben response
      routes.go                    Registra las rutas en el mux

pkg/
  httperr/
    httperr.go                     Helpers para responder errores HTTP en JSON estructurado

test/
  database_test.go                 Test de integración: verifica conexión a PostgreSQL
  mock_repository.go               Mock en memoria de la interfaz Repository
  service_test.go                  Tests unitarios del Service (7 tests)
  handler_test.go                  Tests unitarios del Handler con httptest (7 tests)
```

---

## 2. Qué hace cada capa (flujo de una request)

El flujo de una petición sigue este camino:

```
Request HTTP
    │
    ▼
┌─────────────┐
│  Middlewares │  Recovery → Logging → CORS → (JWT si la ruta lo requiere)
└──────┬──────┘
       ▼
┌─────────────┐
│   Handler    │  Parsea el request, valida input, llama al Service, escribe response
└──────┬──────┘
       ▼
┌─────────────┐
│   Service    │  Lógica de negocio: hashea passwords, orquesta llamadas, mapea a DTOs
└──────┬──────┘
       ▼
┌─────────────┐
│   Store      │  Acceso a datos con GORM. Implementa la interfaz Repository
└──────┬──────┘
       ▼
┌─────────────┐
│  PostgreSQL  │
└─────────────┘
```

### ¿Por qué separar en capas?

**Handler** solo sabe de HTTP: status codes, headers, JSON. No sabe de base de datos ni de hashing.

**Service** sabe de lógica de negocio: hashear un password antes de guardarlo, mapear un `Users` a un `ResponseData` (sin password). No sabe de HTTP ni de SQL.

**Store** solo sabe de SQL/GORM: hace queries, no valida, no hashea, no decide lógica.

Esto permite:
- Cambiar la base de datos (Postgres → Mongo) sin tocar handlers ni services.
- Testear el handler sin necesitar una DB real (usando mocks).
- Testear el service sin necesitar HTTP (pasándole un mock repository).

---

## 3. Inyección de dependencias — ¿Cómo se conectan las capas?

En `cmd/api/server.go` se "arma" todo:

```go
userStore := user.NewStore(db)        // Store necesita *gorm.DB
userService := user.NewService(userStore) // Service necesita algo que cumpla Repository
userHandler := user.NewHandler(userService) // Handler necesita *Service
```

La clave es la **interfaz Repository** en `repository.go`:

```go
type Repository interface {
    GetAll(ctx context.Context) ([]types.Users, error)
    GetByID(ctx context.Context, id int64) (*types.Users, error)
    Create(ctx context.Context, user *types.Users) error
    Update(ctx context.Context, id int64, updates map[string]interface{}) error
    Delete(ctx context.Context, id int64) error
}
```

`Store` implementa esta interfaz. Pero en tests, `MockRepository` (en `test/mock_repository.go`) también la implementa en memoria. Así podemos testear el Service sin tocar la DB real.

---

## 4. Qué cambió respecto a la versión original del proyecto

### 4.1 Se agregó la capa Service
Antes: `Handler → Store` directamente. El handler tenía que saber de hashing, de mapeo de tipos, etc.

Ahora: `Handler → Service → Store`. Cada uno con una sola responsabilidad.

### 4.2 Se agregó la interfaz Repository
Antes: Handler dependía de `*Store` (tipo concreto). No se podía mockear.

Ahora: Service depende de la interfaz `Repository`. Store la implementa, pero también un mock.

### 4.3 El hashing se movió de Store a Service
Antes: `Store.CreateUser` llamaba a `auth.HashPassword`. La capa de datos sabía de seguridad.

Ahora: `Service.Create` hashea el password y luego le pasa el `User` ya hasheado al Store.

### 4.4 Los tipos se organizaron y protegieron
Antes: Un solo struct `Users` servía para todo (modelo DB, request, response). El password se exponía en las respuestas.

Ahora hay tipos distintos para cada propósito, todos en `internal/types/types_users.go`:
- `types.Users` → modelo de dominio (lo que vive en la DB). Password con `json:"-"` para que nunca se serialice.
- `types.UsersCreate` → lo que el cliente envía para crear un usuario (con tags de validación).
- `types.UsersUpdate` → lo que el cliente envía para actualizar (incluye campo `Password` opcional).
- `types.ResponseData` → lo que se le devuelve al cliente (sin password).

### 4.5 Errores estructurados
Antes: `http.Error(w, "error al traer usuarios", 500)` → texto plano, mezcla de idiomas.

Ahora: `httperr.Internal(w, "failed to fetch users")` → JSON estructurado con código, mensaje y detalles.

### 4.6 Validación de input
Antes: No había validación. Cualquier dato pasaba.

Ahora: `validator.Struct(req)` valida los tags definidos en `types.UsersCreate` y `types.UsersUpdate` (formato de email, largo mínimo de password, campos requeridos, etc.).

### 4.7 Middlewares globales
Antes: No había nada entre la request y el handler.

Ahora: Todas las requests pasan por Recovery (captura panics) → Logging (loggea cada request) → CORS (headers cross-origin). Rutas PUT y DELETE además pasan por JWT.

### 4.8 Graceful shutdown
Antes: `http.ListenAndServe` se mataba abruptamente con Ctrl+C.

Ahora: El servidor escucha señales SIGINT/SIGTERM y hace un shutdown elegante, terminando requests en curso antes de cerrar.

### 4.9 AutoMigrate
Antes: Las tablas de la DB se tenían que crear manualmente.

Ahora: GORM crea/actualiza las tablas automáticamente al arrancar basándose en el struct `Users`.

### 4.10 Configuración ampliada
Antes: Solo `URLDATABASE`. Puerto hardcodeado en `:8080`.

Ahora: `URLDATABASE`, `PORT`, `JWT_SECRET`, `ENVIRONMENT`. Todo desde variables de entorno con valores por defecto razonables.

### 4.11 Estructura de servicios reorganizada
Antes: `internal/services/users/service.go` → la palabra "service" aparecía dos veces en el path (carpeta + archivo).

Ahora: `internal/services/user/logic.go` → la carpeta `services/` agrupa todos los módulos de negocio, `user/` es singular (convención Go), y `logic.go` describe lo que hace sin redundar con la carpeta padre.

---

## 5. Rutas y autenticación

| Método | Ruta | Auth | Descripción |
|--------|------|------|-------------|
| GET | `/users` | No | Lista todos los usuarios |
| GET | `/users/{id}` | No | Obtiene un usuario por ID |
| POST | `/users` | No | Crea un usuario nuevo |
| PUT | `/users/{id}` | **JWT** | Actualiza un usuario |
| DELETE | `/users/{id}` | **JWT** | Elimina un usuario |

Para autenticarse, enviar header: `Authorization: Bearer <token>`

---

## 6. Cómo correr el proyecto

```bash
# Levantar la base de datos
make up

# Correr la aplicación
go run cmd/main.go

# Correr todos los tests
go test ./...

# Correr solo los tests unitarios (sin necesitar DB)
go test ./test/ -run "TestHandler|TestService"

# Correr solo el test de conexión a DB
go test ./test/ -run TestDatabaseConnection
```

---

## 7. Variables de entorno

| Variable | Descripción | Default |
|----------|-------------|---------|
| `URLDATABASE` | Connection string de PostgreSQL | *(obligatoria)* |
| `PORT` | Puerto del servidor HTTP | `8080` |
| `JWT_SECRET` | Secreto para firmar tokens JWT | *(obligatoria)* |
| `ENVIRONMENT` | Entorno (`development`, `production`) | `development` |
