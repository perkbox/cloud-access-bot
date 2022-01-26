package messenger

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"text/template"

	"github.com/slack-go/slack"

	"github.com/sirupsen/logrus"
)

func counter() func() int {
	i := -1
	return func() int {
		i++
		return i
	}
}

func GetValuesFromSelectedOptions(selectedOptions []slack.OptionBlockObject) []string {
	var optVals []string
	for _, v := range selectedOptions {
		optVals = append(optVals, v.Value)
	}
	return optVals
}

// SliceToOptions Given a slice will return Options required for a multiselect and other lists
func SliceToOptions(slice []string, TextType string) Options {
	var opts []Option
	if len(slice) > 100 {
		logrus.Errorf("Slice Greater than 100 limit of modal Options: %d", len(slice))
		return Options{}
	}

	for _, service := range slice {
		opts = append(opts, Option{
			Text: Text{
				Text: service,
				Type: TextType,
			},
			Value: service,
		})
	}

	return Options{Options: opts}
}

// MapToOptions Given a map[string]string will return Options required for a multiselect and other lists
// Value will be shown within the Text of the option
// Key will be shown in the Text and as the actual value of the option
// e.g Key (Value)
func MapToOptions(slice map[string]string, TextType string) Options {
	var opts []Option
	for service, id := range slice {
		opts = append(opts, Option{
			Text: Text{
				Text: service,
				Type: TextType,
			},
			Value: id,
		})
	}

	options := Options{Options: opts}
	return options
}

func renderTemplate(fs fs.FS, file string, args interface{}) (bytes.Buffer, error) {
	var tpl bytes.Buffer

	f, err := fs.Open(file)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("error opening asset file %s", err.Error())
	}

	fby, _ := ioutil.ReadAll(f)

	//ParseFS doesn't seem to work
	t := template.Must(template.New("tmpl").Funcs(template.FuncMap{"counter": counter}).Parse(string(fby)))

	err = t.Execute(&tpl, args)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("error Rendering template %s", err.Error())
	}

	return tpl, nil
}
