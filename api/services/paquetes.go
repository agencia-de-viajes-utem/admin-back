// services/paquete_service.go
package services

import (
	"backend-admin/api/models"
	"errors"
	"log"
	"reflect"

	"gorm.io/gorm"
)

// GetPaquetes obtiene todos los paquetes desde la base de datos.
func GetPaquetes(db *gorm.DB) ([]models.Paquete, error) {
	var paquetes []models.Paquete

	result := db.Preload("AeropuertoOrigen.Ciudad.Pais").
		Preload("AeropuertoDestino.Ciudad.Pais").
		Preload("Habitaciones.Hotel.Ciudad.Pais").
		Preload("Habitaciones.TipoHabitacion").
		Preload("Aerolinea").
		Find(&paquetes)
	if result.Error != nil {
		return nil, result.Error
	}

	return paquetes, nil
}

// GetPaqueteByID obtiene un paquete por su ID desde la base de datos.
func GetPaqueteByID(db *gorm.DB, id int) (*models.Paquete, error) {
	var paquete models.Paquete

	// Busca el paquete con el ID especificado
	result := db.Preload("AeropuertoOrigen.Ciudad.Pais").
		Preload("AeropuertoDestino.Ciudad.Pais").
		Preload("Habitaciones.Hotel.Ciudad.Pais").
		Preload("Habitaciones.TipoHabitacion").
		Preload("Aerolinea").
		First(&paquete, "id = ?", id)

	// Verifica si hay errores
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("paquete no encontrado")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &paquete, nil
}

// DeletePaquete elimina un paquete por su ID desde la base de datos.
func DeletePaquete(db *gorm.DB, id int) error {
	// Primero, elimina todas las relaciones en paquetes_habitaciones
	if err := db.Where("id_paquete = ?", id).Delete(models.PaqueteHabitacion{}).Error; err != nil {
		log.Printf("Error al eliminar relaciones en paquetes_habitaciones: %v\n", err)
		return err
	}

	// Luego, elimina el paquete
	if err := db.Delete(&models.Paquete{}, id).Error; err != nil {
		log.Printf("Error al eliminar el paquete: %v\n", err)
		return err
	}

	log.Printf("Paquete eliminado correctamente\n")
	return nil
}

func DeletePaquetesHabitacionesByPaqueteID(db *gorm.DB, paqueteID int) error {
	// Desactivar el modo de eliminación suave (soft delete) para eliminar definitivamente las filas
	return db.Unscoped().Where("id_paquete = ?", paqueteID).Delete(&models.PaqueteHabitacion{}).Error
}

// GetPaqueteByNombre busca un paquete por su nombre en la base de datos y lo retorna.
func GetPaqueteByNombre(db *gorm.DB, nombre string) (*models.Paquete, error) {
	var paquete models.Paquete
	result := db.First(&paquete, "nombre = ?", nombre)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // No se encontró ningún paquete con el nombre especificado
	}

	if result.Error != nil {
		log.Printf("Error al buscar el paquete por nombre: %v\n", result.Error)
		return nil, result.Error
	}

	return &paquete, nil
}

// CreatePaquete crea un nuevo paquete en la base de datos.
func CreatePaquete(db *gorm.DB, nuevoPaquete models.Paquete) error {
	// Verificar si el paquete ya existe
	existingPaquete, err := GetPaqueteByNombre(db, nuevoPaquete.Nombre)
	if err == nil && existingPaquete != nil {
		log.Printf("El paquete ya existe: %s\n", nuevoPaquete.Nombre)
		return errors.New("el paquete ya existe")
	} else if err != nil {
		log.Printf("Error al verificar la existencia del paquete: %v\n", err)
		return err
	}

	result := db.Create(&nuevoPaquete)
	if result.Error != nil {
		log.Printf("Error al crear el paquete: %v\n", result.Error)
		return result.Error
	}

	// A este punto, nuevoPaquete.ID debería estar establecido con el ID generado
	for _, habitacionID := range nuevoPaquete.HabitacionIDs {
		if err := db.Create(&models.PaqueteHabitacion{
			IDPaquete:    nuevoPaquete.ID,
			IDHabitacion: habitacionID,
		}).Error; err != nil {
			log.Printf("Error al crear la relación en paquetes_habitaciones: %v\n", err)
			return err
		}
	}

	log.Printf("Paquete creado correctamente: %s con ID: %d\n", nuevoPaquete.Nombre, nuevoPaquete.ID)
	return nil
}

// UpdatePaquete actualiza un paquete en la base de datos.
func UpdatePaquete(db *gorm.DB, paquete models.Paquete) error {
	// Obtener el paquete actual
	paqueteActual, err := GetPaqueteByID(db, paquete.ID)
	if err != nil {
		log.Printf("Error al obtener el paquete actual: %v\n", err)
		return err
	}

	// Verificar si el nombre del paquete ha cambiado
	if paqueteActual.Nombre != paquete.Nombre {
		// Verificar si el nuevo nombre ya existe
		existingPaquete, err := GetPaqueteByNombre(db, paquete.Nombre)
		if err == nil && existingPaquete != nil {
			log.Printf("El paquete ya existe: %s\n", paquete.Nombre)
			return errors.New("el paquete ya existe")
		} else if err != nil {
			log.Printf("Error al verificar la existencia del paquete: %v\n", err)
			return err
		}
	}

	// Actualizar el paquete
	result := db.Model(&paqueteActual).Updates(paquete)
	if result.Error != nil {
		log.Printf("Error al actualizar el paquete: %v\n", result.Error)
		return result.Error
	}

	// Actualizar las habitaciones
	if !reflect.DeepEqual(paqueteActual.HabitacionIDs, paquete.HabitacionIDs) {
		// Eliminar las relaciones anteriores
		if err := db.Where("id_paquete = ?", paquete.ID).Delete(models.PaqueteHabitacion{}).Error; err != nil {
			log.Printf("Error al eliminar las relaciones anteriores en paquetes_habitaciones: %v\n", err)
			return err
		}

		// Crear las nuevas relaciones
		for _, habitacionID := range paquete.HabitacionIDs {
			if err := db.Create(&models.PaqueteHabitacion{
				IDPaquete:    paquete.ID,
				IDHabitacion: habitacionID,
			}).Error; err != nil {
				log.Printf("Error al crear la nueva relación en paquetes_habitaciones: %v\n", err)
				return err
			}
		}
	}

	log.Printf("Paquete actualizado correctamente: %s\n", paquete.Nombre)
	return nil
}
