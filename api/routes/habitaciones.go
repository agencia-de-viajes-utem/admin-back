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

func ConfigureHabitacionRouter(router *gin.Engine, authService *services.AuthService) {
	habitacionGroup := router.Group("/habitacion")

	// Ruta para GET (obtener todas las habitaciones)
	habitacionGroup.GET("/get", func(c *gin.Context) {
		habitaciones, err := services.GetHabitaciones(authService.DB)
		if err != nil {
			log.Printf("Error al obtener las habitaciones: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las habitaciones"})
			return
		}

		c.JSON(http.StatusOK, habitaciones)
	})

	// Ruta para GET (obtener una habitación)
	habitacionGroup.GET("/get/:id", func(c *gin.Context) {
		habitacionIDStr := c.Param("id")
		habitacionID, err := strconv.Atoi(habitacionIDStr)
		if err != nil {
			log.Printf("Error al convertir el ID de la habitación: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al convertir el ID de la habitación"})
			return
		}

		// Consulta la base de datos para obtener la habitación
		habitacion, err := services.GetHabitacionByID(authService.DB, habitacionID)
		if err != nil {
			log.Printf("Error al obtener la habitación: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la habitación"})
			return
		}

		c.JSON(http.StatusOK, habitacion)
	})

	// Ruta protegida para PATCH (actualizar habitación)
	habitacionGroup.PATCH("/update", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		var updatedHabitacion models.Habitacion
		if err := c.BindJSON(&updatedHabitacion); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		if err := services.UpdateHabitacion(authService.DB, &updatedHabitacion); err != nil {
			log.Printf("Error al actualizar la habitación: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la habitación"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Habitación actualizada correctamente"})
	})

	// Ruta protegida para DELETE (eliminar habitación)
	habitacionGroup.DELETE("/delete", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		habitacionIDStr := c.Query("id")
		habitacionID, err := strconv.Atoi(habitacionIDStr)
		if err != nil {
			log.Printf("Error al convertir el ID de la habitación: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al convertir el ID de la habitación"})
			return
		}

		if err := services.DeleteHabitacion(authService.DB, habitacionID); err != nil {
			log.Printf("Error al eliminar la habitación: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la habitación"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Habitación eliminada correctamente"})
	})

	// Ruta protegida para POST (crear habitación)
	habitacionGroup.POST("/create", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		var nuevaHabitacion models.Habitacion
		if err := c.BindJSON(&nuevaHabitacion); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		if err := services.CreateHabitacion(authService.DB, nuevaHabitacion); err != nil {
			log.Printf("Error al crear la habitación: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la habitación"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Habitación creada correctamente"})
	})
}
