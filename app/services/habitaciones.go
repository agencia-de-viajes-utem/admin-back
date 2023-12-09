package services

import (
	"backend-admin/app/models"
	"errors"
	"log"

	"gorm.io/gorm"
)

func GetHabitaciones(db *gorm.DB) ([]models.HabitacionConPaquete, error) {
	var habitaciones []models.Habitacion
	var habitacionesConPaquete []models.HabitacionConPaquete

	result := db.Preload("TipoHabitacion").
		Preload("Hotel").
		Preload("Hotel.Ciudad").
		Preload("Hotel.Ciudad.Pais").
		Find(&habitaciones)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, h := range habitaciones {
		habitacionConPaquete, err := models.EstaOcupada(h, db)
		if err != nil {
			return nil, err // O manejar el error de manera diferente
		}
		habitacionesConPaquete = append(habitacionesConPaquete, habitacionConPaquete)
	}

	return habitacionesConPaquete, nil
}

func GetHabitacionByID(db *gorm.DB, id int) (*models.HabitacionConPaquete, error) {
	var habitacion models.Habitacion

	result := db.Preload("TipoHabitacion").
		Preload("Hotel").
		Preload("Hotel.Ciudad").
		Preload("Hotel.Ciudad.Pais").
		First(&habitacion, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("habitación no encontrada")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	habitacionConPaquete, err := models.EstaOcupada(habitacion, db)
	if err != nil {
		return nil, err
	}

	return &habitacionConPaquete, nil
}

// CreateHabitacion crea una nueva habitación en la base de datos.
func CreateHabitacion(db *gorm.DB, nuevaHabitacion models.Habitacion) error {
	result := db.Create(&nuevaHabitacion)
	if result.Error != nil {
		log.Printf("Error al crear la habitación: %v\n", result.Error)
		return result.Error
	}

	log.Printf("Habitación creada correctamente\n")
	return nil
}

// UpdateHabitacion actualiza una habitación existente en la base de datos.
func UpdateHabitacion(db *gorm.DB, updatedHabitacion *models.Habitacion) error {
	// Asegurarse de que las relaciones existan o sean actualizadas correctamente
	if _, err := GetTipoHabitacionByID(db, updatedHabitacion.IDTipoHabitacion); err != nil {
		return err
	}
	if _, err := GetHotelByID(db, updatedHabitacion.IDHotel); err != nil {
		return err
	}

	result := db.Save(updatedHabitacion)
	if result.Error != nil {
		log.Printf("Error al actualizar la habitación: %v\n", result.Error)
		return result.Error
	}

	log.Printf("Habitación actualizada correctamente\n")
	return nil
}

// DeleteHabitacion elimina una habitación por su ID desde la base de datos.
func DeleteHabitacion(db *gorm.DB, id int) error {
	var habitacion models.Habitacion

	result := db.Delete(&habitacion, id)
	if result.Error != nil {
		log.Printf("Error al eliminar la habitación: %v\n", result.Error)
		return result.Error
	}

	log.Printf("Habitación eliminada correctamente\n")
	return nil
}
