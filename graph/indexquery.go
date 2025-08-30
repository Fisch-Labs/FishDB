/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package graph

/*
IndexQuery models the interface to the full text search index.
*/
type IndexQuery interface {

	/*
		LookupPhrase finds all nodes where an attribute contains a certain
		phrase. This call returns a list of node keys which contain the phrase
		at least once.
	*/
	LookupPhrase(attr, phrase string) ([]string, error)

	/*
		LookupWord finds all nodes where an attribute contains a certain word.
		This call returns a map which maps node key to a list of word positions.
	*/
	LookupWord(attr, word string) (map[string][]uint64, error)

	/*
		LookupValue finds all nodes where an attribute has a certain value.
		This call returns a list of node keys.
	*/
	LookupValue(attr, value string) ([]string, error)
}
