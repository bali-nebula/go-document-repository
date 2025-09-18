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
	doc "github.com/bali-nebula/go-document-repository/v3/documents"
	fra "github.com/craterdog/go-component-framework/v7"
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
	certificate not.CertificateLike,
) fra.ResourceLike {
	defer v.errorCheck(
		"An error occurred while attempting to save a certificate document.",
	)
	var citation = v.storage_.WriteContract(certificate)
	return citation
}

func (v *documentRepository_) SaveDraft(
	draft not.Parameterized,
) fra.ResourceLike {
	defer v.errorCheck(
		"An error occurred while attempting to save a draft document.",
	)
	var citation = v.storage_.WriteDraft(draft)
	return citation
}

func (v *documentRepository_) RetrieveDraft(
	citation fra.ResourceLike,
) not.Parameterized {
	defer v.errorCheck(
		"An error occurred while attempting to retrieve a draft document.",
	)
	return v.storage_.ReadDraft(citation)
}

func (v *documentRepository_) DiscardDraft(
	citation fra.ResourceLike,
) {
	defer v.errorCheck(
		"An error occurred while attempting to discard a draft document.",
	)
	v.storage_.DeleteDraft(citation)
}

func (v *documentRepository_) NotarizeDocument(
	name fra.NameLike,
	version fra.VersionLike,
	draft not.Parameterized,
) not.Notarized {
	defer v.errorCheck(
		"An error occurred while attempting to notarize a draft document.",
	)
	var contract = v.notary_.NotarizeDraft(draft)
	var citation = v.storage_.WriteContract(contract)
	v.storage_.WriteCitation(name, version, citation)
	return contract
}

func (v *documentRepository_) RetrieveDocument(
	name fra.NameLike,
	version fra.VersionLike,
) not.Notarized {
	defer v.errorCheck(
		"An error occurred while attempting to notarize a contract document.",
	)
	var citation = v.storage_.ReadCitation(name, version)
	return v.storage_.ReadContract(citation)
}

func (v *documentRepository_) CheckoutDocument(
	name fra.NameLike,
	version fra.VersionLike,
	level uti.Cardinal,
) not.Parameterized {
	defer v.errorCheck(
		"An error occurred while attempting to checkout a draft document.",
	)
	var citation = v.storage_.ReadCitation(name, version)
	var content = v.storage_.ReadContract(citation).GetContent()
	var nextVersion = fra.VersionClass().GetNextVersion(
		version,
		uti.Ordinal(level), // TBD - This should be Cardinal since 0 is allowed.
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
	name fra.NameLike,
	capacity uti.Cardinal,
	leasetime uti.Cardinal,
	permissions fra.ResourceLike,
) {
	defer v.errorCheck(
		"An error occurred while attempting to create a message bag.",
	)
	var bag = doc.BagClass().Bag(
		name,
		fra.Number(complex(float64(capacity), 0)), // TBD - Need NumberFromInteger().
		fra.Number(complex(float64(leasetime), 0)),
		permissions,
	)
	var contract = v.notary_.NotarizeDraft(bag)
	var citation = v.storage_.WriteContract(contract)
	var version fra.VersionLike // Bags don't have a version number.
	v.storage_.WriteCitation(name, version, citation)
}

func (v *documentRepository_) RemoveBag(
	name fra.NameLike,
) {
	defer v.errorCheck(
		"An error occurred while attempting to delete a message bag.",
	)
	var version fra.VersionLike // Bags don't have a version number.
	var citation = v.storage_.ReadCitation(name, version)
	v.storage_.DeleteContract(citation)
	v.storage_.DeleteCitation(name, version)
}

func (v *documentRepository_) PostMessage(
	bag fra.NameLike,
	message doc.MessageLike,
) {
	defer v.errorCheck(
		"An error occurred while attempting to send a message via a bag.",
	)
	// TBD - Add checks for bag capacity violations.
	message.SetObject(bag, fra.Symbol("bag"))
	var name = fra.NameClass().Concatenate(
		bag,
		fra.NameFromString("/accessible/"+message.GetTag().AsString()[1:]),
	)
	var version fra.VersionLike // Bags don't have a version number.
	var contract = v.notary_.NotarizeDraft(message)
	var citation = v.storage_.WriteContract(contract)
	v.storage_.WriteCitation(name, version, citation)
}

func (v *documentRepository_) RetrieveMessage(
	bag fra.NameLike,
) not.Notarized {
	defer v.errorCheck(
		"An error occurred while attempting to retrieve a message from a bag.",
	)
	var accessible = fra.NameClass().Concatenate(
		bag,
		fra.NameFromString("/accessible"),
	)
	var processing = fra.NameClass().Concatenate(
		bag,
		fra.NameFromString("/processing"),
	)
	var message not.Notarized
	for retries := 0; retries < 8; retries++ {
		var citations = fra.ListFromSequence(v.storage_.ListCitations(accessible))
		if citations.IsEmpty() {
			// There are no messages currently in the bag.
			continue
		}

		// Select a message from the bag at random.
		var size = uti.Ordinal(citations.GetSize())
		var index = uti.Index(fra.Generator().RandomOrdinal(size))
		var citation = citations.GetValue(index)

		// Read the message.
		message = v.storage_.ReadContract(citation)
		if uti.IsUndefined(message) {
			// Another process got there first.
			continue
		}

		// Move the message citation from accessible to processing.
		var tag = fra.NameFromString(
			"/" + message.GetContent().GetTag().AsString()[1:],
		)
		accessible = fra.NameClass().Concatenate(accessible, tag)
		processing = fra.NameClass().Concatenate(processing, tag)
		var none fra.VersionLike
		v.storage_.DeleteCitation(accessible, none)
		v.storage_.WriteCitation(processing, none, citation)
		break
	}
	return message
}

func (v *documentRepository_) AcceptMessage(
	message not.Notarized,
) {
	defer v.errorCheck(
		"An error occurred while attempting to accept a processed message.",
	)

	// Delete the message citation from the document storage.
	var content = message.GetContent().(doc.MessageLike)
	var bag = content.GetBag()
	var tag = content.GetTag()
	var processing = fra.NameClass().Concatenate(
		bag,
		fra.NameFromString("/processing/"+tag.AsString()),
	)
	var none fra.VersionLike
	v.storage_.DeleteCitation(processing, none)

	// Delete the message from the document storage.
	var citation = v.notary_.CiteDraft(content)
	v.storage_.DeleteContract(citation)
}

func (v *documentRepository_) RejectMessage(
	message not.Notarized,
) {
	defer v.errorCheck(
		"An error occurred while attempting to reject a retrieved  message.",
	)
	var content = message.GetContent().(doc.MessageLike)
	var bag = content.GetBag()
	var tag = content.GetTag()
	var accessible = fra.NameClass().Concatenate(
		bag,
		fra.NameFromString("/accessible/"+tag.AsString()),
	)
	var processing = fra.NameClass().Concatenate(
		bag,
		fra.NameFromString("/processing/"+tag.AsString()),
	)
	var none fra.VersionLike
	v.storage_.DeleteCitation(processing, none)
	var citation = v.notary_.CiteDraft(content)
	v.storage_.WriteCitation(accessible, none, citation)
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
