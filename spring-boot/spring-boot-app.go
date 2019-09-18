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
	"os"
	"fmt"
)

type SpringBootApplication struct {
	AppContext     SpringApplicationContext
	ConfigLocation string
	ConfigParsers  []ConfigParser
}

func NewSpringBootApplication(configLocation string) *SpringBootApplication {
	return &SpringBootApplication{
		AppContext:     NewSpringApplicationContext(),
		ConfigLocation: configLocation,
		ConfigParsers: []ConfigParser{
			new(ConfigParserProperties),
			new(ConfigParserYaml),
			new(ConfigParserToml),
		},
	}
}

//
// 启动 Spring 服务器
//
func (app *SpringBootApplication) Start() {

	if err := app.loadConfigFiles(); err != nil {
		panic(err)
	}

	app.AppContext.RegisterBean(app.AppContext)

	// 初始化各个模块
	for _, fn := range ModuleFuncs {
		fn(app.AppContext)
	}

	// 检查 Bean 的自动绑定
	if err := app.AppContext.AutoWireBeans(); err != nil {
		panic(err)
	}

	// 通知应用启动事件
	var eventBeans []SpringApplicationEvent
	app.AppContext.FindBeansByType(&eventBeans)

	if eventBeans != nil && len(eventBeans) > 0 {
		for _, bean := range eventBeans {
			bean.OnStartApplication(app.AppContext)
		}
	}
}

func (app *SpringBootApplication) loadConfigFiles0(filePath string) error {
	for _, parser := range app.ConfigParsers {
		for _, ext := range parser.FileExt() {
			err := parser.Parse(app.AppContext, filePath+ext)
			if err != nil {
				if _, ok := err.(*os.PathError); !ok {
					return err
				}
			}
		}
	}
	return nil
}

func (app *SpringBootApplication) loadConfigFiles() error {

	// 加载默认的应用配置文件

	filePath := app.ConfigLocation + "application"
	if err := app.loadConfigFiles0(filePath); err != nil {
		return err
	}

	// 加载用户设置的配置文件

	if env := os.Getenv("spring.profile"); len(env) > 0 {
		filePath = fmt.Sprintf(app.ConfigLocation+"application-%s", env)
		if err := app.loadConfigFiles0(filePath); err != nil {
			return err
		}
	}

	return nil
}

func (app *SpringBootApplication) ShutDown() {

	// 通知应用启动事件
	var eventBeans []SpringApplicationEvent
	app.AppContext.FindBeansByType(&eventBeans)

	if eventBeans != nil && len(eventBeans) > 0 {
		for _, bean := range eventBeans {
			bean.OnStopApplication(app.AppContext)
		}
	}

	app.AppContext.Wait()
	fmt.Println("spring exit")
}
