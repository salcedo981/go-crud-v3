package config

import (
	"context"
	"fmt"
	"go_template_v3/pkg/global/model"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/FDSAP-Git-Org/hephaestus/encryption"
	utils_v1 "github.com/FDSAP-Git-Org/hephaestus/utils/v1"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// ADD DATABASE CONNECTION VARIABLE HERE
	DBConnList []gorm.DB

	DBErr error

	RedisClient *redis.Client
	RedisError  error
)

func DecryptDBConfig() (*model.Database, error) {
	decrypted := model.Database{}

	// CREDENTIALS
	decrypted.Host, DBErr = encryption.Decrypt(utils_v1.GetEnv("POSTGRES_HOST"), utils_v1.GetEnv("SECRET_KEY"))
	if DBErr != nil {
		return nil, DBErr
	}
	decrypted.Username, DBErr = encryption.Decrypt(utils_v1.GetEnv("POSTGRES_USERNAME"), utils_v1.GetEnv("SECRET_KEY"))
	if DBErr != nil {
		return nil, DBErr
	}
	decrypted.Password, DBErr = encryption.Decrypt(utils_v1.GetEnv("POSTGRES_PASSWORD"), utils_v1.GetEnv("SECRET_KEY"))
	if DBErr != nil {
		return nil, DBErr
	}
	decrypted.Port, DBErr = strconv.Atoi(utils_v1.GetEnv("POSTGRES_PORT"))
	if DBErr != nil {
		return nil, DBErr
	}
	decrypted.SSLMode = utils_v1.GetEnv("POSTGRES_SSL_MODE")
	decrypted.Timezone = utils_v1.GetEnv("POSTGRES_TIMEZONE")

	// --------------------------
	// GET ALL DATABASES FROM ENV
	// --------------------------
	for _, dbList := range os.Environ() {
		if strings.HasPrefix(dbList, "DB_") {
			dbName := strings.SplitN(dbList, "=", 2)[0]
			dbN, encErr := encryption.Decrypt(utils_v1.GetEnv(dbName), utils_v1.GetEnv("SECRET_KEY"))
			if encErr != nil {
				return nil, encErr
			}
			decrypted.DBList = append(decrypted.DBList, dbN)
		}
	}
	return &decrypted, nil
}

func PostgreSQLConnect() bool {
	// Decrypt the database configuration
	decData, decErr := DecryptDBConfig()
	if decErr != nil {
		fmt.Printf("Database config decryption error: %s\n", decErr.Error())
		return false
	}

	// Connect to the database
	// Note: In able to use the database, you must use the index of the database in the DBConnList which will start with 0
	// ex. config.DBConnList[0].Table("table_name").Find(dbResult) --> this will get the first database
	for _, decDB := range decData.DBList {
		var dbConn *gorm.DB
		dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s timezone=%s",
			decData.Host, decData.Username, decData.Password,
			decDB, decData.Port,
			decData.SSLMode, decData.Timezone)
		dbConn, DBErr = gorm.Open(postgres.Open(dns))

		// Check the database connection
		sqlDB, err := dbConn.DB()
		if err != nil {
			log.Fatalf("FAILED TO GET THE DATABASE INSTANCE: %v", err)
		}

		err = sqlDB.Ping()
		if err != nil {
			log.Fatalf("FAILED TO PING THE DATABASE: %v", err)
		} else {
			fmt.Printf("%s CONNECTION STATUS: ✔\n", strings.ToUpper(decDB))
		}

		DBConnList = append(DBConnList, *dbConn)
		// DBConnList = append(DBConnList, *dbConn)
	}

	fmt.Println("TIMEZONE:", decData.Timezone)
	return true
}

func RedisConnect(address, password string) bool {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	ping, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Can't ping redis:", err)
		return false
	}

	fmt.Println("PING REDIS:", ping)
	return true
}
