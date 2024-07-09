package handlers

import (
	"anotherTZ/components/schemas"
	"anotherTZ/components/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

// GetPeople godoc
// @Summary Get list of people
// @Description Get all people or filter by query parameters
// @Tags people
// @Accept  json
// @Produce  json
// @Param surname query string false "Surname"
// @Param name query string false "Name"
// @Param patronymic query string false "Patronymic"
// @Param address query string false "Address"
// @Param passport_serie query string false "Passport Serie"
// @Param passport_number query string false "Passport Number"
// @Success 200 {array} schemas.PeopleInfo
// @Failure 500 {string} string "Internal Server Error"
// @Router /people [get]
func GetPeople(w http.ResponseWriter, r *http.Request) {
	db, err := utils.ConnectToDB()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	defer db.Close()

	params := utils.GetAllParams(r.URL.Query())
	query := utils.ComposeQuery("SELECT * FROM public.people WHERE ", params, " AND", false)
	rows, err := db.Query(query)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	for rows.Next() {
		var (
			peopleInfo schemas.PeopleInfo
		)
		err = rows.Scan(
			&peopleInfo.Id,
			&peopleInfo.Surname,
			&peopleInfo.Name,
			&peopleInfo.Patronymic,
			&peopleInfo.Address,
			&peopleInfo.PassportSerie,
			&peopleInfo.PassportNumber,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}
		bytes, err := json.Marshal(peopleInfo)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}

		w.Write(bytes)
	}
}

// SetPeople godoc
// @Summary Add a new person
// @Description Add a new person to the database
// @Tags people
// @Accept  json
// @Produce  text/plain
// @Success 201 {string} string "User created successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /people [post]
func SetPeople(w http.ResponseWriter, r *http.Request) {
	db, err := utils.ConnectToDB()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	defer db.Close()

	var params = make(map[string]string)
	err = utils.ParseJSON(r.Body, params)
	passportInfo := strings.Fields(params["passportNumber"])

	if len(passportInfo) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	_, err = db.Query(
		"INSERT INTO "+
			"public.people "+
			"(surname, name, patronymic, address, passport_serie, passport_number) "+
			"VALUES "+
			"($1, $2, $3, $4, $5, $6)",
		params["surname"], params["name"], params["patronymic"], params["address"], passportInfo[0], passportInfo[1],
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	} else {
		fmt.Println("Произведена запись пользователя")
	}
}

// DeletePeople godoc
// @Summary Delete a person
// @Description Delete a person from the database by passport number
// @Tags people
// @Accept  json
// @Produce  json
// @Param passport_number query string true "Passport Number (format: 'serie number')"
// @Success 200 {string} string "User deleted successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /people [delete]
func DeletePeople(w http.ResponseWriter, r *http.Request) {
	db, err := utils.ConnectToDB()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	defer db.Close()

	params := r.URL.Query().Get("passport_number")
	if params == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	passportInfo := strings.Fields(params)

	if len(passportInfo) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Query(
		"DELETE FROM public.people WHERE passport_serie = $1 AND passport_number = $2",
		passportInfo[0], passportInfo[1],
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	} else {
		fmt.Println("Произведено удаление пользователя")
	}
}

// UpdatePeople godoc
// @Summary Update a person
// @Description Update person details in the database
// @Tags people
// @Accept  json
// @Produce  json
// @Param passport_number query string true "Passport Number (format: 'serie number')"
// @Param surname query string false "Surname"
// @Param name query string false "Name"
// @Param patronymic query string false "Patronymic"
// @Param address query string false "Address"
// @Success 200 {string} string "User updated successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /people [put]
func UpdatePeople(w http.ResponseWriter, r *http.Request) {
	db, err := utils.ConnectToDB()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	defer db.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	params := utils.GetAllParams(r.URL.Query())
	devidedPassportInfo := strings.Fields(params["passport_number"])

	query := utils.ComposeQuery("UPDATE public.people SET ", params, ",", true)
	query += fmt.Sprintf(" WHERE passport_serie = '%s' AND passport_number = '%s'", devidedPassportInfo[0], devidedPassportInfo[1])
	_, err = db.Query(query)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	} else {
		fmt.Println("Произведено обновление пользователя")
	}
}
