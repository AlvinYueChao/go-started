package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type City struct {
	ID          int
	Name        string
	CountryCode string
	District    string
	Population  int
}

var ErrDbProvider = "failed to get db provider: %v"

func getCitiesByCountryCode(countryCode string) ([]City, error) {
	var err error
	db, err := getDbProvider()
	if err != nil {
		return nil, fmt.Errorf(ErrDbProvider, err)
	}
	var cities []City
	var dbQuery string
	var rows *sql.Rows

	if countryCode == "" {
		dbQuery = "SELECT * FROM city"
		rows, err = db.Query(dbQuery)
	} else {
		dbQuery = "SELECT * FROM city WHERE CountryCode = ?"
		rows, err = db.Query(dbQuery, countryCode)
	}

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
		return nil, fmt.Errorf(ErrDbProvider, err)
	}

	return cities, nil
}

func getCityByID(id int) (City, error) {
	db, err := getDbProvider()
	var city City
	if err != nil {
		return city, fmt.Errorf(ErrDbProvider, err)
	}
	err = db.QueryRow("SELECT * FROM city WHERE ID = ?", id).Scan(&city.ID, &city.Name, &city.CountryCode, &city.District, &city.Population)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return city, fmt.Errorf("no city with ID %d", id)
		}
		return city, fmt.Errorf("failed to get city by ID %d: %v", id, err)
	}
	return city, nil
}

func addCity(city City) (int64, error) {
	db, err := getDbProvider()
	if err != nil {
		return 0, fmt.Errorf(ErrDbProvider, err)
	}
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
