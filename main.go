package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/oauth2"
)

var (
	githubLink   = "https://api.github.com/graphql"
	exampleQuery = `query {
		repository(owner:"isacikgoz", name:"gitin") {
		  issues(last:20, states:CLOSED) {
			edges {
			  node {
				title
			  }
			}
		  }
		}
	  }`
)

type query struct {
	Query string `json:"query"`
}

type response struct {
	Data *json.RawMessage
}

func main() {
	token := os.Getenv("GITHUB_ACCESS_TOKEN")
	src := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)
	client := oauth2.NewClient(context.Background(), src)

	ctx := context.Background()
	var buf bytes.Buffer
	q := query{Query: exampleQuery}
	err := json.NewEncoder(&buf).Encode(q)
	if err != nil {
		fmt.Printf("could not encode: %v\n", err)
	}

	res, err := ctxhttp.Post(ctx, client, githubLink, "application/json", &buf)
	if err != nil {
		fmt.Printf("could not encode: %v\n", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("non-200 OK status code: %v body: %q\n", res.Status, body)
		os.Exit(2)
	}
	var out response
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		fmt.Printf("could not encode: %v\n", err)
	}
	fmt.Printf("response body: %s\n", out.Data)
}
