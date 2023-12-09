package services

import (
	"backend-admin/api/models"
	"errors"
	"log"
	"reflect"

	"gorm.io/gorm"
)

func GetAeropuertos(db *gorm.DB) ([]models.Aeropuerto, error) {
	var aeropuertos []models.Aeropuerto

	result := db.Preload("Ciudad.Pais").Find(&aeropuertos) // Carga también los datos de Ciudad
	if result.Error != nil {
		return nil, result.Error
	}

	return aeropuertos, nil
}

func GetAeropuertoByID(db *gorm.DB, id int) (*models.Aeropuerto, error) {
	var aeropuerto models.Aeropuerto

	result := db.Preload("Ciudad.Pais").First(&aeropuerto, "id = ?", id) // Carga los datos de Ciudad

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("aeropuerto no encontrado")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &aeropuerto, nil
}

func DeleteAeropuerto(db *gorm.DB, id int) error {
	var aeropuerto models.Aeropuerto

	result := db.Delete(&aeropuerto, id)
	if result.Error != nil {
		log.Printf("Error al eliminar el aeropuerto: %v\n", result.Error)
		return result.Error
	}

	log.Printf("Aeropuerto eliminado correctamente\n")
	return nil
}

func CreateAeropuerto(db *gorm.DB, nuevoAeropuerto models.Aeropuerto) error {
	result := db.Create(&nuevoAeropuerto)
	if result.Error != nil {
		log.Printf("Error al crear el aeropuerto: %v\n", result.Error)
		return result.Error
	}

	return nil
}

func UpdateAeropuerto(db *gorm.DB, updatedAeropuerto *models.Aeropuerto) error {
	existingAeropuerto, err := GetAeropuertoByID(db, updatedAeropuerto.ID)
	if err != nil {
		log.Printf("Error al obtener el aeropuerto: %v\n", err)
		return err
	}

	val := reflect.ValueOf(*updatedAeropuerto)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Interface() != reflect.Zero(field.Type()).Interface() {
			name := val.Type().Field(i).Name
			existingVal := reflect.ValueOf(existingAeropuerto).Elem().FieldByName(name)
			if existingVal.IsValid() && existingVal.CanSet() {
				existingVal.Set(field)
			}
		}
	}

	result := db.Save(&existingAeropuerto)
	if result.Error != nil {
		log.Printf("Error al actualizar el aeropuerto: %v\n", result.Error)
		return result.Error
	}

	log.Printf("Aeropuerto actualizado correctamente\n")
	return nil
}
