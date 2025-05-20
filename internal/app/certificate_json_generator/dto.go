package certificate_json_generator

type Attribute struct {
	TraitType string
	Value     interface{}
}

type Certificate struct {
	Name        string
	Description string
	Image       string
	Attributes  []Attribute
	QuizGrades  []string
}
