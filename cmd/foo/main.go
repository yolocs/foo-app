package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/yolocs/foo-app/pkg/app"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
)

type envConfig struct {
	DbUser      string `envconfig:"DB_USER" required:"true"`
	DbPwd       string `envconfig:"DB_PASS" required:"true"`
	DbName      string `envconfig:"DB_NAME" required:"true"`
	TCPHostName string `envconfig:"TCP_HOST_NAME" default:"127.0.0.1"`
	BucketName  string `envconfig:"BUCKET_NAME" required:"true"`
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("Failed to create logger", err)
	}
	defer logger.Sync()

	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		logger.Fatal("Failed to parse env vars", zap.Error(err))
	}

	db, err := initDbConnPool(env)
	if err != nil {
		logger.Error("Failed to create DB connection pool", zap.Error(err))
	}

	app := &app.FooApp{
		DB:     db,
		Bucket: env.BucketName,
		Logger: logger,
	}
	http.HandleFunc("/", app.Handler)
	logger.Info("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal("Failed to listen http", zap.Error(err))
	}
}

func initDbConnPool(env envConfig) (*sql.DB, error) {
	var dbURI string
	dbURI = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", env.DbUser, env.DbPwd, env.TCPHostName, 3306, env.DbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
}

func configureConnectionPool(dbPool *sql.DB) {
	// [START cloud_sql_mysql_databasesql_limit]

	// Set maximum number of connections in idle connection pool.
	dbPool.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	dbPool.SetMaxOpenConns(7)

	// [END cloud_sql_mysql_databasesql_limit]

	// [START cloud_sql_mysql_databasesql_lifetime]

	// Set Maximum time (in seconds) that a connection can remain open.
	dbPool.SetConnMaxLifetime(1800)

	// [END cloud_sql_mysql_databasesql_lifetime]
}
