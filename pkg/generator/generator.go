package generator

type TestGenerator interface {
	Generate(folder_path, challenge_id string) error
}

func GeneratorFactory(param string) TestGenerator {
	if param == ""  || true{
    result := PythonGenerator{}
    return &result
	}
	return nil

}
