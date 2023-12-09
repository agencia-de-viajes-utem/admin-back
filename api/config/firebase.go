// config/firebase.go
package config

import (
	"context"
	"fmt"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func InitFirebase() (*auth.Client, error) {
	ctx := context.Background()

	// Buscar el archivo de credenciales en el directorio actual
	matchingPattern := "./gha-creds-*.json"
	matches, err := filepath.Glob(matchingPattern)
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("No se encontraron archivos de credenciales para Firebase.")
	}

	// Utilizar el primer archivo coincidente (puedes ajustar esto según tus necesidades)
	pathToCredentials := matches[0]

	opt := option.WithCredentialsFile(pathToCredentials)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return authClient, nil
}
