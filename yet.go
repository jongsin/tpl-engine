package yet

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

type Yet struct {
	Tepmlates []*template.Template
}

const (
	Version = "0.0.1"
	website = "https://yet.yoonit.id"
	banner  = `
 __  __  ____  ______
 \ \/ / / __/ /_   _/
  \  / / _/    /  /
  /_/ /___/   /__/  v.%s
Easy use, fast development, Yoon Easy Template Engine
%s
`
	PathSeparator = string(os.PathSeparator)
)

func New() (y *Yet) {
	// color.Printf(banner, color.Red(Version), color.Green(website))

	return &Yet{
		Tepmlates: make([]*template.Template, 0),
	}
}

func (y *Yet) addTpl(t *template.Template) error {
	for _, tpl := range y.Tepmlates {
		if tpl.Name() == t.Name() {
			return fmt.Errorf("%s already exists", t.Name())
		}
	}
	y.Tepmlates = append(y.Tepmlates, t)
	return nil
}

func (y *Yet) GetTemplate(tName string) (*template.Template, error) {
	for _, tpl := range y.Tepmlates {
		if tpl.Name() == tName {
			return tpl, nil
		}
	}
	return nil, fmt.Errorf("%s not found", tName)
}

func (y *Yet) MailMerge(tName string, params map[string]interface{}) (string, error) {
	t, err := y.GetTemplate(tName)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	// mw := io.MultiWriter(os.Stdout, &b)
	mw := io.MultiWriter(&b)

	err = t.Execute(mw, params)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func (y *Yet) LoadAndCache(path string) error {
	tpls, err := y.LoadTemplates(path)
	if err != nil {
		return err
	}
	for _, tpl := range tpls {
		err = y.addTpl(tpl)
		if err != nil {
			return err
		}
	}
	return nil
}

func (y *Yet) LoadTemplate(path string) (*template.Template, error) {
	tName := path
	if strings.Contains(path, PathSeparator) {
		tName = path[strings.LastIndex(path, PathSeparator):]
	}

	tmp, err := template.New(tName).ParseFiles(path)
	if err != nil {
		return nil, err
	}
	return tmp, nil
}

func (y *Yet) LoadTemplates(path string) ([]*template.Template, error) {
	tplFiles, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	res := make([]*template.Template, 0)

	for _, tpl := range tplFiles {
		if tpl.IsDir() {
			continue
		}

		fPath := path + PathSeparator + tpl.Name()
		tFile, err := y.LoadTemplate(fPath)
		if err != nil {
			return nil, err
		}
		res = append(res, tFile)
	}

	return res, nil
}
