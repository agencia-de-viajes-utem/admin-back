package models

type PaqueteHabitacion struct {
	IDPaquete    int `gorm:"primaryKey;column:id_paquete;foreignKey"`
	IDHabitacion int `gorm:"primaryKey;unique;column:id_habitacion;foreignKey"`
}

func (PaqueteHabitacion) TableName() string {
	return "paquetes_habitaciones"
}
