package services

import (
	"backend-admin/app/models"
	"errors"
	"log"
	"reflect"

	"gorm.io/gorm"
)

func GetTiposHabitacion(db *gorm.DB) ([]models.TipoHabitacion, error) {
	var tiposHabitacion []models.TipoHabitacion

	result := db.Find(&tiposHabitacion)
	if result.Error != nil {
		return nil, result.Error
	}

	return tiposHabitacion, nil
}

func GetTipoHabitacionByID(db *gorm.DB, id int) (*models.TipoHabitacion, error) {
	var tipoHabitacion models.TipoHabitacion

	result := db.First(&tipoHabitacion, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("tipo de habitación no encontrado")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &tipoHabitacion, nil
}

func DeleteTipoHabitacion(db *gorm.DB, id int) error {
	var tipoHabitacion models.TipoHabitacion

	result := db.Delete(&tipoHabitacion, id)
	if result.Error != nil {
		log.Printf("Error al eliminar el tipo de habitación: %v\n", result.Error)
		return result.Error
	}

	log.Printf("Tipo de habitación eliminado correctamente\n")
	return nil
}

func CreateTipoHabitacion(db *gorm.DB, newTipoHabitacion models.TipoHabitacion) error {
	result := db.Create(&newTipoHabitacion)
	if result.Error != nil {
		log.Printf("Error al crear el tipo de habitación: %v\n", result.Error)
		return result.Error
	}

	return nil
}

func UpdateTipoHabitacion(db *gorm.DB, updatedTipoHabitacion *models.TipoHabitacion) error {
	existingTipoHabitacion, err := GetTipoHabitacionByID(db, updatedTipoHabitacion.ID)
	if err != nil {
		log.Printf("Error al obtener el tipo de habitación: %v\n", err)
		return err
	}

	val := reflect.ValueOf(*updatedTipoHabitacion)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Interface() != reflect.Zero(field.Type()).Interface() {
			name := val.Type().Field(i).Name
			existingVal := reflect.ValueOf(existingTipoHabitacion).Elem().FieldByName(name)
			if existingVal.IsValid() && existingVal.CanSet() {
				existingVal.Set(field)
			}
		}
	}

	result := db.Save(existingTipoHabitacion)
	if result.Error != nil {
		log.Printf("Error al actualizar el tipo de habitación: %v\n", result.Error)
		return result.Error
	}

	log.Printf("Tipo de habitación actualizado correctamente\n")
	return nil
}
