// routes/pais_routes.go

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

// ConfigurePaisRoutes configura las rutas relacionadas con "País"
func ConfigurePaisRoutes(router *gin.Engine, authService *services.AuthService) {
	paisGroup := router.Group("/pais")

	// Ruta para GET (obtener todos los países)
	paisGroup.GET("/get", func(c *gin.Context) {
		// Obtener todos los países desde la base de datos
		paises, err := services.GetPaises(authService.DB)
		if err != nil {
			log.Printf("Error al obtener los países: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los países"})
			return
		}

		// Responder con la lista de países
		c.JSON(http.StatusOK, paises)
	})

	// Ruta protegida que requiere autenticación de administrador para PATCH (actualizar país)
	paisGroup.PATCH("/update", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Deserializar el cuerpo de la solicitud en un objeto Pais
		var updatedPais models.Pais
		if err := c.BindJSON(&updatedPais); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		// Actualizar el registro en la base de datos
		if err := services.UpdatePais(authService.DB, &updatedPais); err != nil {
			log.Printf("Error al actualizar el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro actualizado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro actualizado correctamente"})
	})

	// Ruta protegida que requiere autenticación de administrador para DELETE (eliminar país)
	paisGroup.DELETE("/delete", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Obtener el ID del país a eliminar desde la solicitud
		paisIDStr := c.Query("id")

		// Convertir el ID de cadena a un entero
		paisID, err := strconv.Atoi(paisIDStr)
		if err != nil {
			log.Printf("Error al convertir el ID del país a un entero: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID del país inválido"})
			return
		}

		// Eliminar el país de la base de datos
		if err := services.DeletePais(authService.DB, paisID); err != nil {
			log.Printf("Error al eliminar el país: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el país"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("País eliminado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "País eliminado correctamente"})
	})

	// Ruta protegida que requiere autenticación de administrador para POST (crear nuevo país)
	paisGroup.POST("/create", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Deserializar el cuerpo de la solicitud en un objeto Pais
		var newPais models.Pais
		if err := c.BindJSON(&newPais); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		// Crear un nuevo país en la base de datos
		if err := services.CreatePais(authService.DB, newPais); err != nil {
			log.Printf("Error al crear el país: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el país"})
			return
		}

		// Responder con un código de estado 201 (Created)
		log.Printf("País creado correctamente\n")
		c.JSON(http.StatusCreated, gin.H{"message": "País creado correctamente"})
	})
}
