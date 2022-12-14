package govespa

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type scanner struct {
	res []map[string]any
}

// Get scans the first result into a destination.
// The destination needs to be a pointer to a struct which fields are annotated with the "vespa" Tag.
func (s *scanner) Get(dest any) error {
	value := reflect.ValueOf(dest)

	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("expected a pointer but got %T", dest)
	}
	if value.IsNil() {
		return errors.New("expected a pointer but got nil")
	}

	if s.res == nil || len(s.res) == 0 {
		return errors.New("Fields are empty")
	}
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  dest,
		TagName: "vespa",
	})
	if err != nil {
		return err
	}
	err = d.Decode(s.res[0])
	if err != nil {
		return err
	}
	return nil
}

// Select scans all results into a destination, which must be a pointer to a slice.
func (s *scanner) Select(dest any) error {
	value := reflect.ValueOf(dest)

	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("expected a pointer but got %T", dest)
	}
	if value.IsNil() {
		return errors.New("expected a pointer but got nil")
	}

	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  dest,
		TagName: "vespa",
	})
	if err != nil {
		return err
	}

	err = d.Decode(s.res)
	if err != nil {
		return err
	}
	return nil
}
