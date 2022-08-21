package govespa

import (
	"errors"

	"github.com/mitchellh/mapstructure"
)

// TODO: iter should be able to reexectue any request with a new continuation token.
// The iterator should be unique for every Get/Query/Remove/Update instance,
// it Should be able to rexecute it with a different continuation token.
// That way we can implement "Get Visit", "Remove Where", "Update where"
type iter struct {
	res []map[string]any
}

// Get scans the first result into a destination.
// The destination needs to be a pointer to a struct which fields are annotated with the "vespa" Tag.
func (i *iter) Get(dest any) error {
	if i.res == nil || len(i.res) == 0 {
		return errors.New("Fields are empty")
	}
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  dest,
		TagName: "vespa",
	})
	if err != nil {
		return err
	}
	d.Decode(i.res[0])
	return nil
}

// Select scans all results into a destination, which must be a pointer to a slice.
func (i *iter) Select(dest any) {

}
