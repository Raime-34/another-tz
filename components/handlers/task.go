package handlers

import (
	"anotherTZ/components/schemas"
	"anotherTZ/components/utils"
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"
)

// SetTask godoc
// @Summary Create a new task
// @Description Create a new task for a specific person
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param peopleID query string true "ID of the person"
// @Success 201 {string} string "Task created successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /task [post]
func SetTask(w http.ResponseWriter, r *http.Request) {
	db, err := utils.ConnectToDB()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	defer db.Close()

	peopleID := r.URL.Query().Get("peopleID")
	if peopleID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t := time.Now().Format("2006-01-02 15:04:05")

	rows, err := db.Query(
		"INSERT INTO "+
			"public.task "+
			"(peopleId, startT) "+
			"VALUES "+
			"($1, $2)",
		peopleID, t,
	)

	if !rows.Next() {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		fmt.Println("Произведена запись задачи")
	}
}

// CloseTask godoc
// @Summary Close a task
// @Description Close a task for a specific person by updating the end time
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param peopleID query string true "ID of the person"
// @Success 200 {string} string "Task closed successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /task [put]
func CloseTask(w http.ResponseWriter, r *http.Request) {
	db, err := utils.ConnectToDB()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	defer db.Close()

	peopleID := r.URL.Query().Get("peopleID")
	if peopleID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t := time.Now().Format("2006-01-02 15:04:05")

	rows, err := db.Query(
		"UPDATE "+
			"public.task "+
			"SET "+
			"endT = $1 "+
			"WHERE "+
			"peopleId = $2",
		t, peopleID,
	)

	if !rows.Next() {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		fmt.Println("Произведено завершение задачи")
	}
}

// GetTask godoc
// @Summary Get tasks
// @Description Get tasks for a specific person within a given time interval
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param peopleID query string true "ID of the person"
// @Param interval query string false "Time interval (e.g., '1 day', '1 week')"
// @Success 200 {object} schemas.TaskResult
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /task [get]
func GetTask(w http.ResponseWriter, r *http.Request) {
	db, err := utils.ConnectToDB()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	defer db.Close()

	peopleID := r.URL.Query().Get("peopleID")
	interval := r.URL.Query().Get("interval")

	rows, err := db.Query(
		"SELECT "+
			"id, startT, endT "+
			"FROM "+
			"public.task "+
			"WHERE "+
			"peopleId = $1 ",
		peopleID,
	)

	result := []schemas.Task{}
	var leftBorder time.Time
	rigthBorder := time.Now()

	switch interval {
	case "1 day":
		leftBorder = rigthBorder.Add(-24 * time.Hour)
	case "1 week":
		leftBorder = rigthBorder.Add(-7 * 24 * time.Hour)
	default:
		leftBorder = rigthBorder.Add(-3 * time.Hour)
	}

	for rows.Next() {
		var (
			id    int64
			start time.Time
			end   time.Time
		)

		rows.Scan(&id, &start, &end)
		if start.Before(leftBorder) || end.After(rigthBorder) {
			continue
		}
		result = append(result, schemas.Task{
			Id:   id,
			Cost: (end.Unix() - start.Unix()) / 60 * 1000,
		})
	}

	slices.SortFunc(
		result,
		func(a, b schemas.Task) int {
			return cmp.Compare(b.Cost, a.Cost)
		},
	)

	tasks := schemas.TaskResult{
		Tasks: result,
	}

	bytes, err := json.Marshal(tasks)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	w.Write(bytes)
}
