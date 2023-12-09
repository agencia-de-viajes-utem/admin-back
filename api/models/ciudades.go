package models

type Ciudad struct {
	ID     int    `gorm:"primaryKey;column:id"`
	Nombre string `gorm:"column:nombre"`
	IDPais int    `gorm:"column:id_pais"`
	Pais   Pais   `gorm:"foreignKey:IDPais"`
}

func (Ciudad) TableName() string {
	return "ciudades"
}
