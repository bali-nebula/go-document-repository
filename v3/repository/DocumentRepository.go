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

func (c *documentRepositoryClass_) ExtractBag(
	document doc.DocumentLike,
) fra.ResourceLike {
	var bag fra.ResourceLike
	var draft = not.DraftClass().ExtractDraft(document)
	var component = draft.GetComponent()
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

func (c *documentRepositoryClass_) ExtractContent(
	document doc.DocumentLike,
) doc.ComponentLike {
	var content doc.ComponentLike
	var component = document.GetComponent()
	var collection = component.GetAny().(doc.CollectionLike)
	var attributes = collection.GetAny().(doc.AttributesLike)
	var associations = attributes.GetAssociations()
	var iterator = associations.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var element = association.GetPrimitive().GetAny().(doc.ElementLike)
		var symbol = element.GetAny().(string)
		if symbol == "$content" {
			content = association.GetDocument().GetComponent()
			break
		}
	}
	return content
}

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
	draft not.CitationLike,
) not.DraftLike {
	return v.storage_.ReadDraft(draft)
}

func (v *documentRepository_) DiscardDraft(
	draft not.CitationLike,
) {
	v.storage_.RemoveDraft(draft)
}

func (v *documentRepository_) NotarizeDraft(
	name fra.ResourceLike,
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
	contract fra.ResourceLike,
) not.ContractLike {
	var citation = v.storage_.ReadCitation(contract)
	if uti.IsUndefined(citation) {
		var message = fmt.Sprintf(
			"Attempted to retrieve a non-existent contract with name: %v",
			contract,
		)
		panic(message)
	}
	return v.storage_.ReadContract(citation)
}

func (v *documentRepository_) CheckoutDraft(
	contract fra.ResourceLike,
	level int,
) not.DraftLike {
	var citation = v.storage_.ReadCitation(contract)
	if uti.IsUndefined(citation) {
		var message = fmt.Sprintf(
			"Attempted to checkout a non-existent contract with name: %v",
			contract,
		)
		panic(message)
	}
	var draft = v.storage_.ReadContract(citation).GetDraft()
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
	name fra.ResourceLike,
	permissions fra.ResourceLike,
	capacity int,
	leasetime int,
) {
	var source = `[
    $capacity: ` + stc.Itoa(capacity) + `
    $leasetime: ` + stc.Itoa(leasetime) + `
](
    $type: <bali:/types/documents/Bag:v3>
    $tag: ` + fra.TagWithSize(20).AsString() + `
    $version: v1
    $permissions: ` + permissions.AsString() + `
)`
	var draft = not.DraftFromString(source)
	var bag = v.notary_.NotarizeDraft(draft)
	var citation = v.storage_.WriteBag(bag)
	v.storage_.WriteCitation(name, citation)
}

func (v *documentRepository_) DeleteBag(
	bag fra.ResourceLike,
) {
	var citation = v.storage_.ReadCitation(bag)
	v.storage_.RemoveBag(citation)
}

func (v *documentRepository_) MessageCount(
	bag fra.ResourceLike,
) int {
	var citation = v.storage_.ReadCitation(bag)
	return v.storage_.MessageCount(citation)
}

func (v *documentRepository_) SendMessage(
	bag fra.ResourceLike,
	content doc.DocumentLike,
) {
	var source = `[
    $bag: ` + bag.AsString() + `
    $content: ` + doc.FormatDocument(content) + `](
    $type: <bali:/types/documents/Message:v3>
    $tag: ` + fra.TagWithSize(20).AsString() + `
    $version: v1
	$permissions: <bali:/permissions/Public:v3>
)`
	var draft = not.DraftFromString(source)
	var citation = v.storage_.ReadCitation(bag)
	var contract = v.notary_.NotarizeDraft(draft)
	v.storage_.WriteMessage(citation, contract)
}

func (v *documentRepository_) RetrieveMessage(
	bag fra.ResourceLike,
) not.ContractLike {
	var citation = v.storage_.ReadCitation(bag)
	return v.storage_.ReadMessage(citation)
}

func (v *documentRepository_) AcceptMessage(
	message not.ContractLike,
) {
	var source = message.AsString()
	var document = doc.ParseSource(source)
	var bag = documentRepositoryClass().ExtractBag(document)
	var bagCitation = v.storage_.ReadCitation(bag)
	var draft = message.GetDraft()
	var messageCitation = v.notary_.CiteDraft(draft)
	v.storage_.RemoveMessage(bagCitation, messageCitation)
}

func (v *documentRepository_) RejectMessage(
	message not.ContractLike,
) {
	var source = message.AsString()
	var document = doc.ParseSource(source)
	var bag = documentRepositoryClass().ExtractBag(document)
	var bagCitation = v.storage_.ReadCitation(bag)
	var draft = message.GetDraft()
	var messageCitation = v.notary_.CiteDraft(draft)
	v.storage_.ReleaseMessage(bagCitation, messageCitation)
}

func (v *documentRepository_) PublishEvent(
	kind fra.ResourceLike,
	content doc.DocumentLike,
	permissions fra.ResourceLike,
) {
	var source = `[
    $kind: ` + kind.AsString() + `
    $content: ` + doc.FormatDocument(content) + `](
    $type: <bali:/types/documents/Event:v3>
    $tag: ` + fra.TagWithSize(20).AsString() + `
    $version: v1
	$permissions: ` + permissions.AsString() + `
)`
	var draft = not.DraftFromString(source)
	var contract = v.notary_.NotarizeDraft(draft)
	v.storage_.WriteEvent(contract)
}

// Attribute Methods

// PROTECTED INTERFACE

// Private Methods

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
