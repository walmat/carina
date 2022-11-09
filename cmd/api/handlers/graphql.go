package handlers

import (
	"context"
	"github.com/graphql-go/graphql"
	"github.com/kataras/iris/v12"
	gql2 "nebula/cmd/api/gql"
)

// define schema, with our RootQuery and RootMutation
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    gql2.RootQuery,
	Mutation: gql2.RootMutation,
})

type postData struct {
	Query     string                  `json:"query"`
	Operation *string                 `json:"operation"`
	Variables *map[string]interface{} `json:"variables"`
}

func GraphQLHandler(ctx iris.Context) {
	var data postData
	if err := ctx.ReadJSON(&data); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	var gqlOperation string
	if data.Operation != nil {
		gqlOperation = *data.Operation
	}

	var gqlVariables map[string]interface{}
	if data.Variables != nil {
		gqlVariables = *data.Variables
	}

	result := graphql.Do(graphql.Params{
		Context:        context.TODO(),
		Schema:         schema,
		RequestString:  data.Query,
		VariableValues: gqlVariables,
		OperationName:  gqlOperation,
	})

	_, err := ctx.JSON(result)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
}

var graphiql = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.css"/>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/2.0.3/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/16.2.0/umd/react.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react-dom/16.2.0/umd/react-dom.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function fetchGQL(params) {
				return fetch("/api/graphql", {
					method: "post",
					body: JSON.stringify(params),
					credentials: "include",
					headers: {
						'Content-Type': 'application/json',
					},
				}).then(function (resp) {
					return resp.text();
				}).then(function (body) {
					try {
						return JSON.parse(body);
					} catch (error) {
						return body;
					}
				});
			}
			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: fetchGQL}),
				document.getElementById("graphiql")
			)
		</script>
	</body>
</html>
`)

func GraphiQLHandler(ctx iris.Context) {
	ctx.ContentType("text/html")
	_, err := ctx.Write(graphiql)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
}
