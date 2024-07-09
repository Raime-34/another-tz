package schemas

// Для вывода при get запросе
type Task struct {
	Id   int64
	Cost int64
}

type TaskResult struct {
	Tasks []Task
}

// Для импорта данных из json'a при миграции
type Task4Load struct {
	Id     int64  `json:"peopleId"`
	StartT string `json:"startT"`
	EndT   string `json:"endT"`
}
