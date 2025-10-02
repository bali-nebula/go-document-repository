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
	citation, status = v.storage_.WriteDocument(certificate)
	var content = certificate.GetContent()
	var tag = content.GetTag()
	var version = content.GetVersion()
	var name = doc.Name("/certificates/" + tag.AsString()[1:])
	status = v.storage_.WriteCitation(name, version, citation)
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
	citation, status = v.storage_.WriteDocument(document)
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
	citation, status = v.storage_.WriteDocument(document)
	if status != Success {
		return
	}
	status = v.storage_.WriteCitation(name, version, citation)
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
	_, status = v.storage_.WriteDocument(document)
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
	var citation not.CitationLike
	message = v.notary_.NotarizeDocument(message)
	citation, status = v.storage_.WriteDocument(message)
	if status != Success {
		return
	}
	status = v.storage_.WriteMessage(bag, citation)
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
	var citation not.CitationLike
	citation, status = v.storage_.ReadMessage(bag)
	if status != Success {
		return
	}
	message, status = v.storage_.ReadDocument(citation)
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
	var citation = v.notary_.CiteDocument(message)
	status = v.storage_.DeleteMessage(bag, citation)
	if status != Success {
		return
	}
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
	var citation = v.notary_.CiteDocument(message)
	status = v.storage_.UnreadMessage(bag, citation)
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
	citation, status = v.storage_.WriteDocument(event)
	if status != Success {
		return
	}
	status = v.storage_.WriteCitation(name, version, citation)
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
