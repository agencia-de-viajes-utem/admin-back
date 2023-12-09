package models

type Aeropuerto struct {
	ID       int    `gorm:"primaryKey;column:id"`
	Nombre   string `gorm:"column:nombre"`
	IDCiudad int    `gorm:"column:id_ciudad"`
	Ciudad   Ciudad `gorm:"foreignKey:IDCiudad"`
}

func (Aeropuerto) TableName() string {
	return "aeropuertos"
}
