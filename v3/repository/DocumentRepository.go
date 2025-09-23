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
	bal "github.com/bali-nebula/go-bali-documents/v3"
	not "github.com/bali-nebula/go-digital-notary/v3"
	doc "github.com/bali-nebula/go-document-repository/v3/documents"
	uti "github.com/craterdog/go-missing-utilities/v7"
)

// CLASS INTERFACE

// Access Function

func DocumentRepositoryClass() DocumentRepositoryClassLike {
	return documentRepositoryClass()
}

// Constructor Methods

func (c *documentRepositoryClass_) DocumentRepository(
	notary not.DigitalNotaryLike,
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

func (v *documentRepository_) SaveCertificate(
	certificate not.ContractLike,
) bal.ResourceLike {
	defer v.errorCheck(
		"An error occurred while attempting to save a certificate document.",
	)
	var citation = v.storage_.WriteContract(certificate)
	return citation
}

func (v *documentRepository_) SaveDraft(
	draft not.Parameterized,
) bal.ResourceLike {
	defer v.errorCheck(
		"An error occurred while attempting to save a draft document.",
	)
	var citation = v.storage_.WriteDraft(draft)
	return citation
}

func (v *documentRepository_) RetrieveDraft(
	citation bal.ResourceLike,
) not.Parameterized {
	defer v.errorCheck(
		"An error occurred while attempting to retrieve a draft document.",
	)
	return v.storage_.ReadDraft(citation)
}

func (v *documentRepository_) DiscardDraft(
	citation bal.ResourceLike,
) {
	defer v.errorCheck(
		"An error occurred while attempting to discard a draft document.",
	)
	v.storage_.DeleteDraft(citation)
}

func (v *documentRepository_) NotarizeDocument(
	name bal.NameLike,
	version bal.VersionLike,
	draft not.Parameterized,
) not.ContractLike {
	defer v.errorCheck(
		"An error occurred while attempting to notarize a draft document.",
	)
	var contract = v.notary_.NotarizeDocument(draft)
	var citation = v.storage_.WriteContract(contract)
	v.storage_.WriteCitation(name, version, citation)
	return contract
}

func (v *documentRepository_) RetrieveDocument(
	name bal.NameLike,
	version bal.VersionLike,
) not.ContractLike {
	defer v.errorCheck(
		"An error occurred while attempting to notarize a contract document.",
	)
	var citation = v.storage_.ReadCitation(name, version)
	return v.storage_.ReadContract(citation)
}

func (v *documentRepository_) CheckoutDocument(
	name bal.NameLike,
	version bal.VersionLike,
	level uint,
) not.Parameterized {
	defer v.errorCheck(
		"An error occurred while attempting to checkout a draft document.",
	)
	var citation = v.storage_.ReadCitation(name, version)
	var content = v.storage_.ReadContract(citation).GetContent()
	var nextVersion = bal.VersionClass().GetNextVersion(
		version,
		level,
	)
	var draft = not.Draft(
		content.GetEntity(),
		content.GetType(),
		content.GetTag(),
		nextVersion,
		content.GetPermissions(),
		citation,
	)
	return draft
}

func (v *documentRepository_) CreateBag(
	name bal.NameLike,
	capacity uint,
	leasetime uint,
	permissions bal.ResourceLike,
) {
	defer v.errorCheck(
		"An error occurred while attempting to create a message bag.",
	)
	var bag = doc.BagClass().Bag(
		name,
		bal.Number(float64(capacity)),
		bal.Number(float64(leasetime)),
		permissions,
	)
	var contract = v.notary_.NotarizeDocument(bag)
	var citation = v.storage_.WriteContract(contract)
	var version = bal.Version()
	v.storage_.WriteCitation(name, version, citation)
}

func (v *documentRepository_) RemoveBag(
	name bal.NameLike,
) {
	defer v.errorCheck(
		"An error occurred while attempting to delete a message bag.",
	)
	var version = bal.Version()
	var citation = v.storage_.ReadCitation(name, version)
	v.storage_.DeleteContract(citation)
	v.storage_.DeleteCitation(name, version)
}

func (v *documentRepository_) PostMessage(
	bag bal.NameLike,
	message doc.MessageLike,
) {
	defer v.errorCheck(
		"An error occurred while attempting to send a message via a bag.",
	)
	// TBD - Add checks for bag capacity violations.
	message.SetObject(bag, bal.Symbol("bag"))
	var name = bal.NameClass().Concatenate(
		bag,
		bal.Name("/accessible/"+message.GetTag().AsString()[1:]),
	)
	var version = bal.Version()
	var contract = v.notary_.NotarizeDocument(message)
	var citation = v.storage_.WriteContract(contract)
	v.storage_.WriteCitation(name, version, citation)
}

func (v *documentRepository_) RetrieveMessage(
	bag bal.NameLike,
) not.ContractLike {
	defer v.errorCheck(
		"An error occurred while attempting to retrieve a message from a bag.",
	)
	var accessible = bal.NameClass().Concatenate(
		bag,
		bal.Name("/accessible"),
	)
	var processing = bal.NameClass().Concatenate(
		bag,
		bal.Name("/processing"),
	)
	var message not.ContractLike
	for retries := 0; retries < 8; retries++ {
		var citations = bal.List[bal.ResourceLike](
			v.storage_.ListCitations(accessible),
		)
		if citations.IsEmpty() {
			// There are no messages currently in the bag.
			continue
		}

		// Select a message from the bag at random.
		var size = citations.GetSize()
		var index = int(bal.Generator().RandomOrdinal(size))
		var citation = citations.GetValue(index)

		// Read the message.
		message = v.storage_.ReadContract(citation)
		if uti.IsUndefined(message) {
			// Another process got there first.
			continue
		}

		// Move the message citation from accessible to processing.
		var tag = bal.Name(
			"/" + message.GetContent().GetTag().AsString()[1:],
		)
		accessible = bal.NameClass().Concatenate(accessible, tag)
		processing = bal.NameClass().Concatenate(processing, tag)
		var version = bal.Version()
		v.storage_.DeleteCitation(accessible, version)
		v.storage_.WriteCitation(processing, version, citation)
		break
	}
	return message
}

func (v *documentRepository_) AcceptMessage(
	message not.ContractLike,
) {
	defer v.errorCheck(
		"An error occurred while attempting to accept a processed message.",
	)

	// Delete the message citation from the document storage.
	var content = message.GetContent()
	var bag = content.GetObject(
		bal.Symbol("bag"),
	).GetComponent().GetEntity().(bal.NameLike)
	var tag = content.GetTag().AsString()[1:]
	var processing = bal.NameClass().Concatenate(
		bag,
		bal.Name("/processing/"+tag),
	)
	var version = bal.Version()
	v.storage_.DeleteCitation(processing, version)

	// Delete the message from the document storage.
	var citation = v.notary_.CiteDocument(content)
	v.storage_.DeleteContract(citation)
}

func (v *documentRepository_) RejectMessage(
	message not.ContractLike,
) {
	defer v.errorCheck(
		"An error occurred while attempting to reject a retrieved  message.",
	)
	var content = message.GetContent()
	var bag = content.GetObject(
		bal.Symbol("bag"),
	).GetComponent().GetEntity().(bal.NameLike)
	var tag = content.GetTag().AsString()[1:]
	var accessible = bal.NameClass().Concatenate(
		bag,
		bal.Name("/accessible/"+tag),
	)
	var processing = bal.NameClass().Concatenate(
		bag,
		bal.Name("/processing/"+tag),
	)
	var version = bal.Version()
	v.storage_.DeleteCitation(processing, version)
	var citation = v.notary_.CiteDocument(content)
	v.storage_.WriteCitation(accessible, version, citation)
}

// Attribute Methods

// PROTECTED INTERFACE

// Private Methods

func (v *documentRepository_) errorCheck(
	message string,
) {
	if e := recover(); e != nil {
		message = fmt.Sprintf(
			"DocumentRepository: %s:\n    %v",
			message,
			e,
		)
		panic(message)
	}
}

// Instance Structure

type documentRepository_ struct {
	// Declare the instance attributes.
	notary_  not.DigitalNotaryLike
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
