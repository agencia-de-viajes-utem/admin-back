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

func ConfigureFechaPaqueteRouter(router *gin.Engine, authService *services.AuthService) {
	fechaPaqueteGroup := router.Group("/fecha_paquete")

	// Ruta para GET (obtener todas las fechas_paquete)
	fechaPaqueteGroup.GET("/get", func(c *gin.Context) {
		// Obtener todas las fechas_paquete desde la base de datos
		fechasPaquete, err := services.GetFechasPaquete(authService.DB)
		if err != nil {
			log.Printf("Error al obtener las fechas_paquete: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las fechas_paquete"})
			return
		}

		// Responder con la lista de fechas_paquete
		c.JSON(http.StatusOK, fechasPaquete)
	})

	// Ruta protegida que requiere autenticación de administrador para PATCH (actualizar fecha_paquete)
	fechaPaqueteGroup.PATCH("/update", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Deserializar el cuerpo de la solicitud en un objeto FechaPaquete
		var updatedFechaPaquete models.FechaPaquete
		if err := c.BindJSON(&updatedFechaPaquete); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		// Actualizar el registro en la base de datos
		if err := services.UpdateFechaPaquete(authService.DB, &updatedFechaPaquete); err != nil {
			log.Printf("Error al actualizar el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro actualizado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro actualizado correctamente"})
	})

	// Ruta protegida que requiere autenticación de administrador para DELETE (eliminar fecha_paquete)
	fechaPaqueteGroup.DELETE("/delete", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Obtener el ID de la fecha_paquete a eliminar desde la solicitud
		fechaPaqueteIDStr := c.Query("id")

		// Convertir el ID de la fecha_paquete a entero
		fechaPaqueteID, err := strconv.Atoi(fechaPaqueteIDStr)
		if err != nil {
			log.Printf("Error al convertir el ID de la fecha_paquete: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al convertir el ID de la fecha_paquete"})
			return
		}

		// Eliminar el registro de la base de datos
		if err := services.DeleteFechaPaquete(authService.DB, fechaPaqueteID); err != nil {
			log.Printf("Error al eliminar el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro eliminado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro eliminado correctamente"})
	})

	// Ruta protegida que requiere autenticación de administrador para POST (crear fecha_paquete)
	fechaPaqueteGroup.POST("/create", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Deserializar el cuerpo de la solicitud en un objeto FechaPaquete
		var nuevaFechaPaquete models.FechaPaquete
		if err := c.BindJSON(&nuevaFechaPaquete); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		// Crear el registro en la base de datos
		if err := services.CreateFechaPaquete(authService.DB, nuevaFechaPaquete); err != nil {
			log.Printf("Error al crear el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro creado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro creado correctamente"})
	})
}
