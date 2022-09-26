package main

import (
	"context"
	"day7/connection"
	"fmt"
	"html/template"
	"log"

	"math"
	"net/http"
	"strconv"

	"time"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	connection.ConnectionProject()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.PathPrefix("/node_modules/").Handler(http.StripPrefix("/node_modules/", http.FileServer(http.Dir("./node_modules"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/form-project", formProject).Methods("GET")
	route.HandleFunc("/add-project", addProject).Methods("POST")
	route.HandleFunc("/detail-project/{index}", detailProject).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/delete-project/{index}", deleteProject).Methods("GET")

	fmt.Println("server on")
	http.ListenAndServe("localhost:5000", route)
}

type Blog struct {
	Id          int
	NameProject string
	StarDate    string
	EndDate     string
	Duration    string
	Message     string
	Tech        []string
}

var dataBlog = []Blog{}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	var nameProject = r.PostForm.Get("projectName")
	var startDate = r.PostForm.Get("startDate")
	var endDate = r.PostForm.Get("endDate")
	var desc = r.PostForm.Get("Description")
	// var image = r.PostForm.Get("image")

	layout := "2006-01-02"
	date1, _ := time.Parse(layout, startDate)
	date2, _ := time.Parse(layout, endDate)

	hs := date2.Sub(date1).Hours()
	day, _ := math.Modf(hs / 24)
	bulan := int64(day / 30)
	tahun := int64(day / 365)

	var duration string

	if tahun > 0 {
		duration = strconv.FormatInt(tahun, 10) + " Year"
	} else if bulan > 0 {
		duration = strconv.FormatInt(bulan, 10) + " Month"
	} else {
		duration = fmt.Sprintf("%.0f", day) + " Day"
	}

	var tech []string
	for key, values := range r.Form {
		for _, value := range values {
			if key == "technologies" {
				tech = append(tech, value)
			}
		}
	}

	var newData = Blog{
		NameProject: nameProject,
		StarDate:    startDate,
		EndDate:     endDate,
		Duration:    duration,
		Message:     desc,
		Tech:        tech,
	}

	dataBlog = append(dataBlog, newData)
	fmt.Println(dataBlog)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("view/index.html")

	if err != nil {
		w.Write([]byte("messege: " + err.Error()))
		return
	}

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, description,technologies FROM tb_projects")
	fmt.Println(data)

	var result []Blog

	for data.Next() {

		var each = Blog{}
		err := data.Scan(&each.Id, &each.NameProject, &each.Message, &each.Tech)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	response := map[string]interface{}{
		"Blogs": result,
	}

	tmpl.Execute(w, response)
}

func detailProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("view/blog-detail.html")

	if err != nil {
		w.Write([]byte("messege: " + err.Error()))
		return
	}

	var ProjectDetail = Blog{}

	indexData, _ := strconv.Atoi(mux.Vars(r)["index"])

	for index, value := range dataBlog {
		if index == indexData {
			ProjectDetail = Blog{
				NameProject: value.NameProject,
				StarDate:    value.StarDate,
				EndDate:     value.EndDate,
				Duration:    value.Duration,
				Message:     value.Message,
				Tech:        value.Tech,
			}
		}
	}

	response := map[string]interface{}{
		"Blogs": ProjectDetail,
	}

	tmpl.Execute(w, response)
}

func formProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("view/add.html")

	if err != nil {
		w.Write([]byte("messege: " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("view/contact.html")

	if err != nil {
		w.Write([]byte("messege: " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	indexDelete, _ := strconv.Atoi(mux.Vars(r)["index"])

	dataBlog = append(dataBlog[:indexDelete], dataBlog[indexDelete+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}
