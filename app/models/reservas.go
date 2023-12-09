package models

type Reserva struct {
	ID            int     `gorm:"primaryKey;column:id"`
	UIDUsuario    string  `gorm:"column:uid_usuario"`
	IDPaquete     int     `gorm:"column:id_paquete"`
	FechaInicio   string  `gorm:"column:fecha_inicio"` // Puedes usar el tipo de dato apropiado para la fecha según tu base de datos
	FechaFin      string  `gorm:"column:fecha_fin"`    // Puedes usar el tipo de dato apropiado para la fecha según tu base de datos
	Precio        float64 `gorm:"column:precio"`
	NumNoches     int     `gorm:"column:num_noches"`
	EstadoReserva string  `gorm:"column:estado_reserva"`

	// Claves foráneas
	Paquete Paquete `gorm:"foreignKey:IDPaquete"`
	Usuario Usuario `gorm:"foreignKey:UIDUsuario"`
}

func (Reserva) TableName() string {
	return "reservas"
}
