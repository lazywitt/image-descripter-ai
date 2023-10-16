package main

import (
	"bufio"
	"context"
	"fmt"
	apiService "heyalley-server/api"
	"heyalley-server/db"
	"heyalley-server/db/models"
	imageCronService "heyalley-server/image/ai_processor"
	imageService "heyalley-server/image/service"
	"log"
	"net/http"
	"os/exec"
	"time"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

func initFetchHttpServer(ctx context.Context, apiService *apiService.HttpService) {
	mux := http.NewServeMux()
	mux.HandleFunc("/image/register", apiService.StoreImage)
	mux.HandleFunc("/image/search", apiService.SearchImage)
	http.ListenAndServe(":8080", mux)
}

func main() {

	ctx := context.Background()

	// init AI model
	initModel()

	// init DB
	postgresClient, err := initDb(ctx, db.GetConf())
	if err != nil {
		fmt.Printf("error Initialising DB and entities: %v\n", err)
	}
	dbClient := db.GetHandler(postgresClient)

	cron := &imageCronService.ImageProcessingCron{
		SleepTime: time.Second * 5,
		Pipeline: &imageCronService.ImagePipeline{
			QueryEngine: dbClient,
		},
	}

	// init cron
	go initImageProcessorCron(ctx, cron)

	imageRegisteryService := &imageService.ImageRegisteryService{
		DatabaseService: dbClient,
	}
	imageStoreService := &imageService.ImageStoreService{
		ImageRegisteryService: imageRegisteryService,
		DatabaseService:       dbClient,
	}
	httpService := &apiService.HttpService{
		ImageStorageService: imageStoreService,
	}
	// init server
	initFetchHttpServer(ctx, httpService)
}

// initDb Create entities in PGDB
func initDb(ctx context.Context, config *db.Config) (*gorm.DB, error) {
	dbClient := getDbClient(ctx, config)
	err := dbClient.AutoMigrate(&models.Image{})
	if err != nil {
		return nil, err
	}
	// create indexes
	createIndexRes := dbClient.Exec(`CREATE INDEX IF NOT EXISTS text_search_idx ON images USING GIN (to_tsvector('english', description));
	CREATE INDEX IF NOT EXISTS created_at_idx ON images (created_at);`)
	if createIndexRes.Error != nil {
		return nil, createIndexRes.Error
	}
	return dbClient, nil
}

// getDbClient returns postgres client
func getDbClient(ctx context.Context, pgdbConf *db.Config) *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		pgdbConf.Host, pgdbConf.User, pgdbConf.Password, pgdbConf.Dbname, pgdbConf.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("error connecting to postgres", err)
	}
	fmt.Println("successfully connected to postgres")
	return db
}

// Init Model
func initModel() {
	cmd := exec.Command("bash", "./server/initmodel.sh")
	stderr, _ := cmd.StderrPipe()
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(stdout))

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

// Init image processor cron
func initImageProcessorCron(ctx context.Context, cron imageCronService.Cron) {
	cron.Trigger(ctx)
}
