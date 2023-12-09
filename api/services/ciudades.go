package services

import (
	"backend-admin/api/models"
	"errors"
	"log"
	"reflect"

	"gorm.io/gorm"
)

// GetCiudades obtiene todas las ciudades desde la base de datos.
func GetCiudades(db *gorm.DB) ([]models.Ciudad, error) {
	var ciudades []models.Ciudad

	result := db.Preload("Pais").Find(&ciudades) // Carga también los datos de Pais
	if result.Error != nil {
		return nil, result.Error
	}

	return ciudades, nil
}

func GetCiudadByID(db *gorm.DB, id int) (*models.Ciudad, error) {
	var ciudad models.Ciudad

	result := db.Preload("Pais").First(&ciudad, "id = ?", id) // Carga los datos de Pais

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("ciudad no encontrada")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &ciudad, nil
}

// DeleteCiudad elimina una ciudad por su ID desde la base de datos.
func DeleteCiudad(db *gorm.DB, id int) error {
	var ciudad models.Ciudad

	result := db.Delete(&ciudad, id)
	if result.Error != nil {
		log.Printf("Error al eliminar la ciudad: %v\n", result.Error)
		return result.Error
	}

	log.Printf("Ciudad eliminada correctamente\n")
	return nil
}

// CreateCiudad crea una nueva ciudad en la base de datos.
func CreateCiudad(db *gorm.DB, nuevaCiudad models.Ciudad) error {
	result := db.Create(&nuevaCiudad)
	if result.Error != nil {
		log.Printf("Error al crear la ciudad: %v\n", result.Error)
		return result.Error
	}

	log.Printf("Ciudad creada correctamente\n")
	return nil
}

func UpdateCiudad(db *gorm.DB, updatedCiudad *models.Ciudad) error {
	existingCiudad, err := GetCiudadByID(db, updatedCiudad.ID)
	if err != nil {
		log.Printf("Error al obtener la ciudad: %v\n", err)
		return err
	}

	val := reflect.ValueOf(*updatedCiudad)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Interface() != reflect.Zero(field.Type()).Interface() {
			name := val.Type().Field(i).Name
			existingVal := reflect.ValueOf(existingCiudad).Elem().FieldByName(name)
			if existingVal.IsValid() && existingVal.CanSet() {
				existingVal.Set(field)
			}
		}
	}

	result := db.Save(existingCiudad)
	if result.Error != nil {
		log.Printf("Error al actualizar la ciudad: %v\n", result.Error)
		return result.Error
	}

	log.Printf("Ciudad actualizada correctamente\n")
	return nil
}
