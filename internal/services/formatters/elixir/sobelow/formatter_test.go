// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sobelow

import (
	"errors"
	"testing"

	"github.com/ZupIT/horusec-devkit/pkg/entities/analysis"
	"github.com/ZupIT/horusec-devkit/pkg/enums/tools"
	"github.com/stretchr/testify/assert"

	cliConfig "github.com/ZupIT/horusec/config"
	"github.com/ZupIT/horusec/internal/entities/toolsconfig"
	"github.com/ZupIT/horusec/internal/entities/workdir"
	"github.com/ZupIT/horusec/internal/services/formatters"
	"github.com/ZupIT/horusec/internal/utils/testutil"
)

func getOutputString() string {
	return `


		[31m[+][0m Config.CSP: Missing Content-Security-Policy - lib/built_with_elixir_web/router.ex:9
		[31m[+][0m Config.Secrets: Hardcoded Secret - config/travis.exs:24
		[31m[+][0m Config.HTTPS: HTTPS Not Enabled - config/prod.exs:0
		[32m[+][0m XSS.Raw: XSS - lib/built_with_elixir_web/templates/layout/app.html.eex:17
		test

`
}

func TestStartCFlawfinder(t *testing.T) {
	t.Run("should success execute container and process output", func(t *testing.T) {
		dockerAPIControllerMock := testutil.NewDockerMock()
		entity := &analysis.Analysis{}
		config := &cliConfig.Config{}
		config.WorkDir = &workdir.WorkDir{}

		output := getOutputString()

		dockerAPIControllerMock.On("CreateLanguageAnalysisContainer").Return(output, nil)

		service := formatters.NewFormatterService(entity, dockerAPIControllerMock, config)
		formatter := NewFormatter(service)

		formatter.StartAnalysis("")

		assert.NotEmpty(t, entity)
		assert.Len(t, entity.AnalysisVulnerabilities, 4)
	})

	t.Run("should return error when invalid output", func(t *testing.T) {
		dockerAPIControllerMock := testutil.NewDockerMock()
		analysis := &analysis.Analysis{}
		config := &cliConfig.Config{}
		config.WorkDir = &workdir.WorkDir{}

		output := ""

		dockerAPIControllerMock.On("CreateLanguageAnalysisContainer").Return(output, nil)

		service := formatters.NewFormatterService(analysis, dockerAPIControllerMock, config)
		formatter := NewFormatter(service)

		assert.NotPanics(t, func() {
			formatter.StartAnalysis("")
		})
	})

	t.Run("should return error when executing container", func(t *testing.T) {
		dockerAPIControllerMock := testutil.NewDockerMock()
		entity := &analysis.Analysis{}
		config := &cliConfig.Config{}
		config.WorkDir = &workdir.WorkDir{}

		dockerAPIControllerMock.On("CreateLanguageAnalysisContainer").Return("", errors.New("test"))

		service := formatters.NewFormatterService(entity, dockerAPIControllerMock, config)
		formatter := NewFormatter(service)

		assert.NotPanics(t, func() {
			formatter.StartAnalysis("")
		})
	})

	t.Run("should not execute tool because it's ignored", func(t *testing.T) {
		entity := &analysis.Analysis{}
		dockerAPIControllerMock := testutil.NewDockerMock()
		config := &cliConfig.Config{}
		config.WorkDir = &workdir.WorkDir{}
		config.ToolsConfig = toolsconfig.ToolsConfig{
			tools.Sobelow: toolsconfig.Config{
				IsToIgnore: true,
			},
		}

		service := formatters.NewFormatterService(entity, dockerAPIControllerMock, config)
		formatter := NewFormatter(service)

		formatter.StartAnalysis("")
	})
}
