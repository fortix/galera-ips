package monitor

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	db      *sql.DB
	Writers *[]string
	Readers *[]string
)

func MonitorProxySQL() {

	// Initialize the list of readers and writers
	Writers = &[]string{}
	Readers = &[]string{}

	for {
		log.Info().Msg("Connecting to ProxySQL")

		var err error
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/",
			viper.GetString("proxysql.username"),
			viper.GetString("proxysql.password"),
			viper.GetString("proxysql.host"),
			viper.GetString("proxysql.port"),
		))
		if err != nil {
			log.Error().Msgf("Error connecting to ProxySQL, retrying in 5 seconds: %s", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	log.Info().Msg("Entering monitoring loop")
	for {
		runQuery()
		time.Sleep(1 * time.Second)
	}
}

func runQuery() {
	rows, err := db.Query(`SELECT hostgroup_id,hostname FROM runtime_mysql_servers WHERE hostgroup_id IN(2,3) AND status="ONLINE" ORDER BY hostname ASC;`)
	if err != nil {

		// Clear the list of readers / writers
		Writers = &[]string{}
		Readers = &[]string{}

		log.Error().Msgf("Error running query: %s", err)

		time.Sleep(5 * time.Second)
		return
	}
	defer rows.Close()

	w := []string{}
	r := []string{}

	for rows.Next() {
		var hostgroupID int
		var hostname string
		if err := rows.Scan(&hostgroupID, &hostname); err != nil {
			log.Error().Msgf("Error scanning row: %s", err)
			return
		}

		log.Debug().Msgf("hostgroupID: %d, hostname: %s", hostgroupID, hostname)

		if hostgroupID == 2 {
			w = append(w, hostname)
		} else if hostgroupID == 3 {
			r = append(r, hostname)
		}
	}

	// Update the global variables
	Writers = &w
	Readers = &r
}
