package schemas

type PeopleInfo struct {
	Id             string `json:"id"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
	PassportSerie  string `json:"passport_serie"`
	PassportNumber string `json:"passport_number"`
}
