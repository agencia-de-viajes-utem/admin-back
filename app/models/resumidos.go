package models

import "github.com/lib/pq"

// PaqueteResumido representa una versión simplificada de un Paquete para el frontend
type PaqueteResumido struct {
	ID           int                  `json:"ID"`
	Nombre       string               `json:"Nombre"`
	Descripcion  string               `json:"Descripcion"`
	PrecioNormal float64              `json:"PrecioNormal"`
	Imagenes     pq.StringArray       `json:"Imagenes"`
	Origen       AeropuertoResumido   `json:"Origen"`
	Destino      AeropuertoResumido   `json:"Destino"`
	Habitaciones []HabitacionResumida `json:"Habitaciones,omitempty"`
	Hotel        *HotelResumido       `json:"Hotel,omitempty"`
}

// AeropuertoResumido representa una versión simplificada de un Aeropuerto
type AeropuertoResumido struct {
	ID     int            `json:"ID"`
	Nombre string         `json:"Nombre"`
	Ciudad CiudadResumida `json:"Ciudad"`
}

// CiudadResumida representa una versión simplificada de una Ciudad
type CiudadResumida struct {
	ID     int          `json:"ID"`
	Nombre string       `json:"Nombre"`
	Pais   PaisResumido `json:"Pais"`
}

// PaisResumido representa una versión simplificada de un Pais
type PaisResumido struct {
	ID     int    `json:"ID"`
	Nombre string `json:"Nombre"`
}

// HotelResumido representa una versión simplificada de un Hotel
type HotelResumido struct {
	ID          int            `json:"ID"`
	Nombre      string         `json:"Nombre"`
	Descripcion string         `json:"Descripcion"`
	Direccion   string         `json:"Direccion"`
	Ciudad      CiudadResumida `json:"Ciudad"`
	Servicios   pq.StringArray `json:"Servicios"`
	Contacto    Contacto       `json:"Contacto"`
	Imagenes    pq.StringArray `json:"Imagenes"`
}

// Contacto representa la información de contacto de un Hotel
type Contacto struct {
	Email    string `json:"Email"`
	Telefono string `json:"Telefono"`
	SitioWeb string `json:"SitioWeb"`
}

// HabitacionResumida representa una versión simplificada de una Habitación
type HabitacionResumida struct {
	ID             int                    `json:"ID"`
	Nombre         string                 `json:"Nombre"`
	Descripcion    string                 `json:"Descripcion"`
	Servicios      pq.StringArray         `json:"Servicios"`
	PrecioNoche    float64                `json:"PrecioNoche"`
	Imagenes       pq.StringArray         `json:"Imagenes"`
	TipoHabitacion TipoHabitacionResumida `json:"TipoHabitacion"`
}

// TipoHabitacionResumida representa una versión simplificada de un TipoHabitacion
type TipoHabitacionResumida struct {
	ID        int    `json:"ID"`
	Nombre    string `json:"Nombre"`
	Capacidad int    `json:"Capacidad"`
}
