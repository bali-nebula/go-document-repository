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
	ran "crypto/rand"
	not "github.com/bali-nebula/go-digital-notary/v3"
	fra "github.com/craterdog/go-component-framework/v7"
	uti "github.com/craterdog/go-missing-utilities/v7"
	big "math/big"
	osx "os"
	sts "strings"
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
	resource fra.ResourceLike,
) bool {
	var path = v.directory_ + "citations" + v.extractPath(resource)
	return uti.PathExists(path)
}

func (v *localStorage_) ReadCitation(
	resource fra.ResourceLike,
) not.CitationLike {
	var path = v.directory_ + "citations" + v.extractPath(resource)
	var source = uti.ReadFile(path)
	var citation = not.CitationFromString(source)
	return citation
}

func (v *localStorage_) WriteCitation(
	resource fra.ResourceLike,
	citation not.CitationLike,
) {
	var path = v.directory_ + "citations" + v.extractPath(resource)
	var source = citation.AsString()
	uti.WriteFile(path, source)
}

func (v *localStorage_) DraftExists(
	citation not.CitationLike,
) bool {
	var path = v.directory_ + "drafts/" + v.dereference(citation)
	return uti.PathExists(path)
}

func (v *localStorage_) ReadDraft(
	citation not.CitationLike,
) not.DraftLike {
	var path = v.directory_ + "drafts/" + v.dereference(citation)
	var source = uti.ReadFile(path)
	var draft = not.DraftFromString(source)
	return draft
}

func (v *localStorage_) WriteDraft(
	draft not.DraftLike,
) not.CitationLike {
	var citation = v.notary_.CiteDraft(draft)
	var path = v.directory_ + "drafts/" + v.dereference(citation)
	var source = draft.AsString()
	uti.WriteFile(path, source)
	return citation
}

func (v *localStorage_) DeleteDraft(
	citation not.CitationLike,
) {
	var path = v.directory_ + "drafts/" + v.dereference(citation)
	uti.RemovePath(path)
}

func (v *localStorage_) ContractExists(
	citation not.CitationLike,
) bool {
	var path = v.directory_ + "contracts/" + v.dereference(citation)
	return uti.PathExists(path)
}

func (v *localStorage_) ReadContract(
	citation not.CitationLike,
) not.ContractLike {
	var path = v.directory_ + "contracts/" + v.dereference(citation)
	var source = uti.ReadFile(path)
	var contract = not.ContractFromString(source)
	return contract
}

func (v *localStorage_) WriteContract(
	contract not.ContractLike,
) not.CitationLike {
	var draft = contract.GetDraft()
	var citation = v.notary_.CiteDraft(draft)
	var path = v.directory_ + "contracts/" + v.dereference(citation)
	var source = contract.AsString()
	uti.WriteFile(path, source)
	return citation
}

func (v *localStorage_) MessageCount(
	bag not.CitationLike,
) uint {
	var path = v.directory_ + "bags/available/" + v.dereference(bag)
	var messages = uti.ReadDirectory(path)
	return uint(len(messages))
}

func (v *localStorage_) ReadMessage(
	bag not.CitationLike,
) not.ContractLike {
	var contract not.ContractLike
	for tries := 5; tries > 0; tries-- {
		var available = v.directory_ + "bags/" + v.dereference(bag) + "/available/"
		var processing = v.directory_ + "bags/" + v.dereference(bag) + "/processing/"
		var messages = uti.ReadDirectory(available)
		var count = len(messages)
		if count == 0 {
			// No messages in the bag, try again...
			continue
		}
		var index, _ = ran.Int(ran.Reader, big.NewInt(int64(count)))
		var message = messages[index.Int64()]
		var err = osx.Rename(available+message, processing+message)
		if err != nil {
			// Someone got there first, try again...
			continue
		}
		var source = uti.ReadFile(processing + message)
		contract = not.ContractFromString(source)
		break
	}
	return contract
}

func (v *localStorage_) WriteMessage(
	bag not.CitationLike,
	message not.ContractLike,
) {
	var path = v.directory_ + "bags/" + v.dereference(bag) + "/available/"
	var draft = message.GetDraft()
	var citation = v.notary_.CiteDraft(draft)
	path += v.dereference(citation)
	var source = message.AsString()
	uti.WriteFile(path, source)
}

func (v *localStorage_) DeleteMessage(
	bag not.CitationLike,
	message not.CitationLike,
) {
	var path = v.directory_ + "bags/" + v.dereference(bag) + "/processing/"
	path += v.dereference(message)
	uti.RemovePath(path)
}

func (v *localStorage_) ReleaseMessage(
	bag not.CitationLike,
	message not.CitationLike,
) {
	var processing = v.directory_ + "bags/" + v.dereference(bag) + "/processing/"
	var available = v.directory_ + "bags/" + v.dereference(bag) + "/available/"
	var filename = v.dereference(message)
	// If Rename() fails the lease has already expired so nothing else to do.
	var _ = osx.Rename(processing+filename, available+filename)
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

func (v *localStorage_) extractPath(
	resource fra.ResourceLike,
) string {
	return sts.ReplaceAll(resource.GetPath(), ":", "/")
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
