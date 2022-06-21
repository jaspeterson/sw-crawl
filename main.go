package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Film struct {
	Title string `json:"title"`
	Crawl string `json:"opening_crawl"`
}

type CrawlResponse struct {
	Title string `json:"film_name"`
	Crawl []string `json:"opening_crawl"`
}

type Manager struct {
	randomGenerator *rand.Rand
}

func main() {
	seed := rand.NewSource(time.Now().UnixNano())
	mngr := Manager{
		randomGenerator: rand.New(seed),
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/random_crawl", mngr.crawlHandler)

	fmt.Println("Running star wars crawl server on port 8080...")
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("Error encountered running server. Aborting... | %s\n", err.Error())
	}
}

func(m* Manager) crawlHandler(ctx *gin.Context) {
	n := m.randomGenerator.Intn(6)+1

	resp, err := http.Get("https://swapi.dev/api/films/"+ strconv.Itoa(n) +"/")
	if err != nil {
		fmt.Printf("Error encountered accessing swapi: %s\n", err.Error())
		ctx.JSON(http.StatusInternalServerError, "Cannot access swapi at this time")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error encountered reading swapi response: %s\n", err.Error())
		ctx.JSON(http.StatusInternalServerError, "Cannot access swapi at this time")
		return
	}

	var film Film
	err = json.Unmarshal(body, &film)
	if err != nil {
		fmt.Printf("Error encountered reading swapi response: %s\n", err.Error())
		ctx.JSON(http.StatusInternalServerError, "Cannot access swapi at this time")
		return
	}

	crawl := strings.Split(film.Crawl, "\r\n")
	
	ctx.JSON(http.StatusOK, CrawlResponse{
		Title: film.Title,
		Crawl: crawl,
	})
}