package models

type User struct {
	ID             int    `json:"id" db:"id"`
	PassportNumber string `json:"passport_number" db:"passport_number"`
	Name           string `json:"name" db:"name"`
	Patronymic     string `json:"patronymic" db:"patronymic"`
	Surname        string `json:"surname" db:"surname"`
	Address        string `json:"addr" db:"addr"`
}

type QueryParams struct {
	ID             string `form:"id"`
	Name           string `form:"name"`
	Surname        string `form:"surname"`
	Patronymic     string `form:"patronymic"`
	PassportNumber string `form:"passport_number"`
	Address        string `form:"address"`
	Limit          int    `form:"limit"`
	Offset         int    `form:"offset"`
}

type UpdateUserInput struct {
	PassportNumber *string `json:"passport_number"`
	Address        *string `json:"addr"`
}
