package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
)

func getAllCities(c *gin.Context) {
	allCities, err := getCitiesByCountryCode("")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, allCities)
	}
}

func getAllCitiesByCountryCode(c *gin.Context) {
	countryCode := c.Param("countryCode")
	allCities, err := getCitiesByCountryCode(countryCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, allCities)
	}
}

func getSpecificCityByID(c *gin.Context) {
	cityID, err := strconv.Atoi(c.Param("cityID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	city, err := getCityByID(cityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, city)
	}
}

func createCity(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
	}
	var city City
	err = json.Unmarshal(bytes, &city)
	if err != nil {
		log.Println(err)
	}
	cityId, err := addCity(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"cityId": cityId})
	}
}

func main() {
	router := gin.Default()
	router.GET("/cities", getAllCities)
	router.GET("/cities/countryCode/:countryCode", getAllCitiesByCountryCode)
	router.GET("/cities/:cityID", getSpecificCityByID)
	router.POST("/cities", createCity)
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
