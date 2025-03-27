package generator

type TestGenerator interface {
	Generate(folder_path, challenge_id string) error
}

func GeneratorFactory(param... string) TestGenerator {
	if len(param)==0{
    result := PythonGenerator{}
    return &result
	}
	return nil

}
