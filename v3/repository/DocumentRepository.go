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

package repository

import (
	fmt "fmt"
	not "github.com/bali-nebula/go-digital-notary/v3"
	bal "github.com/bali-nebula/go-document-notation/v3"
	fra "github.com/craterdog/go-component-framework/v7"
	uti "github.com/craterdog/go-missing-utilities/v7"
	stc "strconv"
)

// CLASS INTERFACE

// Access Function

func DocumentRepositoryClass() DocumentRepositoryClassLike {
	return documentRepositoryClass()
}

// Constructor Methods

func (c *documentRepositoryClass_) DocumentRepository(
	notary not.NotaryLike,
	storage Persistent,
) DocumentRepositoryLike {
	if uti.IsUndefined(notary) {
		panic("The \"notary\" attribute is required by this class.")
	}
	if uti.IsUndefined(storage) {
		panic("The \"storage\" attribute is required by this class.")
	}
	var instance = &documentRepository_{
		// Initialize the instance attributes.
		notary_:  notary,
		storage_: storage,
	}
	return instance
}

// Constant Methods

// Function Methods

// INSTANCE INTERFACE

// Principal Methods

func (v *documentRepository_) GetClass() DocumentRepositoryClassLike {
	return documentRepositoryClass()
}

func (v *documentRepository_) SaveDraft(
	draft not.DraftLike,
) not.CitationLike {
	var citation = v.storage_.WriteDraft(draft)
	return citation
}

func (v *documentRepository_) RetrieveDraft(
	citation not.CitationLike,
) not.DraftLike {
	var draft = v.storage_.ReadDraft(citation)
	return draft
}

func (v *documentRepository_) DiscardDraft(
	citation not.CitationLike,
) {
	v.storage_.DeleteDraft(citation)
}

func (v *documentRepository_) NotarizeDraft(
	name fra.NameLike,
	draft not.DraftLike,
) not.ContractLike {
	if v.storage_.CitationExists(name) {
		var message = fmt.Sprintf(
			"Attempted to notarize a draft document using an existing name: %v",
			name,
		)
		panic(message)
	}
	var contract = v.notary_.NotarizeDraft(draft)
	var citation = v.storage_.WriteContract(contract)
	v.storage_.WriteCitation(name, citation)
	return contract
}

func (v *documentRepository_) RetrieveContract(
	name fra.NameLike,
) not.ContractLike {
	var citation = v.storage_.ReadCitation(name)
	if uti.IsUndefined(citation) {
		var message = fmt.Sprintf(
			"Attempted to retrieve a non-existent contract with name: %v",
			name,
		)
		panic(message)
	}
	var contract = v.storage_.ReadContract(citation)
	return contract
}

func (v *documentRepository_) CheckoutDraft(
	name fra.NameLike,
	level uint,
) not.DraftLike {
	var citation = v.storage_.ReadCitation(name)
	if uti.IsUndefined(citation) {
		var message = fmt.Sprintf(
			"Attempted to checkout a non-existent contract with name: %v",
			name,
		)
		panic(message)
	}
	var contract = v.storage_.ReadContract(citation)
	var draft = contract.GetDraft()
	var component = draft.GetComponent()
	var type_ = draft.GetType()
	var tag = draft.GetTag()
	var version = draft.GetVersion()
	var permissions = draft.GetPermissions()
	var previous = citation
	var nextVersion = fra.VersionClass().GetNextVersion(
		version,
		uti.Ordinal(level),
	)
	draft = not.Draft(
		component,
		type_,
		tag,
		nextVersion,
		permissions,
		previous,
	)
	return draft
}

func (v *documentRepository_) CreateBag(
	name fra.NameLike,
	permissions fra.ResourceLike,
	capacity uint,
	leasetime uint,
) {
	var component = bal.ParseSource(
		`[
    $capacity: ` + stc.Itoa(int(capacity)) + `
    $leasetime: ` + stc.Itoa(int(leasetime)) + `
]`,
	).GetComponent()
	var type_ = fra.ResourceFromString("<bali:/nebula/types/Bag:v1>")
	var tag = fra.TagWithSize(20)
	var version = fra.VersionFromString("v1")
	var previous not.CitationLike
	var draft = not.Draft(
		component,
		type_,
		tag,
		version,
		permissions,
		previous,
	)
	var contract = v.notary_.NotarizeDraft(draft)
	var citation = v.storage_.WriteContract(contract)
	v.storage_.WriteCitation(name, citation)
}

func (v *documentRepository_) MessageCount(
	bag fra.NameLike,
) uint {
	var citation = v.storage_.ReadCitation(bag)
	return v.storage_.MessageCount(citation)
}

func (v *documentRepository_) PostMessage(
	bag fra.NameLike,
	message not.DraftLike,
) {
	var citation = v.storage_.ReadCitation(bag)
	var contract = v.storage_.ReadContract(citation)
	var source = contract.GetDraft().AsString()
	var document = bal.ParseSource(source)
	source = not.DraftClass().ExtractAttribute("$capacity", document)
	var capacity, _ = stc.Atoi(source)
	var count = v.storage_.MessageCount(citation)
	if count < uint(capacity) {
		var contract = v.notary_.NotarizeDraft(message)
		v.storage_.WriteMessage(citation, contract)
	}
}

func (v *documentRepository_) RetrieveMessage(
	bag fra.NameLike,
) not.DraftLike {
	var citation = v.storage_.ReadCitation(bag)
	var contract = v.storage_.ReadMessage(citation)
	var message = contract.GetDraft()
	return message
}

func (v *documentRepository_) AcceptMessage(
	message not.DraftLike,
) {
	var source = message.AsString()
	var document = bal.ParseSource(source)
	source = not.DraftClass().ExtractAttribute("$bag", document)
	var bag = fra.NameFromString(source)
	var bagCitation = v.storage_.ReadCitation(bag)
	var messageCitation = v.notary_.CiteDraft(message)
	v.storage_.DeleteMessage(bagCitation, messageCitation)
}

func (v *documentRepository_) RejectMessage(
	message not.DraftLike,
) {
	var source = message.AsString()
	var document = bal.ParseSource(source)
	source = not.DraftClass().ExtractAttribute("$bag", document)
	var bag = fra.NameFromString(source)
	var bagCitation = v.storage_.ReadCitation(bag)
	var messageCitation = v.notary_.CiteDraft(message)
	v.storage_.ReturnMessage(bagCitation, messageCitation)
}

func (v *documentRepository_) DeleteBag(
	bag fra.NameLike,
) {
	var citation = v.storage_.ReadCitation(bag)
	v.storage_.DeleteDraft(citation)
}

func (v *documentRepository_) PublishEvent(
	event not.DraftLike,
) {
	var bag = fra.NameFromString("/nebula/bag/Events")
	var citation = v.storage_.ReadCitation(bag)
	var contract = v.notary_.NotarizeDraft(event)
	v.storage_.WriteMessage(citation, contract)
}

// Attribute Methods

// PROTECTED INTERFACE

// Private Methods

// Instance Structure

type documentRepository_ struct {
	// Declare the instance attributes.
	notary_  not.NotaryLike
	storage_ Persistent
}

// Class Structure

type documentRepositoryClass_ struct {
	// Declare the class constants.
}

// Class Reference

func documentRepositoryClass() *documentRepositoryClass_ {
	return documentRepositoryClassReference_
}

var documentRepositoryClassReference_ = &documentRepositoryClass_{
	// Initialize the class constants.
}
