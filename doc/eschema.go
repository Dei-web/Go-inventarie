package docs

import "time"

// =====================================================================
// ENUMS
// =====================================================================

type PieceState string

const (
	Disponible PieceState = "DISPONIBLE"
	Agotado    PieceState = "AGOTADO"
)

type CitaEstado string

const (
	Asignada   CitaEstado = "ASIGNADA"
	Completada CitaEstado = "COMPLETADA"
	Pendiente  CitaEstado = "PENDIENTE"
	Cancelada  CitaEstado = "CANCELADA"
)

// TipoServicio distingue las dos categorías de servicio: el catálogo
// de taller (Servicio) y las solicitudes de alquiler de flota
// (SolicitudFlota). Ambas quedan bajo el mismo concepto de "servicio"
// gracias a este campo en CategoriaServicio.
type TipoServicio string

const (
	Taller TipoServicio = "TALLER"
	Flota  TipoServicio = "FLOTA"
)

type EstadoSolicitud string

const (
	Procesando EstadoSolicitud = "PROCESANDO"
	Rechazado  EstadoSolicitud = "RECHAZADO"
	Aprobado   EstadoSolicitud = "APROBADO"
	Creado     EstadoSolicitud = "CREADO"
)

// =====================================================================
// SEDES
// =====================================================================

type Sede struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre    string    `gorm:"size:100;not null" json:"nombre"`
	Direccion string    `gorm:"size:150" json:"direccion"`
	Telefono  string    `gorm:"size:20" json:"telefono"`
	CreatedAt time.Time `json:"createdAt"`
}

func (Sede) TableName() string { return "sedes" }

// =====================================================================
// ROLES Y USUARIOS
// =====================================================================

// Rol es una tabla abierta: agregar roles nuevos (Coordinador Regional,
// Coordinador de Flota, Solicitante, Mecánico, Admin...) es solo
// insertar filas, no requiere cambio de esquema.
type Rol struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre      string `gorm:"size:50;uniqueIndex;not null" json:"nombre"`
	Descripcion string `gorm:"size:150" json:"descripcion"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   uint      `gorm:"index" json:"-"`
}

func (Rol) TableName() string { return "roles" }

type Usuario struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre    string    `gorm:"size:100;not null" json:"nombre"`
	Email     string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	RolID     uint      `gorm:"index;not null" json:"rolId"`
	SedeID    uint      `gorm:"index;not null" json:"sedeId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt uint      `gorm:"index" json:"-"`
}

func (Usuario) TableName() string { return "usuarios" }

// =====================================================================
// INVENTARIO
// =====================================================================

type CategoriaPieza struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre    string `gorm:"size:50;not null" json:"nombre"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt uint   `gorm:"index" json:"-"`
}

func (CategoriaPieza) TableName() string { return "categorias_pieza" }

type Pieza struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre    string     `gorm:"size:100;not null" json:"nombre"`
	Precio    float64    `gorm:"type:decimal(15,2)" json:"precio"`
	Stock     int        `json:"stock"`
	Estado    PieceState `gorm:"type:varchar(20);default:DISPONIBLE" json:"estado"`
	CategoriaID uint      `gorm:"index;not null" json:"categoriaId"`
	ProveedorID uint      `gorm:"index;not null" json:"proveedorId"`
	SedeID      uint      `gorm:"index;not null" json:"sedeId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   uint      `gorm:"index" json:"-"`
}

func (Pieza) TableName() string { return "piezas" }

// =====================================================================
// SERVICIOS (categoría común para TALLER y FLOTA)
// =====================================================================

type CategoriaServicio struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre    string         `gorm:"size:50;not null" json:"nombre"`
	Tipo      TipoServicio   `gorm:"type:varchar(20);not null" json:"tipo"` // TALLER o FLOTA
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt uint           `gorm:"index" json:"-"`
}

func (CategoriaServicio) TableName() string { return "categorias_servicio" }

// Servicio es el catálogo de taller (ítems con precio fijo). Sus
// CategoriaID deben apuntar a una CategoriaServicio con Tipo=TALLER.
type Servicio struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre      string    `gorm:"size:100;not null" json:"nombre"`
	Descripcion string    `json:"descripcion"`
	Precio      float64   `gorm:"type:decimal(15,2)" json:"precio"`
	CategoriaID uint      `gorm:"index;not null" json:"categoriaId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   uint      `gorm:"index" json:"-"`
}

func (Servicio) TableName() string { return "servicios" }

// =====================================================================
// PROVEEDORES
// =====================================================================

type Proveedor struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre    string    `gorm:"size:100;not null" json:"nombre"`
	Telefono  string    `gorm:"size:20" json:"telefono"`
	Email     string    `gorm:"size:100" json:"email"`
	Direccion string    `gorm:"size:150" json:"direccion"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt uint      `gorm:"index" json:"-"`
}

func (Proveedor) TableName() string { return "proveedores" }

// =====================================================================
// CITAS (taller)
// =====================================================================

