package model

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AboutUserData interface {
}

type AboutCustomer struct {
	Capital         int `json:"capital"`
	AssignedLoaders []Loader
}

type AboutLoader struct {
	Weight     int  `json:"weight"`
	Wage       int  `json:"wage"`
	Alcoholism bool `json:"alcoholism"`
	Fatigue    int  `json:"fatigue"`
}

type TaskToGen struct {
	TaskName string `json:"task_name" db:"task_name"`
	Weight   int    `json:"weight" db:"weight"`
	Status   bool   `json:"status" db:"status"`
}

type StartTask struct {
	TaskID          int   `json:"task_name"`
	AssignedLoaders []int `json:"loaders"`
}
