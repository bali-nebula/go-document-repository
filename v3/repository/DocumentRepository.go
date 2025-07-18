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
	doc "github.com/bali-nebula/go-document-notation/v3"
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

func (v *documentRepository_) SaveDraft(
	draft not.DraftLike,
) not.CitationLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to save a draft document.",
	)

	// Write the draft document out to document storage.
	var citation = v.storage_.WriteDraft(draft)
	return citation
}

func (v *documentRepository_) RetrieveDraft(
	draft not.CitationLike,
) not.DraftLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to retrieve a draft document.",
	)

	// Read the draft document from document storage.
	return v.storage_.ReadDraft(draft)
}

func (v *documentRepository_) DiscardDraft(
	draft not.CitationLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to discard a draft document.",
	)

	// Remove the draft document from document storage.
	v.storage_.RemoveDraft(draft)
}

func (v *documentRepository_) NotarizeDraft(
	name string,
	draft not.DraftLike,
) not.ContractLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to notarize a draft document.",
	)

	// Make sure the name of the document is unique.
	var resource = fra.ResourceFromString(name)
	if v.storage_.CitationExists(resource) {
		var message = fmt.Sprintf(
			"Attempted to notarize a draft document using an existing name: %v",
			name,
		)
		panic(message)
	}

	// Notarize the draft document.
	var contract = v.notary_.NotarizeDraft(draft)

	// Write the notarized contract out to document storage.
	var citation = v.storage_.WriteContract(contract)

	// Write the citation to the contract out to document storage.
	v.storage_.WriteCitation(resource, citation)
	return contract
}

func (v *documentRepository_) RetrieveContract(
	contract string,
) not.ContractLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to notarize a contract document.",
	)

	// Make sure the name of the contract document exists.
	var resource = fra.ResourceFromString(contract)
	var citation = v.storage_.ReadCitation(resource)
	if uti.IsUndefined(citation) {
		var message = fmt.Sprintf(
			"Attempted to retrieve a non-existent contract with name: %v",
			contract,
		)
		panic(message)
	}

	// Read the contract document from document storage.
	return v.storage_.ReadContract(citation)
}

func (v *documentRepository_) CheckoutDraft(
	contract string,
	level int,
) not.DraftLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to checkout a draft document.",
	)

	// Make sure the name of the contract document exists.
	var resource = fra.ResourceFromString(contract)
	var citation = v.storage_.ReadCitation(resource)
	if uti.IsUndefined(citation) {
		var message = fmt.Sprintf(
			"Attempted to checkout a non-existent contract with name: %v",
			contract,
		)
		panic(message)
	}

	// Read the contract document from document storage.
	var draft = v.storage_.ReadContract(citation).GetDraft()

	// Create a draft copy of the contract document with an updated version.
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
	name string,
	permissions string,
	capacity int,
	leasetime int,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to create a message bag.",
	)

	// Create a new message bag.
	var source = `[
    $capacity: ` + stc.Itoa(capacity) + `
    $leasetime: ` + stc.Itoa(leasetime) + `
](
    $type: <bali:/types/documents/Bag:v3>
    $tag: ` + fra.TagWithSize(20).AsString() + `
    $version: v1
    $permissions: ` + permissions + `
)`
	var draft = not.DraftFromString(source)

	// Notarize the message bag.
	var bag = v.notary_.NotarizeDraft(draft)

	// Write the message bag out to document storage.
	var citation = v.storage_.WriteBag(bag)

	// Write the citation to the message bag out to document storage.
	var resource = fra.ResourceFromString(name)
	v.storage_.WriteCitation(resource, citation)
}

func (v *documentRepository_) DeleteBag(
	bag string,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to delete a message bag.",
	)

	// Make sure the name of the message bag exists.
	var resource = fra.ResourceFromString(bag)
	var citation = v.storage_.ReadCitation(resource)
	if uti.IsUndefined(citation) {
		var message = fmt.Sprintf(
			"Attempted to delete a non-existent message bag with name: %v",
			bag,
		)
		panic(message)
	}

	// Remove the message bag (and any remaining messages) from document storage.
	v.storage_.RemoveBag(citation)
}

func (v *documentRepository_) MessageCount(
	bag string,
) int {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to calculate the message count.",
	)

	// Make sure the name of the message bag exists.
	var resource = fra.ResourceFromString(bag)
	var citation = v.storage_.ReadCitation(resource)
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
	bag string,
	content doc.DocumentLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to send a message via a bag.",
	)

	// Create the new message.
	var source = `[
    $bag: ` + bag + `
    $content: ` + doc.FormatDocument(content) + `](
    $type: <bali:/types/documents/Message:v3>
    $tag: ` + fra.TagWithSize(20).AsString() + `
    $version: v1
	$permissions: <bali:/permissions/Public:v3>
)`
	var draft = not.DraftFromString(source)

	// Make sure the name of the message bag exists.
	var resource = fra.ResourceFromString(bag)
	var citation = v.storage_.ReadCitation(resource)
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
	bag string,
) not.ContractLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to retrieve a message from a bag.",
	)

	// Make sure the name of the message bag exists.
	var resource = fra.ResourceFromString(bag)
	var citation = v.storage_.ReadCitation(resource)
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
	message not.ContractLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to accept a processed message.",
	)

	// Extract the bag citation from the message.
	var draft = message.GetDraft()
	var source = draft.AsString()
	var document = doc.ParseSource(source)
	var bag = v.extractBag(document)
	var bagCitation = v.storage_.ReadCitation(bag)

	// Generate a message citation for the message.
	var messageCitation = v.notary_.CiteDraft(draft)

	// Remove the message from the message bag in document storage.
	v.storage_.RemoveMessage(bagCitation, messageCitation)
}

func (v *documentRepository_) RejectMessage(
	message not.ContractLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to reject a retrieved  message.",
	)

	// Extract the bag citation from the message.
	var draft = message.GetDraft()
	var source = draft.AsString()
	var document = doc.ParseSource(source)
	var bag = v.extractBag(document)
	var bagCitation = v.storage_.ReadCitation(bag)

	// Generate a message citation for the message.
	var messageCitation = v.notary_.CiteDraft(draft)

	// Reset the message lease in the message bag in document storage.
	v.storage_.ReleaseMessage(bagCitation, messageCitation)
}

func (v *documentRepository_) PublishEvent(
	kind string,
	content doc.DocumentLike,
	permissions string,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to publish an event.",
	)

	// Create the new event.
	var source = `[
    $kind: ` + kind + `
    $content: ` + doc.FormatDocument(content) + `](
    $type: <bali:/types/documents/Event:v3>
    $tag: ` + fra.TagWithSize(20).AsString() + `
    $version: v1
	$permissions: ` + permissions + `
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

func (v *documentRepository_) extractBag(
	document doc.DocumentLike,
) fra.ResourceLike {
	var bag fra.ResourceLike
	var component = document.GetComponent()
	var collection = component.GetAny().(doc.CollectionLike)
	var attributes = collection.GetAny().(doc.AttributesLike)
	var associations = attributes.GetAssociations()
	var iterator = associations.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var element = association.GetPrimitive().GetAny().(doc.ElementLike)
		var symbol = element.GetAny().(string)
		if symbol == "$bag" {
			var source = doc.FormatDocument(association.GetDocument())
			bag = fra.ResourceFromString(source)
			break
		}
	}
	return bag
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
