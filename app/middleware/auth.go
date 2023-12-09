package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Función para verificar si un usuario es administrador
func isUserAdmin(token string) (bool, error) {
	// Realiza una solicitud HTTP al backend "Usuarios" para obtener la información del usuario en formato JSON
	url := "http://localhost:8080/user/info"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	// Agrega el token JWT en el encabezado de autorización como "Bearer {token}"
	req.Header.Set("Authorization", "Bearer "+token)

	// Realiza la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Leer y analizar la respuesta JSON
	var userInfo map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&userInfo); err != nil {
		return false, err
	}

	// Verificar si el usuario es admin basado en la clave "Rol" del JSON
	if role, ok := userInfo["Rol"].(string); ok && role == "admin" {
		return true, nil // El usuario es admin
	}

	return false, nil // El usuario no es admin
}

// Middleware para verificar si el usuario es administrador
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extraer el token del encabezado 'Authorization'
		authHeader := c.GetHeader("Authorization")

		// Dividir el encabezado en el espacio ' '
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("Encabezado de autorización mal formado o ausente: %s\n", authHeader)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Encabezado de autorización mal formado o ausente"})
			c.Error(fmt.Errorf("encabezado de autorización mal formado o ausente"))
			return
		}

		// El token es la segunda parte del encabezado dividido
		token := parts[1]

		// Verificar si el usuario es admin
		isAdmin, err := isUserAdmin(token)
		if err != nil {
			log.Printf("Error al verificar el rol del usuario: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar el rol del usuario"})
			return
		}

		if isAdmin {
			// El usuario es admin, permite que la solicitud continúe
			c.Next()
		} else {
			// El usuario no es admin, responde con un error o toma las acciones correspondientes
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Acceso no autorizado"})
			c.Abort()
		}
	}
}
