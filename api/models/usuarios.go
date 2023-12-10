package models

type Usuario struct {
	ID              string `gorm:"primaryKey;column:uid"` // Utiliza `gorm:"primaryKey"` si este es tu campo clave
	Nombre          string `gorm:"column:nombre"`
	Apellido        string `gorm:"column:apellido_paterno"`
	SegundoApellido string `gorm:"column:apellido_materno"`
	Email           string `gorm:"column:email"` // `uniqueIndex` si el email debe ser único
	RUT             string `gorm:"column:rut"`
	Fono            string `gorm:"column:telefono"`
	ImgProfile      string `gorm:"column:imagen_perfil"`
	Rol             string `gorm:"column:rol"`
}

// TableName sobrescribe el nombre de la tabla por defecto
func (Usuario) TableName() string {
	return "usuarios"
}
