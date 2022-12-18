//go:generate mockery --output mock --name ProcessorFactory
package service

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type ProcessorFactory interface {
	NewProblemAddProcessor(problemType domain.ProblemTypeName) (ProblemAddProcessor, error)

	NewProblemUpdateProcessor(problemType domain.ProblemTypeName) (ProblemUpdateProcessor, error)

	NewProblemRemoveProcessor(problemType domain.ProblemTypeName) (ProblemRemoveProcessor, error)

	NewProblemImportProcessor(problemType domain.ProblemTypeName) (ProblemImportProcessor, error)

	NewProblemQuotaProcessor(problemType domain.ProblemTypeName) (ProblemQuotaProcessor, error)
}

type processorFactrory struct {
	addProcessors    map[domain.ProblemTypeName]ProblemAddProcessor
	updateProcessors map[domain.ProblemTypeName]ProblemUpdateProcessor
	removeProcessors map[domain.ProblemTypeName]ProblemRemoveProcessor
	importProcessors map[domain.ProblemTypeName]ProblemImportProcessor
	quotaProcessors  map[domain.ProblemTypeName]ProblemQuotaProcessor
}

func NewProcessorFactory(addProcessors map[domain.ProblemTypeName]ProblemAddProcessor, updateProcessors map[domain.ProblemTypeName]ProblemUpdateProcessor, removeProcessors map[domain.ProblemTypeName]ProblemRemoveProcessor, importProcessors map[domain.ProblemTypeName]ProblemImportProcessor, quotaProcessors map[domain.ProblemTypeName]ProblemQuotaProcessor) ProcessorFactory {
	return &processorFactrory{
		addProcessors:    addProcessors,
		updateProcessors: updateProcessors,
		removeProcessors: removeProcessors,
		importProcessors: importProcessors,
		quotaProcessors:  quotaProcessors,
	}
}

func (f *processorFactrory) NewProblemAddProcessor(problemType domain.ProblemTypeName) (ProblemAddProcessor, error) {
	processor, ok := f.addProcessors[problemType]
	if !ok {
		return nil, liberrors.Errorf("NewProblemAddProcessor not found. problemType: %s", problemType)
	}
	return processor, nil
}

func (f *processorFactrory) NewProblemUpdateProcessor(problemType domain.ProblemTypeName) (ProblemUpdateProcessor, error) {
	processor, ok := f.updateProcessors[problemType]
	if !ok {
		return nil, liberrors.Errorf("NewProblemUpdateProcessor not found. problemType: %s", problemType)
	}
	return processor, nil
}

func (f *processorFactrory) NewProblemRemoveProcessor(problemType domain.ProblemTypeName) (ProblemRemoveProcessor, error) {
	processor, ok := f.removeProcessors[problemType]
	if !ok {
		return nil, liberrors.Errorf("NewProblemRemoveProcessor not found. problemType: %s", problemType)
	}
	return processor, nil
}

func (f *processorFactrory) NewProblemImportProcessor(problemType domain.ProblemTypeName) (ProblemImportProcessor, error) {
	processor, ok := f.importProcessors[problemType]
	if !ok {
		return nil, liberrors.Errorf("NewProblemImportProcessor not found. problemType: %s", problemType)
	}
	return processor, nil
}

func (f *processorFactrory) NewProblemQuotaProcessor(problemType domain.ProblemTypeName) (ProblemQuotaProcessor, error) {
	processor, ok := f.quotaProcessors[problemType]
	if !ok {
		return nil, liberrors.Errorf("NewProblemQuotaProcessor not found. problemType: %s", problemType)
	}
	return processor, nil
}
