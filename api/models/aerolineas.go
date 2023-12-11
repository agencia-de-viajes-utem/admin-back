package models

type Aerolinea struct {
	ID          int          `gorm:"primaryKey;column:id"`
	Nombre      string       `gorm:"column:nombre"`
	Imagen      string       `gorm:"column:imagen"`
	Aeropuertos []Aeropuerto `gorm:"many2many:aeropuertos_aerolineas;joinForeignKey:id_aerolinea;joinReferences:id_aeropuerto"`
}

func (Aerolinea) TableName() string {
	return "aerolineas"
}
