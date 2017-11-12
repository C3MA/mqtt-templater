package main

import (
	"io/ioutil"
	"regexp"
	"bytes"
)

type Template struct {
	Template []byte
	Data map[string][]byte
}

func NewTemplate() *Template {
	t := new(Template)
	t.Data = make(map[string][]byte)

	return t
}

func (t *Template) ReadTemplate(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`\$\(.*\)`)
	keys := re.FindAllString(string(buf), -1)
	for _, key := range keys {
		t.Data[key] = nil
	}

	t.Template = buf

	return nil
}

func (t *Template) WriteOutput(filename string) error {
	output := make([]byte, len(t.Template))
	copy(output, t.Template)

	for key, value := range t.Data {
		if value == nil {
			value = []byte("null")
		}
		output = bytes.Replace(output, []byte(key), value, -1)
	}

	err := ioutil.WriteFile(filename, output, 0666)
	return err
}

func (t *Template) SetVariable(key string, value []byte) bool {
	key = "$(" + key + ")"

	if _, hasKey := t.Data[key]; hasKey == false {
		return false
	}

	mapValue := t.Data[key]
	if bytes.Compare(mapValue, value) == 0 {
		return false
	}

	t.Data[key] = value
	return true
}
