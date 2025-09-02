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
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to save a certificate document.",
	)

	// Write the certificate document out to document storage.
	var citation = v.storage_.WriteDocument(certificate)
	return citation
}

func (v *documentRepository_) SaveDraft(
	draft not.Parameterized,
) fra.ResourceLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to save a draft document.",
	)

	// Write the draft document out to document storage.
	var citation = v.storage_.WriteDraft(draft)
	return citation
}

func (v *documentRepository_) RetrieveDraft(
	citation fra.ResourceLike,
) not.Parameterized {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to retrieve a draft document.",
	)

	// Read the draft document from document storage.
	return v.storage_.ReadDraft(citation)
}

func (v *documentRepository_) DiscardDraft(
	citation fra.ResourceLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to discard a draft document.",
	)

	// Delete the draft document from document storage.
	v.storage_.DeleteDraft(citation)
}

func (v *documentRepository_) NotarizeDocument(
	name fra.ResourceLike,
	draft not.Parameterized,
) not.Notarized {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to notarize a draft document.",
	)

	// Make sure the name of the document is unique.
	if v.storage_.CitationExists(name) {
		var message = fmt.Sprintf(
			"Attempted to notarize a draft document using an existing name: %v",
			name,
		)
		panic(message)
	}

	// Notarize the draft document.
	var contract = v.notary_.NotarizeDraft(draft)

	// Write the notarized contract out to document storage.
	var citation = v.storage_.WriteDocument(contract)

	// Write the citation to the contract out to document storage.
	v.storage_.WriteCitation(name, citation)
	return contract
}

func (v *documentRepository_) RetrieveDocument(
	name fra.ResourceLike,
) not.Notarized {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to notarize a contract document.",
	)

	// Make sure the name of the contract document exists.
	var citation = v.storage_.ReadCitation(name)
	if uti.IsUndefined(citation) {
		var message = fmt.Sprintf(
			"Attempted to retrieve a non-existent contract with name: %v",
			name,
		)
		panic(message)
	}

	// Read the contract document from document storage.
	return v.storage_.ReadDocument(citation)
}

func (v *documentRepository_) CheckoutDocument(
	name fra.ResourceLike,
	level uti.Cardinal,
) not.Parameterized {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to checkout a draft document.",
	)

	// Make sure the name of the contract document exists.
	var citation = v.storage_.ReadCitation(name)
	if uti.IsUndefined(citation) {
		var message = fmt.Sprintf(
			"Attempted to checkout a non-existent contract with name: %v",
			name,
		)
		panic(message)
	}

	// Read the contract document from document storage.
	var draft = v.storage_.ReadDocument(citation).GetContent()

	// Create a draft copy of the contract document with an updated version.
	var entity = draft.GetEntity()
	var type_ = draft.GetType()
	var tag = draft.GetTag()
	var version = draft.GetVersion()
	var nextVersion = fra.VersionClass().GetNextVersion(
		version,
		uti.Ordinal(level),
	)
	var permissions = draft.GetPermissions()
	var previous = citation
	draft = not.Draft(
		entity,
		type_,
		tag,
		nextVersion,
		permissions,
		previous,
	)
	return draft
}

func (v *documentRepository_) CreateBag(
	name fra.ResourceLike,
	permissions fra.ResourceLike,
	capacity uti.Cardinal,
	leasetime uti.Cardinal,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to create a message bag.",
	)

	// Create a new message bag.
	var source = `[
    $capacity: ` + stc.Itoa(int(capacity)) + `
    $leasetime: ` + stc.Itoa(int(leasetime)) + `
](
    $type: <bali:/types/documents/Bag:v3>
    $tag: ` + fra.TagWithSize(20).AsString() + `
    $version: v1
    $permissions: ` + permissions.AsString() + `
)`
	var draft = not.DraftFromString(source)

	// Notarize the message bag.
	var bag = v.notary_.NotarizeDraft(draft)

	// Write the message bag out to document storage.
	var citation = v.storage_.WriteBag(bag)

	// Write the citation to the message bag out to document storage.
	v.storage_.WriteCitation(name, citation)
}

