package criteria

import (
	"github.com/open-policy-agent/opa/ast"

	"github.com/pomerium/pomerium/pkg/policy/parser"
)

var corsPreflightBody = ast.Body{
	ast.MustParseExpr(`input.http.method == "OPTIONS"`),
	ast.MustParseExpr(`count(object.get(input.http.headers, "Access-Control-Request-Method", [])) > 0`),
	ast.MustParseExpr(`count(object.get(input.http.headers, "Origin", [])) > 0`),
}

type corsPreflightCriterion struct {
	g *Generator
}

func (corsPreflightCriterion) Names() []string {
	return []string{"cors_preflight"}
}

func (c corsPreflightCriterion) GenerateRule(_ string, _ parser.Value) (*ast.Rule, []*ast.Rule, error) {
	rule := c.g.NewRule("cors_preflight")
	rule.Body = corsPreflightBody
	return rule, nil, nil
}

// CORSPreflight returns a Criterion which returns true if the input request is a CORS preflight request.
func CORSPreflight(generator *Generator) Criterion {
	return corsPreflightCriterion{g: generator}
}

func init() {
	Register(CORSPreflight)
}
