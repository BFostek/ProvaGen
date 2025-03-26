package scraper


type Scrapper interface {
	GetChallenge() (Challenge,error)
}

type TestCase[T any] struct {
	Param string
	Args  *T
}

type Challenge struct {
	Name        string
	Tests       []map[string]string
	Description *string
	InitialFile *string
	Solution    *string
}

