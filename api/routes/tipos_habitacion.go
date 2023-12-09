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

func ConfigureTipoHabitacionRoute(router *gin.Engine, authService *services.AuthService) {
	tipoHabitacionGroup := router.Group("/tipo_habitacion")

	// Ruta para GET (obtener todos los tipos de habitación)
	tipoHabitacionGroup.GET("/get", func(c *gin.Context) {
		// Obtener todos los tipos de habitación desde la base de datos
		tiposHabitacion, err := services.GetTiposHabitacion(authService.DB)
		if err != nil {
			log.Printf("Error al obtener los tipos de habitación: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los tipos de habitación"})
			return
		}

		// Responder con la lista de tipos de habitación
		c.JSON(http.StatusOK, tiposHabitacion)
	})

	// Ruta protegida que requiere autenticación de administrador para POST (crear tipo de habitación)
	tipoHabitacionGroup.POST("/create", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Deserializar el cuerpo de la solicitud en un objeto TipoHabitacion
		var newTipoHabitacion models.TipoHabitacion
		if err := c.BindJSON(&newTipoHabitacion); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		// Crear el registro en la base de datos
		if err := services.CreateTipoHabitacion(authService.DB, newTipoHabitacion); err != nil {
			log.Printf("Error al crear el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro creado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro creado correctamente"})
	})

	// Ruta protegida que requiere autenticación de administrador para PATCH (actualizar tipo de habitación)
	tipoHabitacionGroup.PATCH("/update", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Deserializar el cuerpo de la solicitud en un objeto TipoHabitacion
		var updatedTipoHabitacion models.TipoHabitacion
		if err := c.BindJSON(&updatedTipoHabitacion); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		// Actualizar el registro en la base de datos
		if err := services.UpdateTipoHabitacion(authService.DB, &updatedTipoHabitacion); err != nil {
			log.Printf("Error al actualizar el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro actualizado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro actualizado correctamente"})
	})

	// Ruta protegida que requiere autenticación de administrador para DELETE (eliminar tipo de habitación)
	tipoHabitacionGroup.DELETE("/delete", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Obtener el ID del tipo de habitación a eliminar desde la solicitud
		tipoHabitacionIDStr := c.Query("id")

		// Convertir el ID del tipo de habitación a entero
		tipoHabitacionID, err := strconv.Atoi(tipoHabitacionIDStr)
		if err != nil {
			log.Printf("Error al convertir el ID del tipo de habitación: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al convertir el ID del tipo de habitación"})
			return
		}

		// Eliminar el registro de la base de datos
		if err := services.DeleteTipoHabitacion(authService.DB, tipoHabitacionID); err != nil {
			log.Printf("Error al eliminar el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro eliminado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro eliminado correctamente"})
	})
}
