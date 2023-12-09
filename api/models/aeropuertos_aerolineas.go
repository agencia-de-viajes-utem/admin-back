package models

type AeropuertoAerolinea struct {
	IDAeropuerto int `gorm:"foreignKey:IDAeropuerto;joinForeignKey:IDAeropuerto"`
	IDAerolinea  int `gorm:"foreignKey:IDAerolinea;joinForeignKey:IDAerolinea"`
}

func (AeropuertoAerolinea) TableName() string {
	return "aeropuertos_aerolineas"
}
