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
) (
	citation bal.ResourceLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to save a certificate document.",
	)
	citation, status = v.storage_.WriteContract(certificate)
	return
}

func (v *documentRepository_) SaveDraft(
	draft not.Parameterized,
) (
	citation bal.ResourceLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to save a draft document.",
	)
	citation, status = v.storage_.WriteDraft(draft)
	return
}

func (v *documentRepository_) RetrieveDraft(
	citation bal.ResourceLike,
) (
	draft not.Parameterized,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to retrieve a draft document.",
	)
	draft, status = v.storage_.ReadDraft(citation)
	return
}

func (v *documentRepository_) DiscardDraft(
	citation bal.ResourceLike,
) (
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to discard a draft document.",
	)
	status = v.storage_.DeleteDraft(citation)
	return
}

func (v *documentRepository_) NotarizeDocument(
	name bal.NameLike,
	version bal.VersionLike,
	draft not.Parameterized,
) (
	contract not.ContractLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to notarize a draft document.",
	)
	var citation bal.ResourceLike
	contract = v.notary_.NotarizeDocument(draft)
	citation, status = v.storage_.WriteContract(contract)
	if status != Written {
		return
	}
	status = v.storage_.WriteCitation(name, version, citation)
	return
}

func (v *documentRepository_) RetrieveDocument(
	name bal.NameLike,
	version bal.VersionLike,
) (
	contract not.ContractLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to notarize a contract document.",
	)
	var citation bal.ResourceLike
	citation, status = v.storage_.ReadCitation(name, version)
	if status != Retrieved {
		return
	}
	contract, status = v.storage_.ReadContract(citation)
	return
}

func (v *documentRepository_) CheckoutDocument(
	name bal.NameLike,
	version bal.VersionLike,
	level uint,
) (
	draft not.Parameterized,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to checkout a draft document.",
	)
	var previous bal.ResourceLike
	var contract not.ContractLike
	previous, status = v.storage_.ReadCitation(name, version)
	if status != Retrieved {
		return
	}
	contract, status = v.storage_.ReadContract(previous)
	if status != Retrieved {
		return
	}
	draft = contract.GetContent()
	var entity = draft.GetEntity()
	var type_ = draft.GetType()
	var tag = draft.GetTag()
	var nextVersion = bal.VersionClass().GetNextVersion(
		version,
		level,
	)
	var permissions = draft.GetPermissions()
	draft = not.Draft(
		entity,
		type_,
		tag,
		nextVersion,
		permissions,
		previous,
	)
	return
}

func (v *documentRepository_) CreateBag(
	name bal.NameLike,
	capacity uint,
	leasetime uint,
	permissions bal.ResourceLike,
) (
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to create a message bag.",
	)
	var citation bal.ResourceLike
	var bag = doc.BagClass().Bag(
		name,
		bal.Number(float64(capacity)),
		bal.Number(float64(leasetime)),
		permissions,
	)
	var contract = v.notary_.NotarizeDocument(bag)
	citation, status = v.storage_.WriteContract(contract)
	if status != Written {
		return
	}
	var version = bal.Version()
	status = v.storage_.WriteCitation(name, version, citation)
	return
}

func (v *documentRepository_) RemoveBag(
	name bal.NameLike,
) (
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to delete a message bag.",
	)
	var citation bal.ResourceLike
	var version = bal.Version()
	citation, status = v.storage_.ReadCitation(name, version)
	if status != Retrieved {
		return
	}
	status = v.storage_.DeleteContract(citation)
	if status != Deleted {
		return
	}
	status = v.storage_.DeleteCitation(name, version)
	return
}

func (v *documentRepository_) PostMessage(
	bag bal.NameLike,
	message doc.MessageLike,
) (
	status Status,
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
	var citation bal.ResourceLike
	var contract = v.notary_.NotarizeDocument(message)
	citation, status = v.storage_.WriteContract(contract)
	if status != Written {
		return
	}
	var version = bal.Version()
	status = v.storage_.WriteCitation(name, version, citation)
	return
}

func (v *documentRepository_) RetrieveMessage(
	bag bal.NameLike,
) (
	message not.ContractLike,
	status Status,
) {
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
	for retries := 0; retries < 8; retries++ {
		var sequence bal.Sequential[bal.ResourceLike]
		sequence, status = v.storage_.ListCitations(accessible)
		if status != Retrieved {
			return
		}
		var citations = bal.List[bal.ResourceLike](sequence)
		if citations.IsEmpty() {
			// There are no messages currently in the bag.
			continue
		}

		// Select a message from the bag at random.
		var size = citations.GetSize()
		var index = int(bal.Generator().RandomOrdinal(size))
		var citation = citations.GetValue(index)

		// Read the message.
		message, status = v.storage_.ReadContract(citation)
		if status == Missing {
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
		status = v.storage_.DeleteCitation(accessible, version)
		if status != Deleted {
			return
		}
		status = v.storage_.WriteCitation(processing, version, citation)
		if status != Written {
			return
		}
		break
	}
	return
}

func (v *documentRepository_) AcceptMessage(
	message not.ContractLike,
) (
	status Status,
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
	status = v.storage_.DeleteCitation(processing, version)
	if status != Deleted {
		return
	}

	// Delete the message from the document storage.
	var citation = v.notary_.CiteDocument(content)
	status = v.storage_.DeleteContract(citation)
	return
}

func (v *documentRepository_) RejectMessage(
	message not.ContractLike,
) (
	status Status,
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
	status = v.storage_.DeleteCitation(processing, version)
	if status != Deleted {
		return
	}
	var citation = v.notary_.CiteDocument(content)
	status = v.storage_.WriteCitation(accessible, version, citation)
	return
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
