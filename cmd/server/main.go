package main

import (
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/kitbuilder587/cryptotrack/docs"
	"github.com/kitbuilder587/cryptotrack/internal/coingecko"
	"github.com/kitbuilder587/cryptotrack/internal/config"
	"github.com/kitbuilder587/cryptotrack/internal/db"
	"github.com/kitbuilder587/cryptotrack/internal/handlers"
	"github.com/kitbuilder587/cryptotrack/internal/migrations"
	"github.com/kitbuilder587/cryptotrack/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Инициализация конфига
	cfg := config.MustReturnConfig()

	// Соединение с БД
	dsn := "postgres://" + cfg.DB.User + ":" + cfg.DB.Password + "@" + cfg.DB.Host + ":" +
		strconv.Itoa(cfg.DB.Port) + "/cryptotrack?sslmode=disable"
	database, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to connect DB: %+v", err)
	}
	defer database.Close()

	// Миграции
	if err := migrations.Run(database); err != nil {
		log.Fatalf("failed to apply migrations: %+v", err)
	}

	// Инициализация сервисов и хендлеров
	repo := db.NewRepo(database)
	gecko := coingecko.NewClient()
	svc := service.NewTrackService(repo, gecko)
	api := handlers.NewAPI(svc)

	// Роутинг
	r := chi.NewRouter()
	r.Get("/health", api.Health)
	r.Post("/track", api.Track)
	r.Get("/latest", api.Latest)
	r.Get("/history", api.History)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Запуск
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("CryptoTrack running at :%s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("server error: %+v", err)
	}
}
