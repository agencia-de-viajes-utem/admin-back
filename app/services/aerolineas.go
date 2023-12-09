package services

import (
	"backend-admin/app/models"
	"errors"
	"log"

	"gorm.io/gorm"
)

func GetAerolineas(db *gorm.DB) ([]models.Aerolinea, error) {
	var aerolineas []models.Aerolinea

	result := db.Preload("Aeropuertos.Ciudad.Pais").Find(&aerolineas)
	if result.Error != nil {
		return nil, result.Error
	}

	return aerolineas, nil
}

func GetAerolineaByID(db *gorm.DB, id int) (*models.Aerolinea, error) {
	var aerolinea models.Aerolinea

	result := db.Preload("Aeropuertos.Ciudad.Pais").First(&aerolinea, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("aerolinea no encontrada")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &aerolinea, nil
}

// una estructura temporal para actualizar la aerolinea
type UpdateAerolineaRequest struct {
	ID            int    `json:"id"`
	Nombre        string `json:"nombre"`
	AeropuertoIDs []int  `json:"aeropuertos"`
}

func UpdateAerolinea(db *gorm.DB, data UpdateAerolineaRequest) error {
	// Busca la aerolínea existente por ID
	var aerolinea models.Aerolinea
	if err := db.Preload("Aeropuertos").First(&aerolinea, data.ID).Error; err != nil {
		return err // Manejar si no se encuentra la aerolínea o hay otros errores de BD
	}

	// Actualiza el nombre de la aerolínea
	aerolinea.Nombre = data.Nombre

	// Obtén las relaciones de aeropuertos actuales
	var aeropuertosActuales []models.Aeropuerto
	for _, aeropuerto := range aerolinea.Aeropuertos {
		aeropuertosActuales = append(aeropuertosActuales, aeropuerto)
	}

	// Obtiene las relaciones nuevas, eliminando duplicados
	var aeropuertosNuevos []models.Aeropuerto
	aeropuertoIDs := make(map[int]bool)
	for _, aeropuertoID := range data.AeropuertoIDs {
		if !aeropuertoIDs[aeropuertoID] {
			aeropuertoIDs[aeropuertoID] = true
			aeropuertosNuevos = append(aeropuertosNuevos, models.Aeropuerto{ID: aeropuertoID})
		}
	}

	// Calcula las relaciones a eliminar
	var aeropuertosAEliminar []models.Aeropuerto
	for _, aeropuertoActual := range aeropuertosActuales {
		encontrado := false
		for _, aeropuertoNuevo := range aeropuertosNuevos {
			if aeropuertoActual.ID == aeropuertoNuevo.ID {
				encontrado = true
				break
			}
		}
		if !encontrado {
			aeropuertosAEliminar = append(aeropuertosAEliminar, aeropuertoActual)
		}
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Elimina las relaciones antiguas innecesarias
	if len(aeropuertosAEliminar) > 0 {
		if err := tx.Model(&aerolinea).Association("Aeropuertos").Delete(aeropuertosAEliminar); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Agrega las nuevas relaciones
	if len(aeropuertosNuevos) > 0 {
		if err := tx.Model(&aerolinea).Association("Aeropuertos").Append(aeropuertosNuevos); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Guarda la aerolínea actualizada
	if err := tx.Save(&aerolinea).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func DeleteAerolinea(db *gorm.DB, id int) error {
	var aerolinea models.Aerolinea

	// Obtén la aerolínea existente, incluyendo sus relaciones
	if err := db.Preload("Aeropuertos").First(&aerolinea, id).Error; err != nil {
		return err // Manejar si no se encuentra la aerolínea o hay otros errores de BD
	}

	// Elimina todas las relaciones con aeropuertos
	if err := db.Model(&aerolinea).Association("Aeropuertos").Clear(); err != nil {
		log.Printf("Error al eliminar relaciones con aeropuertos: %v\n", err)
		return err
	}

	// Luego, elimina la aerolínea después de haber eliminado sus relaciones
	if err := db.Delete(&aerolinea, id).Error; err != nil {
		log.Printf("Error al eliminar la aerolinea: %v\n", err)
		return err
	}

	log.Printf("Aerolinea eliminada correctamente\n")
	return nil
}

type CreateAerolineaRequest struct {
	Nombre        string `json:"nombre"`
	AeropuertoIDs []int  `json:"aeropuertos"`
}

func CreateAerolinea(db *gorm.DB, data CreateAerolineaRequest) error {
	var aeropuertos []models.Aeropuerto
	if len(data.AeropuertoIDs) > 0 {
		if err := db.Find(&aeropuertos, "id IN ?", data.AeropuertoIDs).Error; err != nil {
			return err // Handle DB error
		}
	}

	aerolinea := models.Aerolinea{
		Nombre:      data.Nombre,
		Aeropuertos: aeropuertos,
	}

	return db.Create(&aerolinea).Error
}
