package utils

import (
	"anotherTZ/components/schemas"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ParseJSON(body io.ReadCloser, out map[string]string) error {
	bytes, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &out)
	if err != nil {
		panic(err)
	}

	return nil
}

func ConnectToDB() (*sql.DB, error) {
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		port,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("postgres", psqlconn)
	return db, err
}

// Для получения переменных среды
func getConnectionData() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func MakeMigration() {
	getConnectionData()
	db, err := ConnectToDB()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Println("Миграция завершена")
	loadPeopleData(db)
	loadTaskData(db)
}

// Подгружает данные для таблицы people
func loadPeopleData(db *sql.DB) {
	jsonFile, err := os.Open("db/data/people.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var people []schemas.PeopleInfo
	json.Unmarshal(byteValue, &people)

	for _, person := range people {
		_, err := db.Exec("INSERT INTO public.people (surname, name, patronymic, address, passport_serie, passport_number) VALUES ($1, $2, $3, $4, $5, $6)",
			person.Surname, person.Name, person.Patronymic, person.Address, person.PassportSerie, person.PassportNumber)
		if err != nil {
			decoder := charmap.Windows1251.NewDecoder()
			out, _ := decoder.String(err.Error())
			log.Println(out)
			log.Fatal(err)
		}
	}

	fmt.Println("Данные добавлены в БД (пользователи)")
}

// Подгружает данные для таблицы people
func loadTaskData(db *sql.DB) {
	jsonFile, err := os.Open("db/data/tasks.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var tasks []schemas.Task4Load
	json.Unmarshal(byteValue, &tasks)

	for _, task := range tasks {
		_, err := db.Exec("INSERT INTO public.task (peopleId, startT, endT) VALUES ($1, $2, $3)",
			task.Id, task.StartT, task.EndT)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Данные добавлены в БД (задачи)")
}

// Извлекает параметры из URL
func GetAllParams(v url.Values) map[string]string {
	var result = make(map[string]string)

	result["surname"] = v.Get("surname")
	result["name"] = v.Get("name")
	result["patronymic"] = v.Get("patronymic")
	result["address"] = v.Get("address")

	passportInfo := strings.Fields(v.Get("passport_number"))
	if len(passportInfo) == 2 {
		result["passportSerie"] = passportInfo[0]
		result["passportNumber"] = passportInfo[1]
	}

	return result
}

// Собирает запрос по передеанным параметрам
func ComposeQuery(baseQ string, params map[string]string, separator string, ignorePassportInfo bool) string {
	result := []byte(baseQ)

	for k, v := range params {
		if k == "page" {
			continue
		}

		if (k == "passportSerie" || k == "passportNumber") && ignorePassportInfo {
			continue
		}

		if v != "" {
			result = append(result, []byte(fmt.Sprintf("%s = '%s'%s ", k, v, separator))...)
		}
	}

	result = bytes.TrimRight(result, ", ")
	result = bytes.TrimRight(result, " AND ")

	result = append(result, []byte("LIMIT 10 ")...)

	page, err := strconv.Atoi(params["page"])
	if err == nil {
		result = append(result, fmt.Sprintf("OFFSET %d", (page-1)*10)...)
	}

	return string(result)
}
