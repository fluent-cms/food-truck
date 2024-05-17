package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

func ParseYaml[T any](ps ...string) (*T, error) {
	obj := new(T)
	success := false

	var all []string
	for _, p := range ps {
		all = append(all, p)
		ext := path.Ext(p)
		all = append(all, p[:len(p)-len(ext)]+".dev"+ext)
		all = append(all, p[:len(p)-len(ext)]+".prod"+ext)
	}

	for _, p := range all {
		if err := parse(obj, p); err == nil {
			fmt.Println("Load Config from ", p)
			success = true
		}
	}
	if !success {
		return obj, fmt.Errorf("fail to load config from all these path :%v", all)
	}
	return obj, nil
}

func parse[T any](t *T, filePath string) (err error) {
	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		return
	}
	defer func(f *os.File) {
		err = f.Close()
	}(f)

	if err = yaml.NewDecoder(f).Decode(t); err != nil {
		return err
	}
	return nil
}
