package models

type Pais struct {
	ID         int    `gorm:"primaryKey;column:id"`
	Nombre     string `gorm:"column:nombre"`
	CodigoPais string `gorm:"column:codigo_pais"`
}

// TableName sobrescribe el nombre de la tabla por defecto
func (Pais) TableName() string {
	return "paises"
}
