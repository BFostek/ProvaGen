package scraper

type Scrapper interface {
	GetChallenge() Challenge
}

type TestCase[T any] struct{
  Param string
  Args *T
}

type Challenge struct {
	Name string
  Tests map[string]*any

  Description *string
  InitialFile *string
}
