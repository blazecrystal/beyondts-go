package properties

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// a type of properties, whose format is key=value
type Properties struct {
	prop map[string]string
}

// create an empty properties
func EmptyProperties() *Properties {
	return &Properties{make(map[string]string)}
}

// create an empty properties & load key/value maps from a properties xfileutils
func LoadPropertiesFromFile(path string) (*Properties, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	p := EmptyProperties()
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return p, err
			}
		}
		str := string(line)
		indexOfEqual := strings.Index(str, "=")
		isCommont := strings.Index(str, "#") == 0
		if !isCommont && indexOfEqual > 0 {
			// is k=v format
			p.prop[strings.Trim(str[:indexOfEqual], " ")] = strings.Trim(str[indexOfEqual+1:], " ")
		}
	}
	return p, nil
}

// find key/value mappers by filtering their keys with given filter string, if nothing found, return nil, else a map
// of result key/value mappers
// filter string should be a regexp
func (p *Properties) Filter(filter string) (map[string]string, error) {
	regex, err := regexp.Compile(filter)
	if err != nil {
		return nil, err
	}
	var rst map[string]string
	for k, v := range p.prop {
		if regex.MatchString(k) {
			if rst == nil {
				rst = make(map[string]string, len(p.prop))
			}
			rst[k] = v
		}
	}
	return rst, nil
}

// get the value of key, if not exists return "", if properties itself is nil return "" also.
func (p *Properties) Get(key string) string {
	if p.prop == nil {
		return ""
	}
	return p.prop[key]
}

// get the value of key as string type, if porp is nil or key not exists, defaultValue will be returned
// the rule of this fun is different with "Get(string) string"
func (p *Properties) GetString(key, defaultValue string) string {
	if p.prop == nil {
		return defaultValue
	}
	tmp, exist := p.prop[key]
	if !exist {
		return defaultValue
	}
	return tmp
}

// get the value of key as integer type, if not exist or not a number, defaultValue will be returned
func (p *Properties) GetInt(key string, defaultValue int) int {
	tmp, err := strconv.Atoi(p.Get(key))
	if err != nil {
		return defaultValue
	}
	return tmp
}

// get the value of key as bool type, if not exist or not a value can be changed to bool, defaultValue will be returned
// It can parse value of 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
func (p *Properties) GetBool(key string, defaultValue bool) bool {
	tmp, err := strconv.ParseBool(p.Get(key))
	if err != nil {
		return defaultValue
	}
	return tmp
}

// set a value for key
func (p *Properties) Set(key, value string) {
	if p.prop == nil {
		p.prop = make(map[string]string)
	}
	p.prop[key] = value
}

// save current properites into a local xfileutils
func (p *Properties) SaveToFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(path)
			defer file.Close()
			if err != nil {
				return file, err
			}
		} else {
			return file, err
		}
	}
	if p.prop != nil {
		writer := bufio.NewWriter(file)
		for k, v := range p.prop {
			writer.WriteString(strings.Join([]string{k, v}, "="))
			writer.WriteString("\n")
			writer.Flush()
		}
	}
	file.Close()
	return file, nil
}

// get all keys as a string slice
func (p *Properties) Keys() []string {
	if p.prop == nil || len(p.prop) == 0 {
		return make([]string, 0, 0)
	}
	length := len(p.prop)
	keys := make([]string, length, length)
	index := 0
	for key, _ := range p.prop {
		keys[index] = key
		index++
	}
	return keys
}

// get all values as a string slice
func (p *Properties) Values() []string {
	if p.prop == nil || len(p.prop) == 0 {
		return make([]string, 0, 0)
	}
	length := len(p.prop)
	values := make([]string, length, length)
	index := 0
	for _, value := range p.prop {
		values[index] = value
		index++
	}
	return values
}

// judge if this properties contains a key
func (p *Properties) Contains(key string) bool {
	if p.prop == nil || len(p.prop) == 0 {
		return false
	}
	_, exist := p.prop[key]
	return exist
}

// delete a key/value map
func (p *Properties) Delete(key string) string {
	if p.prop == nil || len(p.prop) == 0 {
		return ""
	}
	value := p.prop[key]
	delete(p.prop, key)
	return value
}

// clear all key/value maps
func (p *Properties) Clear() {
	if p.prop != nil {
		for key, _ := range p.prop {
			delete(p.prop, key)
		}
	}
}

// count of key/value maps in the properties
func (p *Properties) Length() int {
	if p.prop == nil {
		return 0
	}
	return len(p.prop)
}
