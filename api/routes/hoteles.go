package routes

import (
	"backend-admin/api/middleware"
	"backend-admin/api/models"
	"backend-admin/api/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ConfigureHotelRouter(router *gin.Engine, authService *services.AuthService) {
	hotelGroup := router.Group("/hotel")

	// Ruta para GET (obtener todos los hoteles)
	hotelGroup.GET("/get", func(c *gin.Context) {
		// Obtener todos los hoteles desde la base de datos
		hoteles, err := services.GetHoteles(authService.DB)
		if err != nil {
			log.Printf("Error al obtener los hoteles: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los hoteles"})
			return
		}

		// Responder con la lista de hoteles
		c.JSON(http.StatusOK, hoteles)
	})

	hotelGroup.GET("/get/:id", func(c *gin.Context) {
		// Obtener el ID del hotel desde la URL
		hotelIDStr := c.Param("id")
		hotelID, err := strconv.Atoi(hotelIDStr)
		if err != nil {
			log.Printf("Error al convertir el ID del hotel: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al convertir el ID del hotel"})
			return
		}

		// Consulta la base de datos para obtener las habitaciones asociadas al hotel
		habitaciones, err := services.GetHabitacionesByHotelID(authService.DB, hotelID)
		if err != nil {
			log.Printf("Error al obtener las habitaciones: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las habitaciones"})
			return
		}

		// Responder con la información de HotelHabitaciones
		c.JSON(http.StatusOK, habitaciones)
	})

	// Ruta protegida que requiere autenticación de administrador para PATCH (actualizar hotel)
	hotelGroup.PATCH("/update", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Deserializar el cuerpo de la solicitud en un objeto Hotel
		var updatedHotel models.Hotel
		if err := c.BindJSON(&updatedHotel); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		// Actualizar el registro en la base de datos
		if err := services.UpdateHotel(authService.DB, &updatedHotel); err != nil {
			log.Printf("Error al actualizar el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro actualizado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro actualizado correctamente"})
	})

	// Ruta protegida que requiere autenticación de administrador para DELETE (eliminar hotel)
	hotelGroup.DELETE("/delete", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Obtener el ID del hotel a eliminar desde la solicitud
		hotelIDStr := c.Query("id")

		// Convertir el ID del hotel a entero
		hotelID, err := strconv.Atoi(hotelIDStr)
		if err != nil {
			log.Printf("Error al convertir el ID del hotel: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al convertir el ID del hotel"})
			return
		}

		// Eliminar el registro de la base de datos
		if err := services.DeleteHotel(authService.DB, hotelID); err != nil {
			log.Printf("Error al eliminar el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el registro"})
			return
		}
	})

	// Ruta protegida que requiere autenticación de administrador para POST (crear hotel)
	hotelGroup.POST("/create", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Deserializar el cuerpo de la solicitud en un objeto Hotel
		var newHotel models.Hotel
		if err := c.BindJSON(&newHotel); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		// Crear el registro en la base de datos
		if err := services.CreateHotel(authService.DB, newHotel); err != nil {
			log.Printf("Error al crear el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro creado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro creado correctamente"})
	})
}
