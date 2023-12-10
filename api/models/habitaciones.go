package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Habitacion struct {
	ID               int            `gorm:"primaryKey;column:id"`
	Nombre           string         `gorm:"column:nombre"`
	Descripcion      string         `gorm:"column:descripcion"`
	Servicios        pq.StringArray `gorm:"type:text[];column:servicios"`
	PrecioNoche      float64        `gorm:"column:precio_noche"`
	Imagenes         pq.StringArray `gorm:"type:text[];column:imagenes"`
	IDTipoHabitacion int            `gorm:"column:id_tipo_habitacion"`
	IDHotel          int            `gorm:"column:id_hotel"`

	// Claves foráneas
	TipoHabitacion TipoHabitacion `gorm:"foreignKey:IDTipoHabitacion"`
	Hotel          *Hotel         `gorm:"foreignKey:IDHotel"`
}

type HabitacionSinHotel struct {
	ID               int            `gorm:"primaryKey;column:id"`
	Nombre           string         `gorm:"column:nombre"`
	Descripcion      string         `gorm:"column:descripcion"`
	Servicios        pq.StringArray `gorm:"type:text[];column:servicios"`
	PrecioNoche      float64        `gorm:"column:precio_noche"`
	Imagenes         pq.StringArray `gorm:"type:text[];column:imagenes"`
	IDTipoHabitacion int            `gorm:"column:id_tipo_habitacion"`

	// Claves foráneas
	TipoHabitacion TipoHabitacion `gorm:"foreignKey:IDTipoHabitacion"`
	Ocupada        bool           `gorm:"column:ocupada, omitempty"`
	IDPaquete      int            `gorm:"column:id_paquete, omitempty"`
}

func (Habitacion) TableName() string {
	return "habitaciones"
}

// HabitacionConPaquete incluye la información de si una habitación está ocupada y, en caso afirmativo, el ID del paquete asociado
type HabitacionConPaquete struct {
	Habitacion
	Ocupada   bool `json:"ocupada"`
	IDPaquete int  `json:"idPaquete,omitempty"` // omitempty para no mostrar si no está ocupada
}

func EstaOcupada(h Habitacion, db *gorm.DB) (HabitacionConPaquete, error) {
	var resultado HabitacionConPaquete
	var relacion struct {
		IDPaquete int
	}

	// Buscar si existe una relación en paquetes_habitaciones
	err := db.Table("paquetes_habitaciones").Select("id_paquete").Where("id_habitacion = ?", h.ID).First(&relacion).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return resultado, err
	}

	resultado.Habitacion = h
	if err == nil { // Si se encuentra un registro
		resultado.Ocupada = true
		resultado.IDPaquete = relacion.IDPaquete
	}

	return resultado, nil
}
