package models

type Pago struct {
	ID        int     `gorm:"primaryKey;column:id"`
	IDFactura int     `gorm:"column:id_factura"`
	FechaPago string  `gorm:"column:fecha_pago"` // Puedes usar el tipo de dato apropiado para la fecha según tu base de datos
	Monto     float64 `gorm:"column:monto"`

	// Clave foránea
	Factura Factura `gorm:"foreignKey:IDFactura"`
}

func (Pago) TableName() string {
	return "pagos"
}
