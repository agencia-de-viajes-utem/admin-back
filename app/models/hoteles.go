package models

import "github.com/lib/pq"

type Hotel struct {
	ID          int            `gorm:"primaryKey;column:id"`
	Nombre      string         `gorm:"column:nombre"`
	Descripcion string         `gorm:"column:descripcion"`
	Direccion   string         `gorm:"column:direccion"`
	Servicios   pq.StringArray `gorm:"type:text[];column:servicios"`
	Email       string         `gorm:"column:email"`
	Telefono    string         `gorm:"column:telefono"`
	SitioWeb    string         `gorm:"column:sitio_web"`
	Imagenes    pq.StringArray `gorm:"type:text[];column:imagenes"`
	CiudadID    int            `gorm:"column:id_ciudad"`    // Cambia a un campo de tipo int que almacena el ID de la ciudad
	Ciudad      Ciudad         `gorm:"foreignKey:CiudadID"` // Define la relación con Ciudad y la clave externa
}

func (Hotel) TableName() string {
	return "hoteles"
}

// Esta es la que hicimos pero no existe en la BD y no se si está bien la pensé para la respuesta del json

type HotelHabitaciones struct {
	HotelID       int   `gorm:"column:hotel_id"`
	HabitacionIDs []int `gorm:"type:integer[];column:habitacion_ids"`
	Sum           int   `gorm:"column:sum"`
}
