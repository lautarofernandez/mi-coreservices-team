package babel

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/leonelquinteros/gotext"
	"golang.org/x/text/language"
)

const (
	pluralForm = `"Plural-Forms: nplurals=2; plural=(n != 1);\n"` + "\n\n"
)

var translations = map[string]*gotext.Po{}

func Tr(locale language.Tag, key string, args ...interface{}) string {
	translation, ok := translations[locale.String()]
	if !ok {
		return fmt.Sprintf(key, args...)
	}
	return translation.Get(key, args...)
}

func Trn(locale language.Tag, count int, singular string, plural string, args ...interface{}) string {
	translation, ok := translations[locale.String()]
	if !ok {
		key := singular
		if count > 1 {
			key = plural
		}
		return fmt.Sprintf(key, args...)
	}
	return translation.GetN(singular, plural, count, args...)
}

func Load(bundle string) error {
	reader, err := zip.OpenReader(bundle)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, entry := range reader.File {
		if !entry.FileInfo().IsDir() && filepath.Ext(entry.Name) == ".po" {

			err := func() error {
				file, err := entry.Open()
				if err != nil {
					return err
				}
				defer file.Close()
				return addTranslation(path.Dir(entry.Name), file)
			}()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func LoadDir(bundleDir string) error {
	err := filepath.Walk(bundleDir, func(fullpath string, f os.FileInfo, err error) error {
		if !f.IsDir() && filepath.Ext(fullpath) == ".po" {
			return func() error {
				fd, err := os.Open(fullpath)
				if err != nil {
					return err
				}
				defer fd.Close()
				// determine locale
				dir, _ := path.Split(fullpath)
				locale := path.Base(dir)
				return addTranslation(locale, fd)

			}()
		}
		return nil
	})
	return err
}

func addTranslation(locale string, reader io.Reader) error {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	po := new(gotext.Po)
	po.Parse(pluralForm + string(content))
	translations[locale] = po
	return nil
}
