package english_word

import (
	"context"
	"errors"
	"strconv"

	appD "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	appS "github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/data"
	pluginCommonDomain "github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/domain"
	pluginEnglishDomain "github.com/kujilabo/cocotola/cocotola-api/src/plugin/english/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

func CreateDemoWorkbook(ctx context.Context, studentService appS.Student) error {
	if err := CreateWorkbook(ctx, studentService, "Example", pluginCommonDomain.PosOther, []string{"butcher", "bakery", "library", "bookstore", "drugstore", "restaurant", "garage", "barbershop", "bank", "market"}); err != nil {
		return err
	}
	return nil
}

func Create20NGSLWorkbook(ctx context.Context, studentService appS.Student) error {
	if err := CreateWorkbook(ctx, studentService, "NGSL-20", pluginCommonDomain.PosOther, []string{
		"know",
		"more",
		"get",
		"who",
		"like",
		"when",
		"think",
		"make",
		"time",
		"see",
		"what",
		"up",
		"some",
		"other",
		"out",
		"good",
		"people",
		"year",
		"take",
		"no",
		"well",
		"because",
		"very",
		"just",
		"come",
		"could",
		"work",
		"use",
		"than",
		"now",
	}); err != nil {
		return err
	}
	return nil
}

func Create300NGSLWorkbook(ctx context.Context, studentService appS.Student) error {
	if err := CreateWorkbook(ctx, studentService, "NGSL-300", pluginCommonDomain.PosOther, []string{
		"know",
		"more",
		"get",
		"who",
		"like",
		"when",
		"think",
		"make",
		"time",
		"see",
		"what",
		"up",
		"some",
		"other",
		"out",
		"good",
		"people",
		"year",
		"take",
		"no",
		"well",
		"because",
		"very",
		"just",
		"come",
		"could",
		"work",
		"use",
		"than",
		"now",
		"then",
		"also",
		"into",
		"only",
		"look",
		"want",
		"give",
		"first",
		"new",
		"way",
		"find",
		"over",
		"any",
		"after",
		"day",
		"where",
		"thing",
		"most",
		"should",
		"need",
		"much",
		"right",
		"how",
		"back",
		"mean",
		"even",
		"may",
		"here",
		"many",
		"such",
		"last",
		"child",
		"tell",
		"really",
		"call",
		"before",
		"company",
		"through",
		"down",
		"show",
		"life",
		"man",
		"change",
		"place",
		"long",
		"between",
		"feel",
		"too",
		"still",
		"problem",
		"write",
		"same",
		"lot",
		"great",
		"try",
		"leave",
		"number",
		"both",
		"own",
		"part",
		"point",
		"little",
		"help",
		"ask",
		"meet",
		"start",
		"talk",
		"something",
		"put",
		"another",
		"become",
		"interest",
		"country",
		"old",
		"each",
		"school",
		"late",
		"high",
		"different",
		"off",
		"next",
		"end",
		"live",
		"why",
		"while",
		"world",
		"week",
		"play",
		"might",
		"must",
		"home",
		"never",
		"include",
		"course",
		"house",
		"report",
		"group",
		"case",
		"woman",
		"around",
		"book",
		"family",
		"seem",
		"let",
		"again",
		"kind",
		"keep",
		"hear",
		"system",
		"every",
		"question",
		"during",
		"always",
		"big",
		"set",
		"small",
		"study",
		"follow",
		"begin",
		"important",
		"since",
		"run",
		"under",
		"turn",
		"few",
		"bring",
		"early",
		"hand",
		"state",
		"move",
		"money",
		"fact",
		"however",
		"area",
		"provide",
		"name",
		"read",
		"friend",
		"month",
		"large",
		"business",
		"without",
		"information",
		"open",
		"order",
		"government",
		"word",
		"issue",
		"market",
		"pay",
		"build",
		"hold",
		"service",
		"against",
		"believe",
		"second",
		"though",
		"yes",
		"love",
		"increase",
		"job",
		"plan",
		"result",
		"away",
		"example",
		"happen",
		"offer",
		"young",
		"close",
		"program",
		"lead",
		"buy",
		"understand",
		"thank",
		"far",
		"today",
		"hour",
		"student",
		"face",
		"hope",
		"idea",
		"cost",
		"less",
		"room",
		"until",
		"reason",
		"form",
		"spend",
		"head",
		"car",
		"learn",
		"level",
		"person",
		"experience",
		"once",
		"member",
		"enough",
		"bad",
		"city",
		"night",
		"able",
		"support",
		"whether",
		"line",
		"present",
		"side",
		"quite",
		"although",
		"sure",
		"term",
		"least",
		"age",
		"low",
		"speak",
		"within",
		"process",
		"public",
		"often",
		"train",
		"possible",
		"actually",
		"rather",
		"view",
		"together",
		"consider",
		"price",
		"parent",
		"hard",
		"party",
		"local",
		"control",
		"already",
		"concern",
		"product",
		"lose",
		"story",
		"almost",
		"continue",
		"stand",
		"whole",
		"yet",
		"rate",
		"care",
		"expect",
		"effect",
		"sort",
		"ever",
		"anything",
		"cause",
		"fall",
		"deal",
		"water",
		"send",
		"allow",
		"soon",
		"watch",
		"base",
		"probably",
		"suggest",
		"past",
		"power",
		"test",
		"visit",
		"center",
		"grow",
		"nothing",
		"return",
		"mother",
		"walk",
		"matter",
	}); err != nil {
		return err
	}
	return nil
}
func CreateWorkbook(ctx context.Context, student appS.Student, workbookName string, pos pluginCommonDomain.WordPos, words []string) error {
	logger := log.FromContext(ctx)

	workbookProperties := map[string]string{
		"audioEnabled": "false",
	}
	param, err := appD.NewWorkbookAddParameter(pluginEnglishDomain.EnglishWordProblemType, workbookName, appD.Lang2JA, "", workbookProperties)
	if err != nil {
		return liberrors.Errorf("failed to NewWorkbookAddParameter. err: %w", err)
	}

	workbook, err := data.CreateWorkbookIfNotExists(ctx, student, workbookName, param)
	if err != nil {
		return liberrors.Errorf("createWorkbookIfNotExists. err: %w", err)
	}

	problems, err := workbook.FindAllProblems(ctx, student)
	if err != nil {
		return liberrors.Errorf("createWorkbookIfNotExists. err: %w", err)
	}

	problemMap := make(map[string]struct{})
	for _, problem := range problems.GetResults() {
		properties := problem.GetProperties(ctx)
		text1, ok := properties["text"]
		if !ok {
			continue
		}
		text2, ok := text1.(string)
		if !ok {
			continue
		}
		problemMap[text2] = struct{}{}
	}

	for i, word := range words {
		// skip if the word is already registered
		if _, ok := problemMap[word]; ok {
			logger.Infof("Skip %s", word)
			continue
		}

		properties := map[string]string{
			"number": strconv.Itoa(i + 1),
			"text":   word,
			"lang2":  "ja",
			"pos":    strconv.Itoa(int(pos)),
		}
		param, err := appD.NewProblemAddParameter(workbook.GetWorkbookID(), properties)
		if err != nil {
			return liberrors.Errorf("failed to NewProblemAddParameter. err: %w", err)
		}

		added, _, _, err := workbook.AddProblem(ctx, student, param)
		if err != nil && !errors.Is(err, appS.ErrProblemAlreadyExists) {
			return liberrors.Errorf("AddProblem. err: %w", err)
		}
		logger.Infof("problemIDs: %v", added)
	}

	logger.Infof("Example %d", workbook.GetWorkbookID())
	return nil
}