func (v *documentRepository_) RemoveBag(
	name fra.ResourceLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to delete a message bag.",
	)

	// Make sure the name of the message bag exists.
	var citation = v.storage_.ReadCitation(name)
	if uti.IsUndefined(citation) {
		var message = fmt.Sprintf(
			"Attempted to delete a non-existent message bag with name: %v",
			name,
		)
		panic(message)
	}

	// Delete the message bag (and any remaining messages) from document storage.
	v.storage_.DeleteBag(citation)

	// Delete the citation to the message bag from document storage.
	v.storage_.DeleteCitation(name)
}

func (v *documentRepository_) MessageCount(
	bag fra.ResourceLike,
) uti.Cardinal {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to calculate the message count.",
	)

	// Make sure the name of the message bag exists.
	var citation = v.storage_.ReadCitation(bag)
	if uti.IsUndefined(citation) {
		var message = fmt.Sprintf(
			"Attempted to access a non-existent message bag with name: %v",
			bag,
		)
		panic(message)
	}

	// Count the messages currently available in the message bag.
	return v.storage_.MessageCount(citation)
}

func (v *documentRepository_) SendMessage(
	bag fra.ResourceLike,
	message doc.ItemsLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to send a message via a bag.",
	)

	// Create the new message.
	var source = doc.FormatComponent(message) + `(
    $type: <bali:/types/documents/Message:v3>
    $tag: ` + fra.TagWithSize(20).AsString() + `
    $version: v1
	$permissions: <bali:/permissions/public:v3>
)`
	var draft = not.DraftFromString(source)
	draft.SetObject(bag, fra.Symbol("bag"))

	// Make sure the name of the message bag exists.
	var citation = v.storage_.ReadCitation(bag)
	if uti.IsUndefined(citation) {
		var message = fmt.Sprintf(
			"Attempted to send a message to a non-existent bag with name: %v",
			bag,
		)
		panic(message)
	}

	// Notarize the message.
	var contract = v.notary_.NotarizeDraft(draft)

	// Write the message out to the message bag in document storage.
	v.storage_.WriteMessage(citation, contract)
}

func (v *documentRepository_) RetrieveMessage(
	bag fra.ResourceLike,
) not.Notarized {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to retrieve a message from a bag.",
	)

	// Make sure the name of the message bag exists.
	var citation = v.storage_.ReadCitation(bag)
	if uti.IsUndefined(citation) {
		var message = fmt.Sprintf(
			"Attempted to retrieve a message from a non-existent bag with name: %v",
			bag,
		)
		panic(message)
	}

	// Read a message from the message bag in document storage.
	return v.storage_.ReadMessage(citation)
}

func (v *documentRepository_) AcceptMessage(
	message not.Notarized,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to accept a processed message.",
	)

	// Extract the bag citation from the message.
	var content = message.GetContent()
	var bag = content.GetObject(fra.Symbol("bag"))
	var component = bag.GetComponent()
	var name = component.GetEntity().(fra.ResourceLike)
	var bagCitation = v.storage_.ReadCitation(name)

	// Generate a message citation for the message.
	var messageCitation = v.notary_.CiteDraft(content)

	// Delete the message from the message bag in document storage.
	v.storage_.DeleteMessage(bagCitation, messageCitation)
}

func (v *documentRepository_) RejectMessage(
	message not.Notarized,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to reject a retrieved  message.",
	)

	// Extract the bag citation from the message.
	var content = message.GetContent()
	var bag = content.GetObject(fra.Symbol("bag"))
	var component = bag.GetComponent()
	var name = component.GetEntity().(fra.ResourceLike)
	var bagCitation = v.storage_.ReadCitation(name)

	// Generate a message citation for the message.
	var messageCitation = v.notary_.CiteDraft(content)

	// Reset the message lease in the message bag in document storage.
	v.storage_.ReleaseMessage(bagCitation, messageCitation)
}

func (v *documentRepository_) PublishEvent(
	event doc.ItemsLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to publish an event.",
	)

	// Create the new event.
	var source = doc.FormatComponent(event) + `(
    $type: <bali:/types/documents/Event:v3>
    $tag: ` + fra.TagWithSize(20).AsString() + `
    $version: v1
	$permissions: <bali:/permissions/public:v3>
)`
	var draft = not.DraftFromString(source)

	// Notarize the event.
	var contract = v.notary_.NotarizeDraft(draft)

	// Write the event out to the notification queue in document storage.
	v.storage_.WriteEvent(contract)
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
