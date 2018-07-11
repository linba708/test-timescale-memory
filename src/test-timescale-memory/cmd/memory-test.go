package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"time"
	"database/sql"
	"strings"
)



func setup(db *sql.DB){


	db.Exec(`CREATE EXTENSION IF NOT EXISTS  timescaledb  CASCADE;`)


	//_, err := db.Exec(`DROP TABLE IF EXISTS trades;`)
	//if err != nil {
	//	panic(err)
	//}

	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS trades
(
  insref numeric(20) NOT NULL,
  tradeprice double precision,
  tradeyield double precision,
  tradequantity numeric(20,0) NULL,
  tradereference character varying(64) NOT NULL,
  tradecode integer NOT NULL,
  tradetype character varying(10),
  buyer character varying(64),
  seller character varying(64),
  executedside character(1),
  "canceltime" time without time zone,
  "time" time without time zone NOT NULL,
  date date NOT NULL,
  PRIMARY KEY (insref, date, tradereference)
);
	`)

	if err != nil {
		panic(err)
	}

	//will only create hypertable in
	//db.Exec("SELECT create_hypertable('trades', 'date', chunk_time_interval => interval '1 day');")

}



func main() {
	fmt.Println("Starting test")


	timescale := make(chan int, 100)
	timescaleDone := make(chan bool, 1)
	pg := make(chan int, 100)
	pgDone := make(chan bool, 1)


	pgTest := func(prog chan int, done chan bool, uri string ) {

		days := 55
		perDay := 500000
		batchSize := 10000
		startDate := time.Now()
		//startDate.Add(time.Hour * 24 * time.Duration(56))

		db, err := sql.Open("postgres", uri)
		setup(db)
		if err != nil {
			panic(err)
		}

		start := time.Now()

		for day := 1; day <= days ; day++  {
			currentDate := startDate.Add(time.Hour * 24 * time.Duration(day))
			dateAsString := currentDate.Format("2006-01-02")
			for insref := 0; insref < perDay; insref += batchSize {
				query := make([]string, 0, batchSize)
				insertString :=
					`INSERT INTO trades
					(insref,
					date,
					tradereference,
					tradecode,
					time)
					VALUES `
				for j := insref; j < perDay && j < insref + batchSize; j ++ {
					value := fmt.Sprintf("(%d, '%s', 'FA12347', 224, '00:00:00')", j, dateAsString)
					query = append(query, value)
				}
				q := insertString + strings.Join(query, ",") + ";"
				//fmt.Printf(q)


				_, err := db.Exec(q)
				if err != nil {
					fmt.Println("Error handling query",err)
				}
				if (insref*(day)) % (perDay/2) == 0{

					elapsed := time.Now().Sub(start)
					start = time.Now()
					opsps := float64(perDay/2) / elapsed.Seconds()
					prog <- int(opsps)
				}
			}
		}

		done <- true
	}


	go pgTest(timescale, timescaleDone, "postgresql://postgres:qwerty@timescale:5432/timescale?sslmode=disable")
	//go pgTest(pg, pgDone, "postgresql://postgres:qwerty@postgres:5432/postgres?sslmode=disable")



	for{


		select {
		case opsps := <- timescale:
			fmt.Printf("TS: %d ops/sec\n", opsps)
		case opsps := <- pg:
			fmt.Printf("PG: %d ops/sec\n", opsps)
		case <-timescaleDone:
			fmt.Println("TS Done")
		case <-pgDone:
			fmt.Println("PG Done")
		}

	}





}

