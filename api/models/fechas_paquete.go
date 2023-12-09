package models

type FechaPaquete struct {
	ID           int     `gorm:"primaryKey;column:id"`
	Nombre       string  `gorm:"column:nombre"`
	FechaInicial string  `gorm:"column:fecha_inicial"` // Puedes usar el tipo de dato apropiado para la fecha según tu base de datos
	FechaFinal   string  `gorm:"column:fecha_final"`   // Puedes usar el tipo de dato apropiado para la fecha según tu base de datos
	PrecioOferta float64 `gorm:"column:precio_oferta"`
	IDPaquete    int     `gorm:"column:id_paquete"`

	// Clave foránea
	Paquete Paquete `gorm:"foreignKey:IDPaquete"`
}

func (FechaPaquete) TableName() string {
	return "fechas_paquete"
}
