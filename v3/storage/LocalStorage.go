/*
................................................................................
.    Copyright (c) 2009-2025 Crater Dog Technologies.  All Rights Reserved.    .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See https://opensource.org/license/MIT)                        .
................................................................................
*/

package storage

import (
	not "github.com/bali-nebula/go-digital-notary/v3"
	fra "github.com/craterdog/go-component-framework/v7"
	uti "github.com/craterdog/go-missing-utilities/v7"
)

// CLASS INTERFACE

// Access Function

func LocalStorageClass() LocalStorageClassLike {
	return localStorageClass()
}

// Constructor Methods

func (c *localStorageClass_) LocalStorage(
	notary not.DigitalNotaryLike,
	directory string,
) LocalStorageLike {
	if uti.IsUndefined(notary) {
		panic("The \"notary\" attribute is required by this class.")
	}
	if uti.IsUndefined(directory) {
		panic("The \"directory\" attribute is required by this class.")
	}
	var instance = &localStorage_{
		// Initialize the instance attributes.
		notary_:    notary,
		directory_: directory,
	}
	return instance
}

// Constant Methods

// Function Methods

// INSTANCE INTERFACE

// Principal Methods

func (v *localStorage_) GetClass() LocalStorageClassLike {
	return localStorageClass()
}

// Attribute Methods

// Persistent Methods

func (v *localStorage_) CitationExists(
	name fra.NameLike,
) bool {
	var path = v.directory_ + "citations" + name.AsString()
	return uti.PathExists(path)
}

func (v *localStorage_) ReadCitation(
	name fra.NameLike,
) not.CitationLike {
	var path = v.directory_ + "citations" + name.AsString()
	var source = uti.ReadFile(path)
	var citation = not.CitationFromString(source)
	return citation
}

func (v *localStorage_) WriteCitation(
	name fra.NameLike,
	citation not.CitationLike,
) {
	var path = v.directory_ + "citations" + name.AsString()
	var source = citation.AsString()
	uti.WriteFile(path, source)
}

func (v *localStorage_) DraftExists(
	citation not.CitationLike,
) bool {
	var path = v.directory_ + "drafts" + v.dereference(citation)
	return uti.PathExists(path)
}

func (v *localStorage_) ReadDraft(
	citation not.CitationLike,
) not.DraftLike {
	var path = v.directory_ + "drafts" + v.dereference(citation)
	var source = uti.ReadFile(path)
	var draft = not.DraftFromString(source)
	return draft
}

func (v *localStorage_) WriteDraft(
	draft not.DraftLike,
) not.CitationLike {
	var citation = v.notary_.CiteDraft(draft)
	var path = v.directory_ + "drafts" + v.dereference(citation)
	var source = draft.AsString()
	uti.WriteFile(path, source)
	return citation
}

func (v *localStorage_) DeleteDraft(
	citation not.CitationLike,
) {
	var path = v.directory_ + "drafts" + v.dereference(citation)
	uti.RemovePath(path)
}

func (v *localStorage_) ContractExists(
	citation not.CitationLike,
) bool {
	var path = v.directory_ + "contracts" + v.dereference(citation)
	return uti.PathExists(path)
}

func (v *localStorage_) ReadContract(
	citation not.CitationLike,
) not.ContractLike {
	var path = v.directory_ + "contracts" + v.dereference(citation)
	var source = uti.ReadFile(path)
	var contract = not.ContractFromString(source)
	return contract
}

func (v *localStorage_) WriteContract(
	contract not.ContractLike,
) not.CitationLike {
	var draft = contract.GetDraft()
	var citation = v.notary_.CiteDraft(draft)
	var path = v.directory_ + "contracts" + v.dereference(citation)
	var source = contract.AsString()
	uti.WriteFile(path, source)
	return citation
}

func (v *localStorage_) MessageCount(
	bag not.CitationLike,
) uint {
	var result_ uint
	// TBD - Add the method implementation.
	return result_
}

func (v *localStorage_) ReadMessage(
	bag not.CitationLike,
) not.ContractLike {
	var result_ not.ContractLike
	// TBD - Add the method implementation.
	return result_
}

func (v *localStorage_) WriteMessage(
	bag not.CitationLike,
	message not.ContractLike,
) {
	// TBD - Add the method implementation.
}

func (v *localStorage_) DeleteMessage(
	bag not.CitationLike,
	message not.CitationLike,
) {
	// TBD - Add the method implementation.
}

func (v *localStorage_) ReleaseMessage(
	bag not.CitationLike,
	message not.CitationLike,
) {
	// TBD - Add the method implementation.
}

// PROTECTED INTERFACE

// Private Methods

func (v *localStorage_) dereference(
	citation not.CitationLike,
) string {
	var tag = citation.GetTag()
	var version = citation.GetVersion()
	return tag.AsString()[1:] + "/" + version.AsString()
}

// Instance Structure

type localStorage_ struct {
	// Declare the instance attributes.
	notary_    not.DigitalNotaryLike
	directory_ string
}

// Class Structure

type localStorageClass_ struct {
	// Declare the class constants.
}

// Class Reference

func localStorageClass() *localStorageClass_ {
	return localStorageClassReference_
}

var localStorageClassReference_ = &localStorageClass_{
	// Initialize the class constants.
}
