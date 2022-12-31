package service

import "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"

type ProblemAddParameterIterator interface {
	Next() (domain.ProblemAddParameter, error)
}
