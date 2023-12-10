package routes

import (
	"backend-admin/api/middleware"
	"backend-admin/api/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ConfigureAerolineaRouter(router *gin.Engine, authService *services.AuthService) {
	aerolineaGroup := router.Group("/aerolinea")

	// Ruta para GET (obtener todos los aerolineas)
	aerolineaGroup.GET("/get", func(c *gin.Context) {
		// Obtener todos los aerolineas desde la base de datos
		aerolineas, err := services.GetAerolineas(authService.DB)
		if err != nil {
			log.Printf("Error al obtener los aerolineas: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los aerolineas"})
			return
		}

		// Responder con la lista de aerolineas
		c.JSON(http.StatusOK, aerolineas)
	})

	// Ruta para PATCH (actualizar aerolínea)
	aerolineaGroup.PATCH("/update", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		var request struct {
			ID          int    `json:"id"` // Incluye explícitamente el campo ID
			Nombre      string `json:"nombre"`
			Aeropuertos []int  `json:"aeropuertos"`
		}

		if err := c.BindJSON(&request); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		updateRequest := services.UpdateAerolineaRequest{
			ID:            request.ID,
			Nombre:        request.Nombre,
			AeropuertoIDs: request.Aeropuertos,
		}

		if err := services.UpdateAerolinea(authService.DB, updateRequest); err != nil {
			log.Printf("Error al actualizar el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el registro"})
			return
		}

		log.Printf("Registro actualizado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro actualizado correctamente"})
	})

	// Ruta protegida que requiere autenticación de administrador para DELETE (eliminar aerolinea)
	aerolineaGroup.DELETE("/delete", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Obtener el ID del aerolinea a eliminar desde la solicitud
		aerolineaIDStr := c.Query("id")

		// Convertir el ID del aerolinea a entero
		aerolineaID, err := strconv.Atoi(aerolineaIDStr)
		if err != nil {
			log.Printf("Error al convertir el ID del aerolinea: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al convertir el ID del aerolinea"})
			return
		}

		// Eliminar el registro de la base de datos
		if err := services.DeleteAerolinea(authService.DB, aerolineaID); err != nil {
			log.Printf("Error al eliminar el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el registro"})
			return
		}

		// Responder con un código de estado 200 (OK)
		log.Printf("Registro eliminado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro eliminado correctamente"})
	})

	aerolineaGroup.POST("/create", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		// Temporary struct to match incoming JSON
		var requestData struct {
			Nombre        string `json:"nombre"`
			AeropuertoIDs []int  `json:"aeropuertos"`
		}

		// Unmarshal JSON into the temporary struct
		if err := c.BindJSON(&requestData); err != nil {
			log.Printf("Error al deserializar el cuerpo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al deserializar el cuerpo de la solicitud"})
			return
		}

		// Construct the data for service function
		createData := services.CreateAerolineaRequest{
			Nombre:        requestData.Nombre,
			AeropuertoIDs: requestData.AeropuertoIDs,
		}

		// Call the service function to create the Aerolinea
		if err := services.CreateAerolinea(authService.DB, createData); err != nil {
			log.Printf("Error al crear el registro: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el registro"})
			return
		}

		// Respond with a status code 200 (OK)
		log.Printf("Registro creado correctamente\n")
		c.JSON(http.StatusOK, gin.H{"message": "Registro creado correctamente"})
	})

}
