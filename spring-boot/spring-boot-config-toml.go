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
	"fmt"
	"github.com/pelletier/go-toml"
)

type ConfigParserToml struct {
}

func (_ *ConfigParserToml) FileExt() []string {
	return []string{".toml"}
}

func (parser *ConfigParserToml) Parse(ctx SpringApplicationContext, path string) error {

	tree, err := toml.LoadFile(path)
	if err != nil {
		return err
	}

	return parser.Tree2FlatMap(ctx, tree, "")
}

func (parser *ConfigParserToml) Tree2FlatMap(ctx SpringApplicationContext, tree *toml.Tree, key string) error {

	keys := tree.Keys()
	for i := range keys {
		k := keys[i]

		subKey := key
		if subKey == "" {
			subKey += k
		} else {
			subKey += "." + k
		}

		v := tree.Get(k)
		switch node := v.(type) {
		case *toml.Tree:
			if err := parser.Tree2FlatMap(ctx, node, subKey); err != nil {
				return err
			}
		default:
			ctx.SetProperties(subKey, fmt.Sprint(node))
		}
	}
	return nil
}
