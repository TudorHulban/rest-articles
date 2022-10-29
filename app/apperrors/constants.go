package apperrors

var Areas = map[int]ErrorArea{
	1: ErrorArea("INFRA"),
	2: ErrorArea("WEB"),
	3: ErrorArea("SERVICE"),
	4: ErrorArea("REPOSITORY"),
}

const (
	ErrorAreaInfra = iota + 1
	ErrorAreaWeb
	ErrorAreaService
	ErrorAreaRepository
)

var (
	OSExitForDatabaseIssues             = 1
	OSExitForRepositoryIssues           = 2
	OSExitForRepositoryMigrationsIssues = 3
	OSExitForServiceIssues              = 4
	OSExitForWebServerIssues            = 4
)
