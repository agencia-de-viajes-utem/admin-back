package services

import (
	"backend-admin/api/models"
	"errors"
	"log"
	"reflect"

	"gorm.io/gorm"
)

// GetFechasPaquete obtiene todas las fechas_paquete desde la base de datos.
func GetFechasPaquete(db *gorm.DB) ([]models.FechaPaquete, error) {
	var fechasPaquete []models.FechaPaquete

	result := db.Preload("Paquete.AeropuertoOrigen").Preload("Paquete.AeropuertoDestino").Find(&fechasPaquete) // Carga también los datos de Paquete
	if result.Error != nil {
		return nil, result.Error
	}

	return fechasPaquete, nil
}

func GetFechaPaqueteByID(db *gorm.DB, id int) (*models.FechaPaquete, error) {
	var fechaPaquete models.FechaPaquete

	result := db.Preload("Paquete").First(&fechaPaquete, "id = ?", id) // Carga los datos de Paquete

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("fecha_paquete no encontrada")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &fechaPaquete, nil
}

// DeleteFechaPaquete elimina una fecha_paquete por su ID desde la base de datos.
func DeleteFechaPaquete(db *gorm.DB, id int) error {
	var fechaPaquete models.FechaPaquete

	result := db.Delete(&fechaPaquete, id)
	if result.Error != nil {
		log.Printf("Error al eliminar la fecha_paquete: %v\n", result.Error)
		return result.Error
	}

	log.Printf("FechaPaquete eliminada correctamente\n")
	return nil
}

// CreateFechaPaquete crea una nueva fecha_paquete en la base de datos.
func CreateFechaPaquete(db *gorm.DB, nuevaFechaPaquete models.FechaPaquete) error {
	result := db.Create(&nuevaFechaPaquete)
	if result.Error != nil {
		log.Printf("Error al crear la fecha_paquete: %v\n", result.Error)
		return result.Error
	}

	log.Printf("FechaPaquete creada correctamente\n")
	return nil
}

func UpdateFechaPaquete(db *gorm.DB, updatedFechaPaquete *models.FechaPaquete) error {
	existingFechaPaquete, err := GetFechaPaqueteByID(db, updatedFechaPaquete.ID)
	if err != nil {
		log.Printf("Error al obtener la fecha_paquete: %v\n", err)
		return err
	}

	val := reflect.ValueOf(*updatedFechaPaquete)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Interface() != reflect.Zero(field.Type()).Interface() {
			name := val.Type().Field(i).Name
			existingVal := reflect.ValueOf(existingFechaPaquete).Elem().FieldByName(name)
			if existingVal.IsValid() && existingVal.CanSet() {
				existingVal.Set(field)
			}
		}
	}

	result := db.Save(existingFechaPaquete)
	if result.Error != nil {
		log.Printf("Error al actualizar la fecha_paquete: %v\n", result.Error)
		return result.Error
	}

	log.Printf("FechaPaquete actualizada correctamente\n")
	return nil
}
