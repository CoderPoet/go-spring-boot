/*
 * Copyright 2012-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package SpringBoot

import (
	"io/ioutil"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
)

type ConfigParserYaml struct {
}

func (_ *ConfigParserYaml) FileExt() []string {
	return []string{".yaml"}
}

func (parser *ConfigParserYaml) Parse(ctx SpringApplicationContext, path string) error {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var config map[interface{}]interface{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	return parser.Struct2FlatMap(ctx, config, "")
}

func (parser *ConfigParserYaml) Struct2FlatMap(ctx SpringApplicationContext, st map[interface{}]interface{}, key string) error {

	for k := range st {
		subKey := key

		if str, ok := k.(string); ok {
			if subKey == "" {
				subKey += str
			} else {
				subKey += "." + str
			}
		} else {
			panic(errors.New(""))
		}

		switch val := st[k].(type) {
		case map[interface{}]interface{}:
			if err := parser.Struct2FlatMap(ctx, val, subKey); err != nil {
				return err
			}
		default:
			ctx.SetProperties(subKey, fmt.Sprint(val))
		}
	}
	return nil
}
