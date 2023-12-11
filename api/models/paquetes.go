package models

import "github.com/lib/pq"

type Paquete struct {
	ID                  int            `gorm:"primaryKey;column:id"`
	Nombre              string         `gorm:"column:nombre"`
	Descripcion         string         `gorm:"column:descripcion"`
	PrecioNormal        float64        `gorm:"column:precio_normal"`
	Imagenes            pq.StringArray `gorm:"type:text[];column:imagenes"`
	IDAeropuertoOrigen  int            `gorm:"column:id_aeropuerto_origen"`
	IDAeropuertoDestino int            `gorm:"column:id_aeropuerto_destino"`
	IDAerolinea         int            `gorm:"column:id_aerolinea"`
	HabitacionIDs       []int          `gorm:"-"` // Campo no persistido, solo para la actualización

	// Claves foráneas
	AeropuertoOrigen  Aeropuerto   `gorm:"foreignKey:IDAeropuertoOrigen"`
	AeropuertoDestino Aeropuerto   `gorm:"foreignKey:IDAeropuertoDestino"`
	Habitaciones      []Habitacion `gorm:"many2many:paquetes_habitaciones;foreignKey:ID;joinForeignKey:IDPaquete;References:ID;joinReferences:IDHabitacion"`
	Aerolinea         Aerolinea    `gorm:"foreignKey:IDAerolinea"`
}

func (Paquete) TableName() string {
	return "paquetes"
}
