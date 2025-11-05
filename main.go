package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go_template_v3/pkg/config"
	"go_template_v3/routers"
	"log"
	"strings"

	"github.com/FDSAP-Git-Org/hephaestus/apilogs"
	utils_v1 "github.com/FDSAP-Git-Org/hephaestus/utils/v1"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment name
	env := utils_v1.GetEnv("ENVIRONMENT")
	loadedEnv := strings.ToLower(env)

	fmt.Println("ENVIRONMENT:", strings.ToUpper(loadedEnv))
	// Load environment settings
	if envErr := godotenv.Load(fmt.Sprintf("./envs/.env-%s", loadedEnv)); envErr != nil {
		log.Fatal("Error loading env file:", envErr)
	}

	fmt.Println("PROJECT: ", utils_v1.GetEnv("PROJECT"))
	fmt.Println("DESCRIPTION: ", utils_v1.GetEnv("DESCRIPTION"))

	folders := []string{"system"}
	apilogs.CreateInitialFolder(folders)

	// Connect to DB
	config.PostgreSQLConnect()
}

func main() {
	app := fiber.New(fiber.Config{
		AppName:          utils_v1.GetEnv("PROJECT"),
		CaseSensitive:    true,
		DisableKeepalive: true,
		JSONEncoder:      json.Marshal,
		JSONDecoder:      json.Unmarshal,
	})

	// CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET,POST,PUT,DELETE"},
		AllowHeaders: []string{"Origin, Content-Type, Accept, Authorization"},
	}))

	app.Use(logger.New())
	app.Use(recover.New())

	// Initialize API Endpoints
	routers.APIRoute(app)

	// TLS Configuration
	if strings.ToUpper(utils_v1.GetEnv("SSL_MODE")) == "ENABLED" {
		fmt.Println("SSL_MODE: ENABLED")
		fmt.Println("CERTIFICATE:", utils_v1.GetEnv("SSL_CERTIFICATE"))
		fmt.Println("KEY:", utils_v1.GetEnv("SSL_KEY"))

		// LOAD CERTIFICATE
		cert, err := tls.LoadX509KeyPair(utils_v1.GetEnv("SSL_CERTIFICATE"), utils_v1.GetEnv("SSL_KEY"))
		if err != nil {
			log.Fatal(err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		// START THE SERVER WITH HTTPS
		tlsPort := fmt.Sprintf(":%s", utils_v1.GetEnv("PORT"))
		listener, err := tls.Listen("tcp", tlsPort, tlsConfig)
		if err != nil {
			log.Fatalf("Failed to create TLS listener: %v", err)
		}

		if err := app.Listener(listener); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	} else {
		fmt.Println("SSL_MODE: DISABLED")
		log.Fatal(app.Listen(fmt.Sprintf(":%s", utils_v1.GetEnv("PORT"))))
	}
}
