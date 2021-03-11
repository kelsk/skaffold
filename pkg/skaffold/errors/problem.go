/*
Copyright 2020 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package errors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/GoogleContainerTools/skaffold/proto/v1"
)

type Description func(error) string
type Suggestion func(Config) []*proto.Suggestion

type problem struct {
	regexp      *regexp.Regexp
	description func(error) string
	errCode     proto.StatusCode
	cfg         Config
	suggestion  func(Config) []*proto.Suggestion
	err         error
}

func NewProblem(d Description, sc proto.StatusCode, config Config, s Suggestion, err error) problem {
	return problem{
		description: d,
		errCode:     sc,
		suggestion:  s,
		err:         err,
		cfg:         config,
	}
}

func (p problem) Error() string {
	description := fmt.Sprintf("%s\n", p.err)
	if p.description != nil {
		description = p.description(p.err)
	}
	if suggestions := p.suggestion(p.cfg); suggestions != nil {
		return fmt.Sprintf("%s. %s", strings.Trim(description, "."), concatSuggestions(suggestions))
	}
	return description
}

func (p problem) withConfigAndErr(cfg Config, err error) problem {
	p.err = err
	p.cfg = cfg
	return p
}

func isProblem(err error) (problem, bool) {
	if p, ok := err.(problem); ok {
		return p, true
	}
	return problem{}, false
}
