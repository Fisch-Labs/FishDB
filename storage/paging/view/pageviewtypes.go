/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package view

/*
TypeFreePage is a free page waiting to be (re)allocated
*/
const TypeFreePage = 0

/*
TypeDataPage is a page which is in use and contains data
*/
const TypeDataPage = 1

/*
TypeTranslationPage is a page which translates between physical and logical row ids
*/
const TypeTranslationPage = 2

/*
TypeFreeLogicalSlotPage is a page which holds free logical slot ids
(used to give stable ids to objects which can grow in size)
*/
const TypeFreeLogicalSlotPage = 3

/*
TypeFreePhysicalSlotPage is a page which holds free physical slot ids
*/
const TypeFreePhysicalSlotPage = 4
