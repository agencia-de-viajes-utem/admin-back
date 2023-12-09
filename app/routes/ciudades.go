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

func ConfigureCiudadRouter(router *gin.Engine, authService *services.AuthService) {
	ciudadGroup := router.Group("/ciudad")

	// Ruta para GET (obtener todas las ciudades)
	ciudadGroup.GET("/get", func(c *gin.Context) {
		// Obtener todas las ciudades desde la base de datos
		ciudades, err := services.GetCiudades(authService.DB)
		if err != nil {
			log.Printf("Error al obtener las ciudades: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las ciudades"})
			return
		}

		// Responder con la lista de ciudades
		c.JSON(http.StatusOK, ciudades)
	})

	// Ruta protegida que requiere autenticación de administrador para PATCH (actualizar ciudad)
	ciudadGroup.PATCH("/update", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Deserializar el cuerpo de la solicitud en un objeto Ciudad
		var updatedCiudad models.Ciudad
		if err := c.BindJSON(&updatedCiudad); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		// Actualizar el registro en la base de datos
		if err := services.UpdateCiudad(authService.DB, &updatedCiudad); err != nil {
			log.Printf("Error al actualizar el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro actualizado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro actualizado correctamente"})
	})

	// Ruta protegida que requiere autenticación de administrador para DELETE (eliminar ciudad)
	ciudadGroup.DELETE("/delete", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Obtener el ID de la ciudad a eliminar desde la solicitud
		ciudadIDStr := c.Query("id")

		// Convertir el ID de la ciudad a entero
		ciudadID, err := strconv.Atoi(ciudadIDStr)
		if err != nil {
			log.Printf("Error al convertir el ID de la ciudad: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al convertir el ID de la ciudad"})
			return
		}

		// Eliminar el registro de la base de datos
		if err := services.DeleteCiudad(authService.DB, ciudadID); err != nil {
			log.Printf("Error al eliminar el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro eliminado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro eliminado correctamente"})
	})

	// Ruta protegida que requiere autenticación de administrador para POST (crear ciudad)
	ciudadGroup.POST("/create", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Deserializar el cuerpo de la solicitud en un objeto Ciudad
		var nuevaCiudad models.Ciudad
		if err := c.BindJSON(&nuevaCiudad); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		// Crear el registro en la base de datos
		if err := services.CreateCiudad(authService.DB, nuevaCiudad); err != nil {
			log.Printf("Error al crear el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro creado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro creado correctamente"})
	})
}
