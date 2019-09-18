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
	Logger "github.com/didi/go-spring/spring-logger"
	"strings"
	"os"
	"errors"
	"bufio"
	"io"
)

type ConfigParserProperties struct {
}

func (_ *ConfigParserProperties) FileExt() []string {
	return []string{".properties"}
}

func (parser *ConfigParserProperties) Parse(ctx SpringApplicationContext, path string) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	Logger.Infoln("load properties from ", path)

	parse := func(s string) (err error) {
		if s[0] != '#' { // 忽略配置文件中的注释行
			if ss := strings.Split(s, "="); len(ss) != 2 {
				err = errors.New("error content: " + s)
			} else {
				k := strings.TrimSpace(ss[0])

				v := strings.TrimSpace(ss[1])
				v = strings.Trim(v, "\"")

				ctx.SetProperties(k, v)
			}
		}
		return
	}

	reader := bufio.NewReader(file)
	for {

		str, err := reader.ReadString('\n')
		if err == nil || err == io.EOF {

			str = strings.TrimSpace(str)
			if len(str) > 0 {

				err = parse(str)
				if err != nil {
					return err
				}
			}

			if err == io.EOF {
				return nil
			}

		} else {
			return err
		}
	}
}
