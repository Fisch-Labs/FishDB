/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package parser

/*
RuntimeProvider provides runtime components for a parse tree.
*/
type RuntimeProvider interface {

	/*
	   Runtime returns a runtime component for a given ASTNode.
	*/
	Runtime(node *ASTNode) Runtime
}

/*
Runtime provides the runtime for an ASTNode.
*/
type Runtime interface {

	/*
	   Validate this runtime component and all its child components.
	*/
	Validate() error

	/*
		Eval evaluate this runtime component.
	*/
	Eval() (interface{}, error)
}
