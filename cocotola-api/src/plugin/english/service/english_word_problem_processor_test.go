package service_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	appD "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	appDM "github.com/kujilabo/cocotola/cocotola-api/src/app/domain/mock"
	appS "github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	appSM "github.com/kujilabo/cocotola/cocotola-api/src/app/service/mock"
	pluginD "github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/domain"
	pluginDM "github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/domain/mock"
	pluginSM "github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/service/mock"
	"github.com/kujilabo/cocotola/cocotola-api/src/plugin/english/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/plugin/english/service"
	libD "github.com/kujilabo/cocotola/lib/domain"
)

var anythingOfContext = mock.MatchedBy(func(_ context.Context) bool { return true })

func englishWordProblemProcessor_Init(t *testing.T) (
	operator *appDM.StudentModel,
	workbookModel *appDM.WorkbookModel,
	rf *appSM.RepositoryFactory,
	problemRepo *appSM.ProblemRepository,
	englishWordProblemProcessor service.EnglishWordProblemProcessor,
	synthesizerClient *appSM.SynthesizerClient,
	translatorClient *pluginSM.TranslatorClient,
	tatoebaClient *pluginSM.TatoebaClient) {

	synthesizerClient = new(appSM.SynthesizerClient)
	translatorClient = new(pluginSM.TranslatorClient)
	tatoebaClient = new(pluginSM.TatoebaClient)
	operator = new(appDM.StudentModel)
	problemRepo = new(appSM.ProblemRepository)
	rf = new(appSM.RepositoryFactory)
	rf.On("NewProblemRepository", anythingOfContext, domain.EnglishWordProblemType).Return(problemRepo, nil)
	workbookModel = new(appDM.WorkbookModel)
	englishWordProblemProcessor = service.NewEnglishWordProblemProcessor(synthesizerClient, translatorClient, tatoebaClient, nil)
	return
}

func testNewTranslation(pos pluginD.WordPos, translated string) *pluginDM.Translation {
	translation := new(pluginDM.Translation)
	translation.On("GetPos").Return(pos)
	translation.On("GetTranslated").Return(translated)
	return translation
}

func Test_englishWordProblemProcessor_AddProblem_singleProblem_audioDisabled(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	operator, workbookModel, rf, problemRepo, processor, _, _, _ := englishWordProblemProcessor_Init(t)

	// given
	// - workbook
	workbookModel.On("GetProperties").Return(map[string]string{
		"audioEnabled": "false",
	})
	// - problemRepo
	problemRepo.On("AddProblem", anythingOfContext, operator, mock.Anything).Return(appD.ProblemID(100), nil)
	// when
	// - param
	param := new(appSM.ProblemAddParameter)
	param.On("GetWorkbookID").Return(appD.WorkbookID(1))
	param.On("GetIntProperty", "number").Return(2, nil)
	param.On("GetProperties").Return(map[string]string{
		"pos":        "6",
		"text":       "pen",
		"translated": "ペン",
		"lang2":      "ja",
	})
	added, updated, removed, err := processor.AddProblem(ctx, rf, operator, workbookModel, param)
	require.NoError(t, err)
	// then
	require.Equal(t, 1, len(added))
	require.Equal(t, 100, int(added[0]))
	require.Equal(t, 0, len(updated))
	require.Equal(t, 0, len(removed))
	paramCheck := mock.MatchedBy(func(p appS.ProblemAddParameter) bool {
		require.Equal(t, 1, int(p.GetWorkbookID()))
		require.Equal(t, "ペン", p.GetProperties()["translated"])
		require.Equal(t, "pen", p.GetProperties()["text"])
		require.Equal(t, "ja", p.GetProperties()["lang2"])
		require.Equal(t, "0", p.GetProperties()["audioId"])
		require.Equal(t, "6", p.GetProperties()["pos"])
		require.Len(t, p.GetProperties(), 5)
		return true
	})
	problemRepo.AssertCalled(t, "AddProblem", anythingOfContext, operator, paramCheck)
	problemRepo.AssertNumberOfCalls(t, "AddProblem", 1)
}

func Test_englishWordProblemProcessor_AddProblem_multipleProblem_audioDisabled(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	operator, workbookModel, rf, problemRepo, processor, _, translatorClient, _ := englishWordProblemProcessor_Init(t)

	// given
	// - workbook
	workbookModel.On("GetProperties").Return(map[string]string{
		"audioEnabled": "false",
	})
	// - problemRepo
	problemRepo.On("AddProblem", anythingOfContext, operator, mock.Anything).Return(appD.ProblemID(100), nil)
	// - translatorClient
	translatorClient.On("DictionaryLookup", anythingOfContext, appD.Lang2EN, appD.Lang2JA, "book").Return([]pluginD.Translation{
		testNewTranslation(pluginD.PosNoun, "本"),
		testNewTranslation(pluginD.PosVerb, "予約する"),
	}, nil)
	// when
	// - param
	param := new(appSM.ProblemAddParameter)
	param.On("GetWorkbookID").Return(appD.WorkbookID(1))
	param.On("GetIntProperty", "number").Return(2, nil)
	param.On("GetProperties").Return(map[string]string{
		"pos":        "99",
		"text":       "book",
		"translated": "",
		"lang2":      "ja",
	})
	added, updated, removed, err := processor.AddProblem(ctx, rf, operator, workbookModel, param)
	require.NoError(t, err)
	// then
	require.Equal(t, 2, len(added))
	require.Equal(t, 100, int(added[0]))
	require.Equal(t, 0, len(updated))
	require.Equal(t, 0, len(removed))
	problemRepo.AssertNumberOfCalls(t, "AddProblem", 2)
	{
		param := (problemRepo.Calls[0].Arguments[2]).(appS.ProblemAddParameter)
		assert.Equal(t, 1, int(param.GetWorkbookID()))
		assert.Equal(t, "本", param.GetProperties()["translated"])
		assert.Equal(t, "book", param.GetProperties()["text"])
		assert.Equal(t, "ja", param.GetProperties()["lang2"])
		assert.Equal(t, "0", param.GetProperties()["audioId"])
		assert.Equal(t, "6", param.GetProperties()["pos"])
	}
	{
		param := (problemRepo.Calls[1].Arguments[2]).(appS.ProblemAddParameter)
		assert.Equal(t, 1, int(param.GetWorkbookID()))
		assert.Equal(t, "予約する", param.GetProperties()["translated"])
		assert.Equal(t, "book", param.GetProperties()["text"])
		assert.Equal(t, "ja", param.GetProperties()["lang2"])
		assert.Equal(t, "0", param.GetProperties()["audioId"])
		assert.Equal(t, "9", param.GetProperties()["pos"])
	}
}

