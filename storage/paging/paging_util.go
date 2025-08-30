/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package paging

/*
CountPages counts the number of pages of a certain type of a given PagedStorageFile.
*/
func CountPages(pager *PagedStorageFile, pagetype int16) (int, error) {

	var err error

	cursor := NewPageCursor(pager, pagetype, 0)

	page, _ := cursor.Next()
	counter := 0

	for page != 0 {
		counter++

		page, err = cursor.Next()
		if err != nil {
			return -1, err
		}
	}

	return counter, nil
}
