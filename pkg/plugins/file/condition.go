package file

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/olblak/updateCli/pkg/core/scm"
)

// Condition test if a file content match the content provided via configuration.
// If the configuration doesn't specify a value then it fall back to the source output
func (f *File) Condition(source string) (bool, error) {

	// If f.Content is not provided then we use the value returned from source
	if len(f.Content) == 0 {
		f.Content = source
	}

	data, err := Read(f.File, "")
	if err != nil {
		return false, err
	}

	content := string(data)

	f.Content, err = f.Line.ContainsExcluded(f.Content)

	if err != nil {
		return false, err
	}

	if ok, err := f.Line.ContainsIncluded(f.Content); err != nil || !ok {
		if err != nil {
			return false, err
		}

		if !ok {
			return false, fmt.Errorf(ErrLineNotFound)
		}

	}

	f.Content, err = f.Line.ContainsIncludedOnly(f.Content)

	if err != nil {
		return false, err
	}

	if strings.Compare(f.Content, content) == 0 {
		logrus.Infof("\u2714 Content from file '%v' is correct'", filepath.Join(f.File))
		return true, nil
	}

	logrus.Infof("\u2717 Wrong content from file '%v'. \n%s",
		f.File, Diff(f.Content, content))

	return false, nil
}

// ConditionFromSCM test if a file content from SCM match the content provided via configuration.
// If the configuration doesn't specify a value then it fall back to the source output
func (f *File) ConditionFromSCM(source string, scm scm.Scm) (bool, error) {

	if len(f.Content) > 0 {
		logrus.Infof("Key content defined from updatecli configuration")
	} else {
		f.Content = source
	}

	data, err := Read(f.File, scm.GetDirectory())
	if err != nil {
		return false, err
	}

	if strings.Compare(f.Content, string(data)) == 0 {
		logrus.Infof("\u2714 Content from file '%v' is correct'", f.File)
		return true, nil
	}

	logrus.Infof("\u2717 Wrong content from file '%v'. \n%s",
		f.File, Diff(f.Content, string(data)))

	return false, nil
}
