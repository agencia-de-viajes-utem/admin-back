// config/firebase.go
package config

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func InitFirebase() (*auth.Client, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("credentials.json")

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
