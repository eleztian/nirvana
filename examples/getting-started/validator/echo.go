/*
Copyright 2018 Caicloud Authors

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

package main

import (
	"context"

	"github.com/caicloud/nirvana/config"
	"github.com/caicloud/nirvana/definition"
	"github.com/caicloud/nirvana/errors"
	"github.com/caicloud/nirvana/log"
	"github.com/caicloud/nirvana/operators/validator"
)

type Body struct {
	Name string
}

var echo = definition.Descriptor{
	Path:        "/echo",
	Description: "Echo API",
	Definitions: []definition.Definition{
		{
			Method:   definition.Create,
			Function: Echo,
			Consumes: []string{definition.MIMEJSON},
			Produces: []string{definition.MIMEJSON},
			Parameters: []definition.Parameter{
				{
					Source:      definition.Query,
					Name:        "msg",
					Description: "Corresponding to the second parameter",
					Operators:   []definition.Operator{validator.String("gt=10")},
				},
				{
					Source:      definition.Body,
					Name:        "body",
					Description: "How to do custom validation",
					Operators: []definition.Operator{
						validator.NewCustom(
							func(ctx context.Context, body *Body) error {
								if body.Name == "" {
									return errors.BadRequest.Error("you should have a name!")
								}
								if body.Name != "nirvana" {
									return errors.BadRequest.Error("name ${name} must be nirvana!", body.Name)
								}
								return nil
							},
							"validate your name"),
					},
				},
			},
			Results: []definition.Result{
				{
					Destination: definition.Data,
					Description: "Corresponding to the first result",
				},
				{
					Destination: definition.Error,
					Description: "Corresponding to the second result",
				},
			},
		},
	},
}

// API function.
func Echo(ctx context.Context, msg string, body *Body) (string, error) {
	return msg, nil
}

func main() {
	cmd := config.NewDefaultNirvanaCommand()
	if err := cmd.Execute(echo); err != nil {
		log.Fatal(err)
	}
}
