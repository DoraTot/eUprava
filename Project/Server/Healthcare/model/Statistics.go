package model

type ChildStatistics struct {
	TotalRecords        int `json:"totalRecords"`
	TotalAbsent         int `json:"totalAbsent"`
	TotalPresent        int `json:"totalPresent"`
	TotalJustified      int `json:"totalJustified"`
	MonthlyAppointments int `json:"monthlyAppointments"`
	YearlyAppointments  int `json:"yearlyAppointments"`
}
