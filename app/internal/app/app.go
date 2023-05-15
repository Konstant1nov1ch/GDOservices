package app

import (
	"GDOservice/internal/config"
	"GDOservice/pkg/logging"
	"GDOservice/pkg/metric"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pressly/goose/v3"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

type App struct {
	cfg        *config.Config
	logger     *logging.Logger
	router     *httprouter.Router
	httpServer *http.Server
	db         *sql.DB
}

func NewApp(config *config.Config, logger *logging.Logger) (App, error) {
	logger.Println("router initializing")
	router := httprouter.New()

	logger.Println("swagger docs initializing")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	logger.Println("heartbeat metric initializing")
	metricHandler := metric.Handler{}
	metricHandler.Register(router)

	db, err := config.GetDB() // Получаем соединение с базой данных из config
	if err != nil {
		return App{}, err
	}

	return App{
		cfg:    config,
		logger: logger,
		router: router,
		db:     db, // Сохраняем соединение с базой данных в поле приложения
	}, nil
}

func (a *App) Run() {
	a.logger.Info("ping database")
	err := a.db.Ping()
	if err != nil {
		a.logger.Fatal(err)
	}
	// Запуск миграций
	err = runMigrations(a.db)
	if err != nil {
		a.logger.Fatal(err)
	}

	a.startHTTP()
}

// runMigrations ToDo panic: goose: duplicate version 1 detected: Users/konstantin/dev/pets/GDOservice/migrations/00001_init.up.sql
// /Users/konstantin/dev/pets/GDOservice/migrations/00001_init.down.sql
func runMigrations(db *sql.DB) error {
	// Укажите вашу конфигурацию миграций (путь к миграциям и источник данных)
	goose.SetDialect("postgres") // Замените на вашу используемую базу данных
	goose.SetTableName("goose_migrations")
	// Получение пути к текущему файлу
	_, filename, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(filename), "../../../migrations") // Замените на путь к вашим миграциям

	// Применение миграций
	err := goose.Up(db, migrationsDir)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) startHTTP() {
	a.logger.Info("start HTTP")

	var listener net.Listener

	if a.cfg.Listen.Type == config.LISTEN_TYPE_SOCK {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			a.logger.Fatal(err)
		}
		socketPath := path.Join(appDir, a.cfg.Listen.SocketFile)
		a.logger.Infof("socket path: %s", socketPath)

		a.logger.Info("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			a.logger.Fatal(err)
		}
	} else {
		a.logger.Infof("bind application to host: %s and port: %s", a.cfg.Listen.BindIP, a.cfg.Listen.Port)
		var err error
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
		if err != nil {
			a.logger.Fatal(err)
		}
	}

	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodOptions, http.MethodDelete},
		AllowedOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type", "Origin", "Accept", "Content-Length", "Accept-Encoding", "X-CSRF-Token"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Location", "Authorization", "Content-Disposition"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	a.logger.Println("application completely initialized and started")

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warn("server shutdown")
		default:
			a.logger.Fatal(err)
		}
	}
	err := a.httpServer.Shutdown(context.Background())
	if err != nil {
		a.logger.Fatal(err)
	}
}
