package main

import (
	"backend-admin/api/config"
	"backend-admin/api/routes"
	"backend-admin/api/services"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Carga las variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar las variables de entorno: ", err)
		return
	}

	// Inicializa Firebase y obtén el cliente de autenticación
	authClient, err := config.InitFirebase()
	if err != nil {
		fmt.Printf("Error initializing Firebase client: %v\n", err)
		return
	}

	// Inicializa la conexión a la base de datos PostgreSQL
	db := config.InitDatabase()

	// Crea una instancia de Gin
	r := gin.Default()

	// Configura CORS para permitir solicitudes desde el puerto 5173 y permitir el encabezado "Authorization"
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175", "http://localhost:5176"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	r.Use(cors.New(config))

	// Configura los servicios y controladores
	authService := services.NewAuthService(authClient, db)

	// Configura las rutas con Gin
	routes.ConfigurePaisRoutes(r, authService)
	routes.ConfigureCiudadRouter(r, authService)
	routes.ConfigureAeropuertoRouter(r, authService)
	routes.ConfigureAerolineaRouter(r, authService)
	routes.ConfigureHotelRouter(r, authService)
	routes.ConfigureTipoHabitacionRoute(r, authService)
	routes.ConfigureHabitacionRouter(r, authService)
	routes.ConfigurePaqueteRoutes(r, authService)
	routes.ConfigureFechaPaqueteRouter(r, authService)

	// Configura y ejecuta tu servidor Gin aquí
	port := ":8081"
	fmt.Printf("Servidor en ejecución en http://localhost%s\n", port)
	log.Fatal(r.Run(port))
}
