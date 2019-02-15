package semver

import (
	s "github.com/blang/semver"
)

func Parse(str string) *s.Version {
	v, err := s.Parse(str)
	if err != nil {
		return nil
	}

	return &v
}