type Cita struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Fecha           time.Time  `json:"fecha"`
	Estado          CitaEstado `gorm:"type:varchar(20);default:PENDIENTE" json:"estado"`
	ClienteNombre   string     `gorm:"size:100;not null" json:"clienteNombre"`
	ClienteTelefono string     `gorm:"size:20" json:"clienteTelefono"`
	UsuarioID       uint       `gorm:"index;not null" json:"usuarioId"` // empleado asignado
	ServicioID      uint       `gorm:"index;not null" json:"servicioId"`
	SedeID          uint       `gorm:"index;not null" json:"sedeId"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       uint       `gorm:"index" json:"-"`
}

func (Cita) TableName() string { return "citas" }

// =====================================================================
// VEHÍCULOS Y DISPONIBILIDAD (flota)
// =====================================================================

type Vehiculo struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Marca      string    `gorm:"size:100;not null" json:"marca"`
	Modelo     string    `gorm:"size:100;not null" json:"modelo"`
	AnioModelo int       `json:"anioModelo"`
	Estado     string    `gorm:"size:50" json:"estado"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	DeletedAt  uint      `gorm:"index" json:"-"`
}

func (Vehiculo) TableName() string { return "vehiculos" }

// VehiculoDisponible marca qué vehículos están libres para asignar y
// por qué (ej: "devuelto de préstamo", "recién adquirido").
type VehiculoDisponible struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	VehiculoID       uint      `gorm:"uniqueIndex;not null" json:"vehiculoId"`
	MotivoDisponible string    `json:"motivoDisponible"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	DeletedAt        uint      `gorm:"index" json:"-"`
}

func (VehiculoDisponible) TableName() string { return "vehiculos_disponibles" }

// =====================================================================
// CENTROS DE COSTO Y COORDINADORES DE FLOTA
// =====================================================================

type CentroCosto struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CoordinadorAsignado uint      `gorm:"index;not null" json:"coordinadorAsignado"` // FK -> Usuario
	NombreEmpresa       string    `gorm:"size:255;not null" json:"nombreEmpresa"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
	DeletedAt           uint      `gorm:"index" json:"-"`
}

func (CentroCosto) TableName() string { return "centros_costo" }

// CoordinadorFlota liga un usuario (rol "Coordinador de Flota") a un
// centro de costo que administra.
type CoordinadorFlota struct {
	ID            uint `gorm:"primaryKey;autoIncrement" json:"id"`
	CoordinadorID uint `gorm:"index;not null" json:"coordinadorId"` // FK -> Usuario
	CentroCostoID uint `gorm:"index;not null" json:"centroCostoId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	DeletedAt     uint      `gorm:"index" json:"-"`
}

func (CoordinadorFlota) TableName() string { return "coordinadores_flota" }

// VehiculoAsignado: qué vehículo quedó bajo el control de qué
// coordinador de flota.
type VehiculoAsignado struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CoordinadorFlotaID uint      `gorm:"index;not null" json:"coordinadorFlotaId"`
	VehiculoID         uint      `gorm:"index;not null" json:"vehiculoId"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	DeletedAt          uint      `gorm:"index" json:"-"`
}

func (VehiculoAsignado) TableName() string { return "vehiculos_asignados" }

// =====================================================================
// SOLICITUDES DE SERVICIO DE FLOTA (ALQUILER) Y SU HISTORIAL DE ESTADO
// =====================================================================

// SolicitudFlota es la solicitud de alquiler de vehículo: cuándo, con
// qué requisitos, y quién la creó / quién la regula regionalmente. Su
// CategoriaServicioID debe apuntar a una CategoriaServicio con
// Tipo=FLOTA.
type SolicitudFlota struct {
	ID                    uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FechaInicio           time.Time `gorm:"type:date;not null" json:"fechaInicio"`
	FechaFin              time.Time `gorm:"type:date;not null" json:"fechaFin"`
	RequisitosVehiculo    string    `gorm:"not null" json:"requisitosVehiculo"`
	CoordinadorRegionalID uint      `gorm:"index;not null" json:"coordinadorRegionalId"` // FK -> Usuario
	CreadoPorID           uint      `gorm:"index;not null" json:"creadoPorId"`           // FK -> Usuario
	CategoriaServicioID   uint      `gorm:"index;not null" json:"categoriaServicioId"`   // FK -> CategoriaServicio (Tipo=FLOTA)
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

func (SolicitudFlota) TableName() string { return "solicitudes_flota" }

// EstadoSolicitudFlota es el historial de cambios de estado de una
// solicitud: cada cambio queda como registro nuevo, nunca se
// sobreescribe el anterior, para trazabilidad completa.
type EstadoSolicitudFlota struct {
	ID               uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	SolicitudID      uint            `gorm:"index;not null" json:"solicitudId"`
	Comentarios      string          `gorm:"not null" json:"comentarios"`
	Estado           EstadoSolicitud `gorm:"type:varchar(20);default:CREADO" json:"estado"`
	ActualizadoPorID uint            `gorm:"index;not null" json:"actualizadoPorId"` // FK -> Usuario
	ActualizadoAt    time.Time       `json:"actualizadoAt"`
}

func (EstadoSolicitudFlota) TableName() string { return "estados_solicitud_flota" }
