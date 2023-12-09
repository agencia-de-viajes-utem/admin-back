package services

import (
	"backend-admin/api/models"
	"errors"
	"log"
	"reflect"

	"gorm.io/gorm"
)

func GetHoteles(db *gorm.DB) ([]models.Hotel, error) {
	var hoteles []models.Hotel

	result := db.Preload("Ciudad.Pais").Preload("Ciudad").Find(&hoteles) // Cargar Ciudad y luego Pais
	if result.Error != nil {
		return nil, result.Error
	}

	return hoteles, nil
}

func GetHotelByID(db *gorm.DB, id int) (*models.Hotel, error) {
	var hotel models.Hotel

	result := db.Preload("Ciudad.Pais").First(&hotel, "id = ?", id) // Carga los datos de Ciudad

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("hotel no encontrado")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &hotel, nil
}
func GetHabitacionesByHotelID(db *gorm.DB, hotelID int) ([]models.HabitacionConPaquete, error) {
	var habitaciones []models.Habitacion
	var habitacionesConPaquete []models.HabitacionConPaquete

	// Consulta la base de datos para obtener las habitaciones asociadas al hotel
	result := db.Preload("TipoHabitacion").
		Preload("Hotel").
		Preload("Hotel.Ciudad").
		Preload("Hotel.Ciudad.Pais").
		Where("id_hotel = ?", hotelID).
		Find(&habitaciones)
	if result.Error != nil {
		return nil, result.Error
	}

	// Para cada habitación, verifica si está ocupada y crea un objeto HabitacionConPaquete
	for _, h := range habitaciones {
		habitacionConPaquete, err := models.EstaOcupada(h, db)
		if err != nil {
			return nil, err
		}
		habitacionesConPaquete = append(habitacionesConPaquete, habitacionConPaquete)
	}

	return habitacionesConPaquete, nil
}

func DeleteHotel(db *gorm.DB, id int) error {
	var hotel models.Hotel

	result := db.Delete(&hotel, id)
	if result.Error != nil {
		log.Printf("Error al eliminar el hotel: %v\n", result.Error)
		return result.Error
	}

	log.Printf("Hotel eliminado correctamente\n")
	return nil
}

func CreateHotel(db *gorm.DB, nuevoHotel models.Hotel) error {
	result := db.Create(&nuevoHotel)
	if result.Error != nil {
		log.Printf("Error al crear el hotel: %v\n", result.Error)
		return result.Error
	}

	return nil
}

func UpdateHotel(db *gorm.DB, updatedHotel *models.Hotel) error {
	var hotel models.Hotel

	// Obtener el registro de la base de datos
	result := db.First(&hotel, "id = ?", updatedHotel.ID)
	if result.Error != nil {
		log.Printf("Error al obtener el registro: %v\n", result.Error)
		return result.Error
	}

	// Actualizar el registro
	result = db.Model(&hotel).Updates(updatedHotel)
	if result.Error != nil {
		log.Printf("Error al actualizar el registro: %v\n", result.Error)
		return result.Error
	}

	// Actualizar el registro en el objeto original
	reflect.ValueOf(&hotel).Elem().Set(reflect.ValueOf(updatedHotel).Elem())

	return nil
}
