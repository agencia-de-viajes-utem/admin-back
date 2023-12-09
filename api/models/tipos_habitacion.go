package models

type TipoHabitacion struct {
	ID        int    `gorm:"primaryKey;column:id"`
	Nombre    string `gorm:"column:nombre"`
	Capacidad int    `gorm:"column:capacidad"`
}

func (TipoHabitacion) TableName() string {
	return "tipos_habitacion"
}
