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
	uti "github.com/craterdog/go-essential-utilities/v8"
	syn "sync"
)

// CLASS INTERFACE

// Access Function

func DocumentRepositoryClass() DocumentRepositoryClassLike {
	return documentRepositoryClass()
}

// Constructor Methods

func (c *documentRepositoryClass_) DocumentRepository(
	group doc.Synchronized,
	notary not.DigitalNotaryLike,
	storage Persistent,
) DocumentRepositoryLike {
	if uti.IsUndefined(group) {
		panic("The \"group\" attribute is required by this class.")
	}
	if uti.IsUndefined(notary) {
		panic("The \"notary\" attribute is required by this class.")
	}
	if uti.IsUndefined(storage) {
		panic("The \"storage\" attribute is required by this class.")
	}
	var instance = &documentRepository_{
		// Initialize the instance attributes.
		group_:   group,
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
	var content = not.Certificate(certificate.GetContent().AsSource())
	var tag = content.GetTag()
	var version = content.GetVersion()
	var name = doc.Name(
		"/certificates/" + tag.AsSource()[1:] + "/" + version.AsSource(),
	)
	defer v.mutex_.Unlock()
	v.mutex_.Lock()
	citation, status = v.storage_.WriteDocument(certificate)
	if status != Success {
		return
	}
	status = v.storage_.WriteCitation(name, citation)
	return
}

func (v *documentRepository_) SaveDraft(
	draft not.DocumentLike,
) (
	citation not.CitationLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to save a draft document.",
	)
	defer v.mutex_.Unlock()
	v.mutex_.Lock()
	citation, status = v.storage_.WriteDraft(draft)
	return
}

func (v *documentRepository_) RetrieveDraft(
	citation not.CitationLike,
) (
	draft not.DocumentLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to retrieve a draft document.",
	)
	defer v.mutex_.Unlock()
	v.mutex_.Lock()
	draft, status = v.storage_.ReadDraft(citation)
	return
}

func (v *documentRepository_) DiscardDraft(
	citation not.CitationLike,
) (
	draft not.DocumentLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to discard a draft document.",
	)
	defer v.mutex_.Unlock()
	v.mutex_.Lock()
	draft, status = v.storage_.DeleteDraft(citation)
	return
}

func (v *documentRepository_) NotarizeDocument(
	name doc.NameLike,
	document not.DocumentLike,
) (
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to notarize a draft document.",
	)
	defer v.mutex_.Unlock()
	v.mutex_.Lock()
	var citation = v.notary_.CiteDocument(document)
	_, status = v.storage_.DeleteDraft(citation)
	v.notary_.NotarizeDocument(document)
	citation, status = v.storage_.WriteDocument(document)
	if status != Success {
		return
	}
	status = v.storage_.WriteCitation(name, citation)
	return
}

func (v *documentRepository_) RetrieveDocument(
	name doc.NameLike,
) (
	document not.DocumentLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to retrieve a document.",
	)
	defer v.mutex_.Unlock()
	v.mutex_.Lock()
	var citation not.CitationLike
	citation, status = v.storage_.ReadCitation(name)
	if status != Success {
		return
	}
	document, status = v.storage_.ReadDocument(citation)
	return
}

func (v *documentRepository_) CheckoutDocument(
	name doc.NameLike,
	level uint,
) (
	draft not.DocumentLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to checkout a draft document.",
	)
	defer v.mutex_.Unlock()
	v.mutex_.Lock()
	var citation not.CitationLike
	citation, status = v.storage_.ReadCitation(name)
	if status != Success {
		return
	}
	var document not.DocumentLike
	document, status = v.storage_.ReadDocument(citation)
	if status != Success {
		return
	}
	var content = document.GetContent()
	var entity = content.GetEntity()
	var type_ = content.GetType()
	var tag = content.GetTag()
	var version = content.GetVersion()
	var nextVersion = doc.VersionClass().GetNextVersion(
		version,
		level,
	)
	var permissions = content.GetPermissions()
	var previous = citation.AsResource()
	content = not.Content(
		entity,
		type_,
		tag,
		nextVersion,
		permissions,
		previous,
	)
	draft = not.Document(content)
	return
}

func (v *documentRepository_) SendMessage(
	bag doc.NameLike,
	message not.DocumentLike,
) (
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to send a message via a bag.",
	)
	defer v.mutex_.Unlock()
	v.mutex_.Lock()
	v.notary_.NotarizeDocument(message)
	status = v.storage_.WriteMessage(bag, message)
	return
}

func (v *documentRepository_) ReceiveMessage(
	bag doc.NameLike,
) (
	message not.DocumentLike,
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to receive a message from a bag.",
	)
	defer v.mutex_.Unlock()
	v.mutex_.Lock()
	message, status = v.storage_.ReadMessage(bag)
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
	defer v.mutex_.Unlock()
	v.mutex_.Lock()
	status = v.storage_.DeleteMessage(bag, message)
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
	defer v.mutex_.Unlock()
	v.mutex_.Lock()
	status = v.storage_.UnreadMessage(bag, message)
	return
}

func (v *documentRepository_) SubscribeEvents(
	bag doc.NameLike,
	type_ doc.NameLike,
) (
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to subscribe to an event type.",
	)
	defer v.mutex_.Unlock()
	v.mutex_.Lock()
	status = v.storage_.WriteSubscription(bag, type_)
	return
}

func (v *documentRepository_) UnsubscribeEvents(
	bag doc.NameLike,
	type_ doc.NameLike,
) (
	status Status,
) {
	defer v.errorCheck(
		"An error occurred while attempting to unsubscribe from an event type.",
	)
	defer v.mutex_.Unlock()
	v.mutex_.Lock()
	status = v.storage_.DeleteSubscription(bag, type_)
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
	defer v.mutex_.Unlock()
	v.mutex_.Lock()
	v.notary_.NotarizeDocument(event)
	var content = event.GetContent()
	var type_ = content.GetType()
	var bags doc.Sequential[doc.NameLike]
	bags, status = v.storage_.ReadSubscriptions(type_)
	var iterator = bags.GetIterator()
	for iterator.HasNext() {
		var bag = iterator.GetNext()
		v.group_.Go(func() {
			var message = v.copyEvent(event)
			v.notary_.NotarizeDocument(message)
			v.storage_.WriteMessage(bag, message)
		})
	}
	return
}

// Attribute Methods

// PROTECTED INTERFACE

// Private Methods

func (v *documentRepository_) copyEvent(
	event not.DocumentLike,
) not.DocumentLike {
	var content = event.GetContent()
	content = not.Content(
		content.GetEntity(),
		content.GetType(),
		doc.Tag(), // Only the tag changes.
		content.GetVersion(),
		content.GetPermissions(),
		content.GetOptionalPrevious(),
	)
	return not.Document(content)
}

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
	mutex_   syn.Mutex
	group_   doc.Synchronized
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
