package model

type User struct {
	UserID   int    `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	UserRole string `json:"user_role" db:"user_role"`
}

type Customer struct {
	UserID  int `json:"user_id" db:"user_id"`
	Capital int `json:"capital" db:"capital"`
}

type Loader struct {
	UserID     int  `json:"user_id" db:"user_id"`
	MaxWeight  int  `json:"max_weight" db:"max_weight"`
	Alcoholism bool `json:"alcoholism" db:"alcoholism"`
	Fatigue    int  `json:"fatigue" db:"fatigue"`
	Wage       int  `json:"wage" db:"wage"`
}

type Task struct {
	TaskID     int    `json:"task_id" db:"task_id"`
	TaskName   string `json:"task_name" db:"task_name"`
	Weight     int    `json:"weight" db:"weight"`
	Status     bool   `json:"status" db:"status"`
	CustomerID int    `json:"customer_id" db:"customer_id"`
}

type AssignedLoader struct {
	LoaderID int `json:"loader_id" db:"loader_id"`
	TaskID   int `json:"task_id" db:"task_id"`
}

type Game struct {
	GameID     int  `json:"game_id" db:"game_id"`
	UserID     int  `json:"user_id" db:"user_id"`
	GameResult bool `json:"game_result" db:"game_result"`
}
