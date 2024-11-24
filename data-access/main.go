package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

type City struct {
	ID          int
	Name        string
	CountryCode string
	District    string
	Population  int
}

func getCitiesByCountryCode(db *sql.DB, countryCode string) ([]City, error) {
	var cities []City

	rows, err := db.Query("SELECT * FROM city WHERE CountryCode = ?", countryCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get cities by countryCode %q: %v", countryCode, err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)
	for rows.Next() {
		var city City
		if err := rows.Scan(&city.ID, &city.Name, &city.CountryCode, &city.District, &city.Population); err != nil {
			return nil, fmt.Errorf("failed to scan city: %v", err)
		}
		cities = append(cities, city)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %v", err)
	}

	return cities, nil
}

func getCityByID(db *sql.DB, id int) (City, error) {
	var city City
	err := db.QueryRow("SELECT * FROM city WHERE ID = ?", id).Scan(&city.ID, &city.Name, &city.CountryCode, &city.District, &city.Population)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return city, fmt.Errorf("no city with ID %d", id)
		}
		return city, fmt.Errorf("failed to get city by ID %d: %v", id, err)
	}
	return city, nil
}

func addCity(db *sql.DB, city City) (int64, error) {
	//if city == nil {
	//	return 0, fmt.Errorf("city cannot be nil")
	//}
	result, err := db.Exec("INSERT INTO city (Name, CountryCode, District, Population) VALUES (?, ?, ?, ?)", city.Name, city.CountryCode, city.District, city.Population)
	if err != nil {
		return 0, fmt.Errorf("failed to add city: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %v", err)
	}
	return id, nil
}

func main() {
	var db *sql.DB
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "world",
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalln(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatalln(pingErr)
	}
	fmt.Println("Connected to the database!")

	cities, err := getCitiesByCountryCode(db, "AIA")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(cities)

	city, err := getCityByID(db, 1)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(city)

	cityId, err := addCity(db, City{
		Name:        "Hirano Yamato",
		CountryCode: "AIA",
		District:    "v6WRpXhnKj",
		Population:  587,
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Id of new added city: %v\n", cityId)
}
