package utils

import (
	"backend-admin/api/models" // Asegúrate de reemplazar esto con el path correcto a tu paquete de modelos
)

func TransformarAPaqueteResumido(paquete models.Paquete) models.PaqueteResumido {
	paqueteResumido := models.PaqueteResumido{
		ID:           paquete.ID,
		Nombre:       paquete.Nombre,
		Descripcion:  paquete.Descripcion,
		PrecioNormal: paquete.PrecioNormal,
		Imagenes:     paquete.Imagenes,
		Origen:       TransformarAeropuertoResumido(paquete.AeropuertoOrigen),
		Destino:      TransformarAeropuertoResumido(paquete.AeropuertoDestino),
	}

	if len(paquete.Habitaciones) > 0 && paquete.Habitaciones[0].Hotel != nil && paquete.Habitaciones[0].Hotel.ID != 0 {
		paqueteResumido.Hotel = TransformarHotelResumido(paquete.Habitaciones[0].Hotel)
	} else {
		paqueteResumido.Hotel = nil // O nil, dependiendo de cómo quieras manejarlo
	}

	for _, habitacion := range paquete.Habitaciones {
		habitacionResumida := TransformarHabitacionResumida(habitacion)
		paqueteResumido.Habitaciones = append(paqueteResumido.Habitaciones, habitacionResumida)
	}

	return paqueteResumido
}

// TransformarAeropuertoResumido convierte un Aeropuerto en su versión resumida
func TransformarAeropuertoResumido(aeropuerto models.Aeropuerto) models.AeropuertoResumido {
	return models.AeropuertoResumido{
		ID:     aeropuerto.ID,
		Nombre: aeropuerto.Nombre,
		Ciudad: TransformarCiudadResumida(aeropuerto.Ciudad),
	}
}

// TransformarCiudadResumida convierte una Ciudad en su versión resumida
func TransformarCiudadResumida(ciudad models.Ciudad) models.CiudadResumida {
	return models.CiudadResumida{
		ID:     ciudad.ID,
		Nombre: ciudad.Nombre,
		Pais:   TransformarPaisResumido(ciudad.Pais),
	}
}

// TransformarPaisResumido convierte un Pais en su versión resumida
func TransformarPaisResumido(pais models.Pais) models.PaisResumido {
	return models.PaisResumido{
		ID:     pais.ID,
		Nombre: pais.Nombre,
	}
}

// TransformarHotelResumido convierte un Hotel en su versión resumida
func TransformarHotelResumido(hotel *models.Hotel) *models.HotelResumido {
	if hotel == nil {
		return nil
	}

	HotelResumido := models.HotelResumido{
		ID:          hotel.ID,
		Nombre:      hotel.Nombre,
		Descripcion: hotel.Descripcion,
		Direccion:   hotel.Direccion,
		Ciudad:      TransformarCiudadResumida(hotel.Ciudad),
		Servicios:   hotel.Servicios,
		Contacto: models.Contacto{
			Email:    hotel.Email,
			Telefono: hotel.Telefono,
			SitioWeb: hotel.SitioWeb,
		},
		Imagenes: hotel.Imagenes,
	}

	return &HotelResumido

}

// TransformarHabitacionResumida convierte una Habitacion en su versión resumida
func TransformarHabitacionResumida(habitacion models.Habitacion) models.HabitacionResumida {
	return models.HabitacionResumida{
		ID:             habitacion.ID,
		Nombre:         habitacion.Nombre,
		Descripcion:    habitacion.Descripcion,
		Servicios:      habitacion.Servicios,
		PrecioNoche:    habitacion.PrecioNoche,
		Imagenes:       habitacion.Imagenes,
		TipoHabitacion: TransformarTipoHabitacionResumida(habitacion.TipoHabitacion),
	}
}

// TransformarTipoHabitacionResumida convierte un TipoHabitacion en su versión resumida
func TransformarTipoHabitacionResumida(tipo models.TipoHabitacion) models.TipoHabitacionResumida {
	return models.TipoHabitacionResumida{
		ID:        tipo.ID,
		Nombre:    tipo.Nombre,
		Capacidad: tipo.Capacidad,
	}
}
