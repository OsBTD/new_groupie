package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Artists struct {
	Image          string   `json:"image"`
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	CreationDate   int      `json:"creationDate"`
	FirstAlbum     string   `json:"firstAlbum"`
	RelationsURL   string   `json:"relations"`
	DatesLocations Relations
}

type Relations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func fetchData(url string, target interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

func main() {
	artists := []Artists{}
	err := fetchData("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		fmt.Println("Error fetching artists")
		return
	}

	for i := 0; i < len(artists); i++ {
		var relation Relations
		err := fetchData(artists[i].RelationsURL, &relation)
		if err != nil {
			fmt.Println("Error fetching relation for artist", artists[i].ID)
		} else {
			artists[i].DatesLocations = relation
		}
	}
	Search := map[string][]Artists{}
	for _, entry := range artists {
		nameEntry := strings.ToLower(entry.Name)
		Search[nameEntry] = append(Search[nameEntry], entry)

		for _, member := range entry.Members {
			memberEntry := strings.ToLower(member)
			Search[memberEntry] = append(Search[memberEntry], entry)
		}

		for location := range entry.DatesLocations.DatesLocations {
			locationEntry := strings.ToLower(location)
			Search[locationEntry] = append(Search[locationEntry], entry)
		}

		creationEntry := strconv.Itoa(entry.CreationDate)
		Search[creationEntry] = append(Search[creationEntry], entry)

		firstalbumEntry := strings.ToLower(entry.FirstAlbum)
		Search[firstalbumEntry] = append(Search[firstalbumEntry], entry)
	}

	tmpl := template.Must(template.ParseFiles("templates/index3.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, artists)
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("templates"))))
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
