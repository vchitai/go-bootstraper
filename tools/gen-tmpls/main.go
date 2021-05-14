package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), os.ModePerm); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func main() {
	err := os.RemoveAll("./assets/templates")
	if err != nil {
		log.Fatal(err)
	}

	baskets := map[string][]string{
		"base": {
			"cmd/server/main.go",
			"configs/configs.go",
			"cmd/server/service_gen.go",
			"internal/services/base.go",
			".gitignore",
			".gitlab-ci.yml",
			".golangci.yml",
			"Dockerfile",
			"Makefile",
			"tools.go",
			"README.md",
		},
		//"configs": {
		//	"configs/config_gen.go",
		//},
		//"mapping": {
		//	"internal/mapping/wish.go",
		//	"internal/mapping/wish_gen.go",
		//},
		//"models": {
		//	"internal/models/model_name_lower.go",
		//},
		//"cmd": {
		//	"cmd/server/server_name_lower.go",
		//},
		"services": {
			"internal/services/server_name_lower.go",
			"internal/services/server_name_lower_gen.go",
		},
		//"stores": {
		//	"internal/stores/store_name_lower.go",
		//	"internal/stores/store_name_lower_gen.go",
		//},
		"proto": {
			"proto/proto_name_lower.proto",
			"proto/proto_name_lower_message.proto",
			"pkg/client/proto_name_lower.go",
		},
	}

	basketIndex := make(map[string]string)
	for l, b := range baskets {
		for _, v := range b {
			basketIndex[v] = l
		}
	}

	ignored := []string{
		"vendor",
		"docs",
		"pb",
	}

	ifRegex := regexp.MustCompile("\\/\\* if (.*) \\*\\/")
	rangeRegex := regexp.MustCompile("\\/\\* range (.*) \\*\\/")
	withRegex := regexp.MustCompile("\\/\\* with (.*) \\*\\/")

	os.Chdir("assets")

	params := []string{
		"ServiceTitle",
		"serviceCamel",
		"service_lower",
		"service_dash",
		"DomainTitle",
		"domainCamel",
		"domain_lower",
		"domain_dash",
		"ProtoNameTitle",
		"protoNameCamel",
		"proto_name_lower",
		"proto_name_dash",
		"ServerNameTitle",
		"serverNameCamel",
		"server_name_lower",
		"server_name_dash",
		"StoreNameTitle",
		"storeNameCamel",
		"store_name_lower",
		"store_name_dash",
		"ModelNameTitle",
		"modelNameCamel",
		"model_name_lower",
		"model_name_dash",
		"DBNameTitle",
		"dbNameCamel",
		"db_name_lower",
		"db_name_dash",
	}
	if err := filepath.Walk("plasma", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		for _, ignore := range ignored {
			if strings.Contains(path, ignore) {
				return nil
			}
		}
		dir, file := filepath.Split(path)
		nestedDir := strings.Split(dir, string(os.PathSeparator))
		if len(nestedDir) == 0 || nestedDir[0] != "plasma" {
			return fmt.Errorf("not right")
		}

		x := filepath.Join(filepath.Join(nestedDir[1:]...), file)
		basket, ok := basketIndex[x]
		if !ok {
			log.Printf("not basket found for %s\n", x)
			return nil
		}
		dir = filepath.Join("templates", basket, filepath.Join(nestedDir[1:]...))

		if strings.Contains(dir, "plasma") {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		s := string(b)

		for _, p := range params {
			s = strings.ReplaceAll(s, p, fmt.Sprintf("{{.%s}}", p))
		}
		xx := ifRegex.FindAllStringSubmatch(s, -1)
		for _, x := range xx {
			s = strings.ReplaceAll(s, x[0], fmt.Sprintf("{{ if %s }}", x[1]))
		}
		xxx := rangeRegex.FindAllStringSubmatch(s, -1)
		for _, x := range xxx {
			s = strings.ReplaceAll(s, x[0], fmt.Sprintf("{{ range %s }}", x[1]))
		}
		xxxx := withRegex.FindAllStringSubmatch(s, -1)
		for _, x := range xxxx {
			s = strings.ReplaceAll(s, x[0], fmt.Sprintf("{{ with %s }}", x[1]))
		}
		s = strings.ReplaceAll(s, "/* end */", "{{ end }}")

		//s = strings.ReplaceAll(s, "Base", "{{.ServiceTitle}}")
		//if s != oldS {
		//	templateGenerated = true
		//}

		file += ".tmpl"
		if file[0] == '.' {
			file = "dot@" + file[1:]
		}
		p := filepath.Join(dir, file)
		fo, err := create(p)
		if err != nil {
			return err
		}
		defer fo.Close()
		err = ioutil.WriteFile(p, []byte(s), 0644)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
