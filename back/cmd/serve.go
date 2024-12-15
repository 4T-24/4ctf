package cmd

import (
	v1 "4ctf/api/v1"
	"4ctf/config"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
	"github.com/savsgio/atreugo/v11"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	keyCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Run: func(cmd *cobra.Command, args []string) {
		config := config.LoadConfig()
		server := atreugo.New(atreugo.Config{
			Addr: fmt.Sprintf(":%d", config.Server.Port),
		})

		// Verify that key is 64 bytes long
		if len(config.Server.Key) != 64 {
			// Generate a random key
			rdm := make([]byte, 64)
			for i := range rdm {
				r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(keyCharset))))
				rdm[i] = keyCharset[r.Int64()]
			}
			logrus.Fatal("Server key must be 64 bytes long, you can use the following key: ", string(rdm))
		}

		dsn := fmt.Sprintf("mysql://%s:%s@%s:%d/%s", config.MySql.User, config.MySql.Password, config.MySql.Host, config.MySql.Port, config.MySql.Database)
		uri, _ := url.Parse(dsn)
		migrate := dbmate.New(uri)
		migrate.SchemaFile = config.MySql.SchemaFile
		migrate.MigrationsDir = []string{config.MySql.MigrationsFolder}
		migrate.Log = logrus.New().Writer()

		err := migrate.CreateAndMigrate()
		if err != nil {
			logrus.
				WithError(err).
				Fatal("cannot run migration")
		}

		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.MySql.User, config.MySql.Password, config.MySql.Host, config.MySql.Port, config.MySql.Database))
		if err != nil {
			logrus.
				WithError(err).
				Fatal("cannot connect to database")
		}

		boil.SetDB(db)

		v1.SetupRoutes(server, config)

		logrus.Infof("Starting server on port %d", config.Server.Port)
		if err := server.ListenAndServe(); err != nil {
			logrus.Fatalf("Error starting server: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
