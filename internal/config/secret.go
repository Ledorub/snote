package config

type secretString struct {
	value string
}

func (s secretString) String() string {
	if s.value == "" {
		return ""
	}
	return "***"
}

func (s secretString) GoString() string {
	return s.String()
}

func (s secretString) GetValue() string {
	return s.value
}

func (s secretString) MarshalYAML() (any, error) {
	return s.String(), nil
}

func (s *secretString) UnmarshalYAML(unmarshal func(any) error) error {
	return unmarshal(&s.value)
}
