package services

import (
	"backend-admin/app/models"
	"errors"
	"log"
	"reflect"

	"gorm.io/gorm"
)

// GetPaises obtiene todos los países desde la base de datos.
func GetPaises(db *gorm.DB) ([]models.Pais, error) {
	var paises []models.Pais

	result := db.Find(&paises)
	if result.Error != nil {
		return nil, result.Error
	}

	return paises, nil
}

func GetPaisByID(db *gorm.DB, id int) (*models.Pais, error) {
	var pais models.Pais

	// Busca el pais con el ID especificado
	result := db.First(&pais, "id = ?", id)

	// Verifica si hay errores
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("usuario no encontrado")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &pais, nil

}

// DeletePais elimina un país por su ID desde la base de datos. Retorna un error si el país no existe.
func DeletePais(db *gorm.DB, id int) error {
	existingPais, err := GetPaisByID(db, id)
	if err != nil {
		log.Printf("Error al obtener el país por ID: %v\n", err)
		return err
	}

	result := db.Delete(existingPais)
	if result.Error != nil {
		log.Printf("Error al eliminar el país: %v\n", result.Error)
		return result.Error
	}

	log.Printf("País eliminado correctamente\n")
	return nil
}

// GetPaisByNombre busca un país por su nombre en la base de datos y lo retorna. Retorna nil si no se encuentra ningún país con el nombre especificado.
func GetPaisByNombre(db *gorm.DB, nombre string) (*models.Pais, error) {
	var pais models.Pais
	result := db.First(&pais, "nombre = ?", nombre)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // No se encontró ningún país con el nombre especificado
	}

	if result.Error != nil {
		log.Printf("Error al buscar el país por nombre: %v\n", result.Error)
		return nil, result.Error
	}

	return &pais, nil
}

// CreatePais crea un nuevo país en la base de datos. Retorna un error si el país ya existe.
func CreatePais(db *gorm.DB, nuevoPais models.Pais) error {
	// Verificar si el país ya existe
	existingPais, err := GetPaisByNombre(db, nuevoPais.Nombre)
	if err == nil && existingPais != nil {
		log.Printf("El país ya existe: %s\n", nuevoPais.Nombre)
		return errors.New("el país ya existe")
	} else if err != nil {
		log.Printf("Error al verificar la existencia del país: %v\n", err)
		return err
	}

	result := db.Create(&nuevoPais)
	if result.Error != nil {
		log.Printf("Error al crear el país: %v\n", result.Error)
		return result.Error
	}

	log.Printf("País creado correctamente: %s\n", nuevoPais.Nombre)
	return nil
}

func UpdatePais(db *gorm.DB, updatedPais *models.Pais) error {

	existingPais, err := GetPaisByID(db, updatedPais.ID)
	if err != nil {
		log.Printf("Error al obtener el registro: %v\n", err)
		return err
	}

	// Usa la reflexión para actualizar los campos del pais existente con los campos del pais actualizado
	val := reflect.ValueOf(*updatedPais)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Interface() != reflect.Zero(field.Type()).Interface() {
			name := val.Type().Field(i).Name
			existingVal := reflect.ValueOf(existingPais).Elem().FieldByName(name)
			if existingVal.IsValid() && existingVal.CanSet() {
				existingVal.Set(field)
			}
		}
	}

	result := db.Save(existingPais)

	if result.Error != nil {
		log.Printf("Error al actualizar el registro: %v\n", result.Error)
		return result.Error
	}

	log.Printf("Registro actualizado correctamente\n")
	return nil

}
