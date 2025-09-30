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
	doc "github.com/bali-nebula/go-bali-documents/v3"
	not "github.com/bali-nebula/go-digital-notary/v3"
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
	certificate not.DocumentLike,
) (
	citation not.CitationLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to save a certificate document.",
	)
	citation, status = v.storage_.CreateDocument(certificate)
	var content = certificate.GetContent()
	var tag = content.GetTag()
	var version = content.GetVersion()
	var name = doc.Name("/certificates/" + tag.AsString()[1:])
	status = v.storage_.CreateCitation(name, version, citation)
	return
}

func (v *documentRepository_) SaveDraft(
	document not.DocumentLike,
) (
	citation not.CitationLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to save a draft document.",
	)
	citation, status = v.storage_.UpdateDocument(document)
	return
}

func (v *documentRepository_) RetrieveDraft(
	citation not.CitationLike,
) (
	document not.DocumentLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to retrieve a draft document.",
	)
	document, status = v.storage_.ReadDocument(citation)
	return
}

func (v *documentRepository_) DiscardDraft(
	citation not.CitationLike,
) (
	document not.DocumentLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to discard a draft document.",
	)
	document, status = v.storage_.DeleteDocument(citation)
	return
}

func (v *documentRepository_) NotarizeDocument(
	name doc.NameLike,
	version doc.VersionLike,
	document not.DocumentLike,
) (
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to notarize a draft document.",
	)
	var citation not.CitationLike
	document = v.notary_.NotarizeDocument(document)
	citation, status = v.storage_.CreateDocument(document)
	if status != Success {
		return
	}
	status = v.storage_.CreateCitation(name, version, citation)
	return
}

func (v *documentRepository_) RetrieveDocument(
	name doc.NameLike,
	version doc.VersionLike,
) (
	document not.DocumentLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to notarize a document.",
	)
	var citation not.CitationLike
	citation, status = v.storage_.ReadCitation(name, version)
	if status != Success {
		return
	}
	document, status = v.storage_.ReadDocument(citation)
	return
}

func (v *documentRepository_) CheckoutDocument(
	name doc.NameLike,
	version doc.VersionLike,
	level uint,
) (
	document not.DocumentLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to checkout a draft document.",
	)
	var previous not.CitationLike
	previous, status = v.storage_.ReadCitation(name, version)
	if status != Success {
		return
	}
	document, status = v.storage_.ReadDocument(previous)
	if status != Success {
		return
	}
	var content = document.GetContent()
	var entity = content.GetEntity()
	var type_ = content.GetType()
	var tag = content.GetTag()
	var nextVersion = doc.VersionClass().GetNextVersion(
		version,
		level,
	)
	var permissions = content.GetPermissions()
	content = not.Content(
		entity,
		type_,
		tag,
		nextVersion,
		previous.AsResource(),
		permissions,
	)
	document = not.Document(content)
	return
}

func (v *documentRepository_) CreateBag(
	name doc.NameLike,
	bag not.DocumentLike,
) (
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to create a message bag.",
	)
	var citation not.CitationLike
	bag = v.notary_.NotarizeDocument(bag)
	citation, status = v.storage_.CreateDocument(bag)
	if status != Success {
		return
	}
	var version = doc.Version()
	status = v.storage_.CreateCitation(name, version, citation)
	return
}

func (v *documentRepository_) RemoveBag(
	name doc.NameLike,
) (
	bag not.DocumentLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to delete a message bag.",
	)
	var citation not.CitationLike
	var version = doc.Version()
	citation, status = v.storage_.ReadCitation(name, version)
	if status != Success {
		return
	}
	bag, status = v.storage_.DeleteDocument(citation)
	if status != Success {
		return
	}
	_, status = v.storage_.DeleteCitation(name, version)
	return
}

func (v *documentRepository_) PostMessage(
	bag doc.NameLike,
	message not.DocumentLike,
) (
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to send a message via a bag.",
	)
	// TBD - Add checks for bag capacity violations.
	message = v.notary_.NotarizeDocument(message)
	var content = message.GetContent()
	var name = doc.NameClass().Concatenate(
		bag,
		doc.Name("/accessible/"+content.GetTag().AsString()[1:]),
	)
	var citation not.CitationLike
	citation, status = v.storage_.CreateDocument(message)
	if status != Success {
		return
	}
	var version = doc.Version()
	status = v.storage_.CreateCitation(name, version, citation)
	return
}

func (v *documentRepository_) RetrieveMessage(
	bag doc.NameLike,
) (
	message not.DocumentLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to retrieve a message from a bag.",
	)
	var accessible = doc.Name(bag.AsString() + "/accessible")
	var processing = doc.Name(bag.AsString() + "/processing")
	for retries := 0; retries < 8; retries++ {
		var sequence doc.Sequential[not.CitationLike]
		sequence, status = v.storage_.ListCitations(accessible)
		if status != Success {
			return
		}
		var citations = doc.List[not.CitationLike](sequence)
		if citations.IsEmpty() {
			// There are no messages currently in the bag.
			continue
		}

		// Select a message from the bag at random.
		var size = citations.GetSize()
		var index = int(doc.Generator().RandomOrdinal(size))
		var citation = citations.GetValue(index)

		// Move the message citation from accessible to processing.
		var tag = doc.Name("/" + citation.GetTag().AsString()[1:])
		accessible = doc.NameClass().Concatenate(accessible, tag)
		processing = doc.NameClass().Concatenate(processing, tag)
		var version = doc.Version()
		citation, status = v.storage_.MoveCitation(
			accessible,
			processing,
			version,
		)
		if status == Missing {
			// Another process got there first.
			continue
		}

		// Read the message.
		message, status = v.storage_.ReadDocument(citation)
		if status != Success {
			return
		}
		break
	}
	return
}

func (v *documentRepository_) AcceptMessage(
	bag doc.NameLike,
	message not.DocumentLike,
) (
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to accept a processed message.",
	)

	// Delete the message citation from the document storage.
	var content = message.GetContent()
	var tag = content.GetTag().AsString()[1:]
	var name = doc.NameClass().Concatenate(
		bag,
		doc.Name("/processing/"+tag),
	)
	var version = doc.Version()
	var citation not.CitationLike
	citation, status = v.storage_.DeleteCitation(name, version)
	if status != Success {
		return
	}

	// Delete the message from the document storage.
	_, status = v.storage_.DeleteDocument(citation)
	return
}

func (v *documentRepository_) RejectMessage(
	bag doc.NameLike,
	message not.DocumentLike,
) (
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to reject a retrieved message.",
	)
	var content = message.GetContent()
	var tag = content.GetTag().AsString()[1:]
	var accessible = doc.NameClass().Concatenate(
		bag,
		doc.Name("/accessible/"+tag),
	)
	var processing = doc.NameClass().Concatenate(
		bag,
		doc.Name("/processing/"+tag),
	)
	var version = doc.Version()
	_, status = v.storage_.MoveCitation(processing, accessible, version)
	return
}

func (v *documentRepository_) PublishEvent(
	event not.DocumentLike,
) (
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to publish an event.",
	)
	event = v.notary_.NotarizeDocument(event)
	var content = event.GetContent()
	var name = doc.Name("/events/" + content.GetTag().AsString()[1:])
	var version = content.GetVersion()
	var citation not.CitationLike
	citation, status = v.storage_.CreateDocument(event)
	if status != Success {
		return
	}
	status = v.storage_.CreateCitation(name, version, citation)
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
			"DocumentRepository: %s:\n    %s",
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
