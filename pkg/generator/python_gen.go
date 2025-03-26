package generator

import "github.com/BFostek/ProvaGen/pkg/scraper"

type PythonGenerator struct {
	scrapper scraper.Scrapper
}

func (pg *PythonGenerator) Generate(project_destination, challenge_id string) error {
	var err error
	pg.scrapper, err = scraper.NCodeInit(challenge_id)
	if err != nil {
		return err
	}
	if val, err := pg.scrapper.GetChallenge(); err == nil {
    println(*val.InitialFile)

	}
	return nil
}
