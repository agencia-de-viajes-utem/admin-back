package routes

import (
	"backend-admin/app/middleware"
	"backend-admin/app/models"
	"backend-admin/app/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ConfigureAeropuertoRouter(router *gin.Engine, authService *services.AuthService) {
	aeropuertoGroup := router.Group("/aeropuerto")

	// Ruta para GET (obtener todos los aeropuertos)
	aeropuertoGroup.GET("/get", func(c *gin.Context) {
		// Obtener todos los aeropuertos desde la base de datos
		aeropuertos, err := services.GetAeropuertos(authService.DB)
		if err != nil {
			log.Printf("Error al obtener los aeropuertos: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los aeropuertos"})
			return
		}

		// Responder con la lista de aeropuertos
		c.JSON(http.StatusOK, aeropuertos)
	})

	// Ruta protegida que requiere autenticación de administrador para PATCH (actualizar aeropuerto)
	aeropuertoGroup.PATCH("/update", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Deserializar el cuerpo de la solicitud en un objeto Aeropuerto
		var updatedAeropuerto models.Aeropuerto
		if err := c.BindJSON(&updatedAeropuerto); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		// Actualizar el registro en la base de datos
		if err := services.UpdateAeropuerto(authService.DB, &updatedAeropuerto); err != nil {
			log.Printf("Error al actualizar el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro actualizado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro actualizado correctamente"})
	})

	// Ruta protegida que requiere autenticación de administrador para DELETE (eliminar aeropuerto)
	aeropuertoGroup.DELETE("/delete", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Obtener el ID del aeropuerto a eliminar desde la solicitud
		aeropuertoIDStr := c.Query("id")

		// Convertir el ID del aeropuerto a entero
		aeropuertoID, err := strconv.Atoi(aeropuertoIDStr)
		if err != nil {
			log.Printf("Error al convertir el ID del aeropuerto: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al convertir el ID del aeropuerto"})
			return
		}

		// Eliminar el registro de la base de datos
		if err := services.DeleteAeropuerto(authService.DB, aeropuertoID); err != nil {
			log.Printf("Error al eliminar el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro eliminado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro eliminado correctamente"})
	})

	// Ruta protegida que requiere autenticación de administrador para POST (crear aeropuerto)
	aeropuertoGroup.POST("/create", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Deserializar el cuerpo de la solicitud en un objeto Aeropuerto
		var newAeropuerto models.Aeropuerto
		if err := c.BindJSON(&newAeropuerto); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		// Crear el registro en la base de datos
		if err := services.CreateAeropuerto(authService.DB, newAeropuerto); err != nil {
			log.Printf("Error al crear el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro creado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro creado correctamente"})
	})
}
