package namespace

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Namespace struct {
	Dir         string `yaml:"dir"`
	Httpauth    string `yaml:"httpauth"`
	AllowUpload bool   `yaml:"allowUpload"`
	AllowDelete bool   `yaml:"allowDelete"`
}

func (ns Namespace) String() string {
	bytes, _ := json.Marshal(&ns)
	return string(bytes)
}

var namespaces = make(map[string]*Namespace)

func init() {
	file, err := os.Open("./namespace/namespace.yml")
	if err != nil {
		log.Println("namespace err:", err.Error())
	}
	ymlData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("namespace err:", err.Error())
	}

	var _namespaces = []*Namespace{}
	err = yaml.Unmarshal(ymlData, &_namespaces)
	if err != nil {
		log.Println("namespace err:", err.Error())
	}

	for _, ns := range _namespaces {
		arr := strings.Split(ns.Httpauth, ":")
		if len(arr) != 2 {
			continue
		}
		namespaces[arr[0]] = ns
	}

	log.Println("namespaces =>", namespaces)

}

func BasicAuth() map[string][]string {
	m := make(map[string][]string)

	for _, ns := range namespaces {
		userpass := strings.SplitN(ns.Httpauth, ":", 2)
		if len(userpass) == 2 {
			user, pass := userpass[0], userpass[1]

			if m[user] != nil {
				m[user] = append(m[user], pass)
			} else {
				m[user] = []string{pass}
			}

		}
	}
	return m
}

// // return: allow_access, allow_upload, allow_delete
// func Check(dir, user string) (allow_access bool, allow_upload bool, allow_delete bool) {
// 	if namespaces == nil {
// 		return
// 	}

// 	namespace := namespaces[dir]
// 	if namespace == nil {
// 		return
// 	}

// 	if strings.HasPrefix(namespace.Httpauth, user+":") {
// 		allow_access = true
// 		allow_upload = namespace.AllowUpload
// 		allow_delete = namespace.AllowDelete
// 	}

// 	return
// }

func Get(user string) *Namespace {
	if namespaces == nil {
		return nil
	}

	return namespaces[user]
}
