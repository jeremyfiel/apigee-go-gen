//  Copyright 2024 Google LLC
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package apiproxy

import (
	"fmt"
	v1 "github.com/apigee/apigee-go-gen/pkg/apigee/v1"
	"github.com/apigee/apigee-go-gen/pkg/common/resources"
	"github.com/apigee/apigee-go-gen/pkg/flags"
	"github.com/apigee/apigee-go-gen/pkg/render"
	"github.com/go-errors/errors"
	"github.com/spf13/cobra"
	"strings"
)

var cFlags = render.NewCommonFlags()
var debug = flags.NewBool(false)
var dryRun = flags.NewEnum([]string{"xml", "yaml"})
var validate = flags.NewBool(true)
var setValue = flags.NewSetAny(cFlags.Values)
var setValueStr = flags.NewSetString(cFlags.Values)
var setValueFile = flags.NewValues(cFlags.Values)
var setFile = flags.NewSetFile(cFlags.Values)
var setOAS = flags.NewSetOAS(cFlags.Values)
var setGraphQL = flags.NewSetGraphQL(cFlags.Values)
var setGRPC = flags.NewSetGRPC(cFlags.Values)
var setJSON = flags.NewSetJSON(cFlags.Values)

var Cmd = &cobra.Command{
	Use:   "apiproxy",
	Short: "Generate an API proxy bundle from a template",
	Long:  Usage(),
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(string(cFlags.OutputFile)) == "" && dryRun.IsUnset() && bool(debug) == false {
			return errors.New("required flag(s) \"output\" not set")
		}

		createModelFunc := func(input string) (v1.Model, error) {
			return v1.NewAPIProxyModel(input)
		}

		return render.GenerateBundle(createModelFunc, cFlags, bool(validate), dryRun.Value, bool(debug))
	},
}

func init() {
	Cmd.Flags().SortFlags = false
	Cmd.Flags().VarP(&cFlags.TemplateFile, "template", "t", `path to main template"`)
	Cmd.Flags().VarP(&cFlags.IncludeList, "include", "i", `path to helper templates (globs allowed)`)
	Cmd.Flags().VarP(&cFlags.OutputFile, "output", "o", `output directory or file`)
	Cmd.Flags().VarP(&debug, "debug", "", `prints rendered template before transforming into API proxy"`)
	Cmd.Flags().VarP(&dryRun, "dry-run", "d", `prints rendered template after transforming into API Proxy"`)
	Cmd.Flags().VarP(&validate, "validate", "v", "check for unknown elements")
	Cmd.Flags().Var(&setValue, "set", `sets a key=value (bool,float,string), e.g. "use_ssl=true"`)
	Cmd.Flags().Var(&setValueStr, "set-string", `sets key=value (string), e.g. "base_path=/v1/hello" `)
	Cmd.Flags().Var(&setValueFile, "values", `sets keys/values from YAML file, e.g. "./values.yaml"`)
	Cmd.Flags().Var(&setFile, "set-file", `sets key=value where value is the content of a file, e.g. "my_data=./from/file.txt"`)
	Cmd.Flags().Var(&setOAS, "set-oas", `sets key=value where value is an OpenAPI spec, e.g. "my_spec=./petstore.yaml"`)
	Cmd.Flags().Var(&setGRPC, "set-grpc", `sets key=value where value is a gRPC proto, e.g. "my_proto=./greeter.proto"`)
	Cmd.Flags().Var(&setGraphQL, "set-graphql", `sets key=value where value is a GraphQL schema, e.g. "my_schema=./resorts.graphql"`)
	Cmd.Flags().Var(&setJSON, "set-json", `sets key=value where value is JSON, e.g. 'servers=["server1","server2"]'`)

	_ = Cmd.MarkFlagRequired("template")
}

func Usage() string {
	usageText := `
This command takes template, renders it, and finally packages the result into an API proxy bundle.

The rendering context includes the following data:

%[1]s

Helper functions:

%[2]s

`
	helpersText, err := resources.FS.ReadFile("helper_functions.txt")
	if err != nil {
		panic(err)
	}

	renderContextText, err := resources.FS.ReadFile("render_context.txt")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(usageText, renderContextText, helpersText)
}