func Test_englishWordProblemProcessor_UpdateProblem(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	operator, workbookModel, rf, problemRepo, processor, _, _, _ := englishWordProblemProcessor_Init(t)

	// given
	// - workbook
	workbookModel.On("GetProperties").Return(map[string]string{
		"audioEnabled": "false",
	})
	// - problemRepo
	problemRepo.On("UpdateProblem", anythingOfContext, operator, mock.Anything, mock.Anything).Return(nil)
	// when
	// - idParam
	idParam := new(appSM.ProblemSelectParameter2)
	idParam.On("GetProblemID").Return(appD.ProblemID(1))
	// - param
	param := new(appSM.ProblemUpdateParameter)
	param.On("GetWorkbookID").Return(appD.WorkbookID(1))
	param.On("GetIntProperty", "number").Return(2, nil)
	param.On("GetProperties").Return(map[string]string{
		"pos":        "6",
		"text":       "pen",
		"translated": "ペン",
		"lang2":      "ja",
	})
	added, updated, removed, err := processor.UpdateProblem(ctx, rf, operator, workbookModel, idParam, param)
	require.NoError(t, err)
	// then
	require.Equal(t, 0, len(added))
	require.Equal(t, 1, len(updated))
	require.Equal(t, 0, len(removed))
	problemRepo.AssertNumberOfCalls(t, "UpdateProblem", 1)
	{
		param := (problemRepo.Calls[0].Arguments[3]).(appS.ProblemUpdateParameter)
		assert.Equal(t, "ペン", param.GetProperties()["translated"])
		assert.Equal(t, "pen", param.GetProperties()["text"])
		assert.Equal(t, "0", param.GetProperties()["audioId"])
		assert.Equal(t, "0", param.GetProperties()["sentenceId1"])
		assert.Len(t, param.GetProperties(), 4)
	}
}

func testNewProblemAddParameter_EnglishWord(properties map[string]string) appS.ProblemAddParameter {
	param := new(appSM.ProblemAddParameter)
	param.On("GetProperties").Return(properties)
	param.On("GetIntProperty", "number").Return(1, nil)
	return param
}

func TestNewEnglishWordProblemAddParemeter(t *testing.T) {
	t.Parallel()
	type args struct {
		param appS.ProblemAddParameter
	}
	tests := []struct {
		name    string
		args    args
		want    *service.EnglishWordProblemAddParemeter
		wantErr error
	}{
		{
			name: "pos is not defined",
			args: args{
				param: testNewProblemAddParameter_EnglishWord(map[string]string{}),
			},
			wantErr: libD.ErrInvalidArgument,
		},
		{
			name: "text is not defined",
			args: args{
				param: testNewProblemAddParameter_EnglishWord(map[string]string{
					"pos": "6",
				}),
			},
			wantErr: libD.ErrInvalidArgument,
		},
		{
			name: "lang2 is not defined",
			args: args{
				param: testNewProblemAddParameter_EnglishWord(map[string]string{
					"pos":  "6",
					"text": "pen",
				}),
			},
			wantErr: libD.ErrInvalidArgument,
		},
		{
			name: "lang2 format is invalid",
			args: args{
				param: testNewProblemAddParameter_EnglishWord(map[string]string{
					"pos":   "6",
					"text":  "pen",
					"lang2": "jpn",
				}),
			},
			wantErr: libD.ErrInvalidArgument,
		},
		{
			name: "parameter is valid",
			args: args{
				param: testNewProblemAddParameter_EnglishWord(map[string]string{
					"pos":   "6",
					"text":  "pen",
					"lang2": "ja",
				}),
			},
			want: &service.EnglishWordProblemAddParemeter{
				Number: 1,
				Pos:    pluginD.PosNoun,
				Text:   "pen",
				Lang2:  appD.Lang2JA,
			},
			wantErr: nil,
		},
		{
			name: "parameter is valid, translated is defined",
			args: args{
				param: testNewProblemAddParameter_EnglishWord(map[string]string{
					"pos":        "6",
					"text":       "pen",
					"lang2":      "ja",
					"translated": "ペン",
				}),
			},
			want: &service.EnglishWordProblemAddParemeter{
				Number:     1,
				Pos:        pluginD.PosNoun,
				Text:       "pen",
				Lang2:      appD.Lang2JA,
				Translated: "ペン",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := service.NewEnglishWordProblemAddParemeter(tt.args.param)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewEnglishWordProblemAddParemeter() err = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEnglishWordProblemAddParemeter() = %v, want %v", got, tt.want)
			}
		})
	}
}
