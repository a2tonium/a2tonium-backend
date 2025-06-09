package json_generator

type Attribute struct {
	TraitType string
	Value     interface{}
}

type Certificate struct {
	Name, IIN, Image string
	Attributes       []Attribute
	QuizGrades       []string
}
