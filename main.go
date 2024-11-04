package main

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

func main() {
	// Create CEL environment with map declarations
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("tier1Approvers", decls.NewMapType(decls.String, decls.String)),
			decls.NewVar("tier2Approvers", decls.NewMapType(decls.String, decls.String)),
			decls.NewVar("tier3Approvers", decls.NewMapType(decls.String, decls.String)),
		),
	)
	if err != nil {
		panic(err)
	}

	// Compile the expression
	// This checks if any value in each map is "APPROVED" and combines with AND
	expr := `
        tier1Approvers.exists(k, tier1Approvers[k] == "REJECTED") ||
        tier2Approvers.exists(k, tier2Approvers[k] == "REJECTED") ||
        tier3Approvers.exists(k, tier3Approvers[k] == "REJECTED")
            ? "REJECTED"
            : (tier1Approvers.exists(k, tier1Approvers[k] == "APPROVED") &&
               tier2Approvers.exists(k, tier2Approvers[k] == "APPROVED") &&
               tier3Approvers.exists(k, tier3Approvers[k] == "APPROVED")
                ? "APPROVED"
                : "PENDING")
	`

	ast, iss := env.Compile(expr)
	if iss.Err() != nil {
		panic(iss.Err())
	}

	// Create the program
	prg, err := env.Program(ast)
	if err != nil {
		panic(err)
	}

	// Test data
	vars := map[string]interface{}{
		"tier1Approvers": map[string]string{
			"user1": "PENDING",
			"user2": "APPROVED",
			"user3": "PENDING",
		},
		"tier2Approvers": map[string]string{
			"user4": "PENDING",
			"user5": "APPROVED",
		},
		"tier3Approvers": map[string]string{
			"user6": "APPROVED",
			"user7": "PENDING",
		},
	}

	// Evaluate
	result, _, err := prg.Eval(vars)
	if err != nil {
		panic(err)
	}

	fmt.Printf("All groups have at least one APPROVED user: %v\n", result)
}
