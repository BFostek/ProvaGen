package generator

import (
	"github.com/BFostek/ProvaGen/pkg/scraper"
	"github.com/BFostek/ProvaGen/pkg/utils"
)

type PythonGenerator struct {
	scrapper scraper.Scrapper
}

func (pg *PythonGenerator) Generate(project_destination, challenge_id string) error {
	var err error
	pg.scrapper, err = scraper.NCodeInit(challenge_id)
	if err != nil {
		return err
	}
	var val scraper.Challenge
	if val, err = pg.scrapper.GetChallenge(); err != nil {
		return err
	}
  utils.CreateStructure(project_destination,challenge_id)
	println(val.Slug)
  

	return nil
}
