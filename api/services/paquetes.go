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
		Preload("Habitaciones.TipoHabitacion"). // Agregando Preload para TipoHabitacion
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
	result := db.Preload("AeropuertoOrigen").
		Preload("AeropuertoDestino").
		Preload("Habitaciones.Hotel.Ciudad.Pais").
		Preload("Habitaciones.TipoHabitacion"). // Agregando Preload para TipoHabitacion
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
	existingPaquete, err := GetPaqueteByID(db, id)
	if err != nil {
		log.Printf("Error al obtener el paquete por ID: %v\n", err)
		return err
	}

	result := db.Delete(existingPaquete)
	if result.Error != nil {
		log.Printf("Error al eliminar el paquete: %v\n", result.Error)
		return result.Error
	}

	log.Printf("Paquete eliminado correctamente\n")
	return nil
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

	log.Printf("Paquete creado correctamente: %s\n", nuevoPaquete.Nombre)
	return nil
}

// UpdatePaquete actualiza un paquete en la base de datos.
func UpdatePaquete(db *gorm.DB, updatedPaquete *models.Paquete) error {
	existingPaquete, err := GetPaqueteByID(db, updatedPaquete.ID)
	if err != nil {
		log.Printf("Error al obtener el paquete por ID: %v\n", err)
		return err
	}

	// Usa la reflexión para actualizar los campos del paquete existente con los campos del paquete actualizado
	val := reflect.ValueOf(*updatedPaquete)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Interface() != reflect.Zero(field.Type()).Interface() {
			name := val.Type().Field(i).Name
			existingVal := reflect.ValueOf(existingPaquete).Elem().FieldByName(name)
			if existingVal.IsValid() && existingVal.CanSet() {
				existingVal.Set(field)
			}
		}
	}

	result := db.Save(existingPaquete)
	if result.Error != nil {
		log.Printf("Error al actualizar el paquete: %v\n", result.Error)
		return result.Error
	}

	log.Printf("Paquete actualizado correctamente\n")
	return nil
}
