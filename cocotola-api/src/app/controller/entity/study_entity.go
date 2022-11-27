package entity

import "time"

type StudyResultParameter struct {
	Result   bool `json:"result"`
	Mastered bool `json:"mastered"`
}

type StudyRecord struct {
	ProblemID      uint       `json:"problemId"`
	Level          int        `json:"level"`
	ResultPrev1    bool       `json:"resultPrev1"`
	Mastered       bool       `json:"mastered"`
	LastAnsweredAt *time.Time `json:"lastAnsweredAt"`
}

type StudyRecords struct {
	Records []*StudyRecord `json:"records" validate:"dive"`
}

type IntValue struct {
	Value int `json:"value"`
}
