//go:generate mockery --output mock --name Problem
package service

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type ProblemFeature interface {
	// FindAudioByAudioID(ctx context.Context, audioID domain.AudioID) (Audio, error)
}

type Problem interface {
	domain.ProblemModel
	ProblemFeature
}

type problem struct {
	domain.ProblemModel
	synthesizerClient SynthesizerClient
}

func NewProblem(synthesizerClient SynthesizerClient, problemModel domain.ProblemModel) (Problem, error) {
	m := &problem{
		ProblemModel:      problemModel,
		synthesizerClient: synthesizerClient,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

// func (s *problem) FindAudioByAudioID(ctx context.Context, audioID domain.AudioID) (Audio, error) {
// 	logger := log.FromContext(ctx)
// 	if strconv.Itoa(int(audioID)) != s.GetProperties(ctx)["audioId"] {
// 		logger.Debugf("properties: %+v", s.GetProperties(ctx))
// 		logger.Warnf("audioID: %d, %s", audioID, s.GetProperties(ctx)["audioId"])
// 		message := "invalid audio id"
// 		return nil, domain.NewPluginError(domain.ErrorType(domain.ErrorTypeClient), message, []string{message}, libD.ErrInvalidArgument)
// 	}

// 	return s.synthesizerClient.FindAudioByAudioID(ctx, audioID)
// }

type ProblemWithResults interface {
	domain.ProblemModel
	GetResults() []bool
	GetLevel() int
}

type problemWithResults struct {
	domain.ProblemModel
	results []bool
	level   int
}

func NewProblemWithResults(problem domain.ProblemModel, results []bool, level int) ProblemWithResults {
	return &problemWithResults{
		ProblemModel: problem,
		results:      results,
		level:        level,
	}
}

func (m *problemWithResults) GetResults() []bool {
	return m.results
}

func (m *problemWithResults) GetLevel() int {
	return m.level
}
