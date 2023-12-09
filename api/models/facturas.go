package models

type Factura struct {
	ID           int     `gorm:"primaryKey;column:id"`
	IDReserva    int     `gorm:"column:id_reserva"`
	FechaEmision string  `gorm:"column:fecha_emision"` // Puedes usar el tipo de dato apropiado para la fecha según tu base de datos
	Total        float64 `gorm:"column:total"`
	Neto         float64 `gorm:"column:neto"`
	IVA          float64 `gorm:"column:iva"`

	// Clave foránea
	Reserva Reserva `gorm:"foreignKey:IDReserva"`
}

func (Factura) TableName() string {
	return "facturas"
}
