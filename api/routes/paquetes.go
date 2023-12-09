// routes/paquete_routes.go
package routes

import (
	"backend-admin/api/middleware"
	"backend-admin/api/models"
	"backend-admin/api/services"
	"backend-admin/api/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ConfigurePaqueteRoutes configures routes related to "Paquete"
func ConfigurePaqueteRoutes(router *gin.Engine, authService *services.AuthService) {
	paqueteGroup := router.Group("/paquete")

	// Ruta para GET (obtener todos los paquetes)
	paqueteGroup.GET("/get", func(c *gin.Context) {
		// Obtener todos los paquetes desde la base de datos
		paquetes, err := services.GetPaquetes(authService.DB)
		if err != nil {
			log.Printf("Error al obtener los paquetes: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los paquetes"})
			return
		}

		// Transformar los paquetes a su versión resumida
		paquetesResumidos := make([]models.PaqueteResumido, 0, len(paquetes))
		for _, paquete := range paquetes {
			paqueteResumido := utils.TransformarAPaqueteResumido(paquete)
			paquetesResumidos = append(paquetesResumidos, paqueteResumido)
		}

		// Responder con la lista de paquetes resumidos
		c.JSON(http.StatusOK, paquetesResumidos)
	})

	// Ruta protegida que requiere autenticación de administrador para PATCH (actualizar paquete)
	paqueteGroup.PATCH("/update", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Deserializar el cuerpo de la solicitud en un objeto Paquete
		var updatedPaquete models.Paquete
		if err := c.BindJSON(&updatedPaquete); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		// Actualizar el registro en la base de datos
		if err := services.UpdatePaquete(authService.DB, &updatedPaquete); err != nil {
			log.Printf("Error al actualizar el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro actualizado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro actualizado correctamente"})
	})

	// Ruta protegida que requiere autenticación de administrador para DELETE (eliminar paquete)
	paqueteGroup.DELETE("/delete", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Obtener el ID del paquete a eliminar desde la solicitud
		paqueteIDStr := c.Query("id")

		// Convertir el ID de cadena a un entero
		paqueteID, err := strconv.Atoi(paqueteIDStr)
		if err != nil {
			log.Printf("Error al convertir el ID del paquete a un entero: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID del paquete inválido"})
			return
		}

		// Eliminar el paquete de la base de datos
		if err := services.DeletePaquete(authService.DB, paqueteID); err != nil {
			log.Printf("Error al eliminar el paquete: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el paquete"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Paquete eliminado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Paquete eliminado correctamente"})
	})

	// Ruta protegida que requiere autenticación de administrador para POST (crear nuevo paquete)
	paqueteGroup.POST("/create", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Deserializar el cuerpo de la solicitud en un objeto Paquete
		var newPaquete models.Paquete
		if err := c.BindJSON(&newPaquete); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		// Crear un nuevo paquete en la base de datos
		if err := services.CreatePaquete(authService.DB, newPaquete); err != nil {
			log.Printf("Error al crear el paquete: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el paquete"})
			return
		}

		// Responder con un código de estado 201 (Created)
		log.Printf("Paquete creado correctamente\n")
		c.JSON(http.StatusCreated, gin.H{"message": "Paquete creado correctamente"})
	})
}
