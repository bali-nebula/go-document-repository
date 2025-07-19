/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package module_test

import (
	not "github.com/bali-nebula/go-digital-notary/v3"
	doc "github.com/bali-nebula/go-document-notation/v3"
	rep "github.com/bali-nebula/go-document-repository/v3"
	fra "github.com/craterdog/go-component-framework/v7"
	uti "github.com/craterdog/go-missing-utilities/v7"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

const testDirectory = "./test/"

func TestLocalStorage(t *tes.T) {
	// Initialize the document repository.
	uti.RemakeDirectory(testDirectory)
	var ssm = not.Ssm(testDirectory)
	var hsm = ssm
	var notary = not.DigitalNotary(testDirectory, ssm, hsm)
	notary.ForgetKey()
	var certificate = notary.GenerateKey()
	var storage rep.Persistent = rep.LocalStorage(notary, testDirectory)
	storage = rep.ValidatedStorage(notary, storage)
	storage = rep.CachedStorage(storage)
	var repository = rep.DocumentRepository(notary, storage)

	// Save the certificate.
	var citation = repository.SaveCertificate(certificate)
	ass.Equal(t, citation.AsString(), notary.GetCitation().AsString())
	certificate = notary.RefreshKey()
	citation = repository.SaveCertificate(certificate)
	ass.Equal(t, citation.AsString(), notary.GetCitation().AsString())

	// Save a draft document.
	var component = doc.ParseSource("~pi").GetComponent()
	var type_ = fra.ResourceFromString("<bali:/examples/Angle:v1>")
	var tag = fra.TagWithSize(20)
	var version = fra.VersionFromString("v1.2.3")
	var permissions = fra.ResourceFromString("<bali:/permissions/Public:v1>")
	var previous not.CitationLike
	var draft = not.Draft(
		component,
		type_,
		tag,
		version,
		permissions,
		previous,
	)
	citation = repository.SaveDraft(draft)

	// Retrieve the draft document.
	var same = repository.RetrieveDraft(citation)
	ass.Equal(t, draft.AsString(), same.AsString())

	// Discard the draft document
	repository.DiscardDraft(citation)

	// Create a notarized contract document.
	var resource = "<bali:/test/Contract:v1.2.3>"
	var contract = repository.NotarizeDraft(resource, draft)
	var same2 = repository.RetrieveContract(resource)
	ass.Equal(t, contract.AsString(), same2.AsString())
	var same3 = repository.RetrieveContract(resource)
	ass.Equal(t, same2.AsString(), same3.AsString())

	// Checkout a new draft of the contract document.
	draft = repository.CheckoutDraft(resource, 2)
	ass.NotEqual(t, draft.AsString(), same.AsString())

	// Create a new message bag.
	var bag = "<bali:/test/Bag:v1>"
	repository.CreateBag(bag, permissions.AsString(), 8, 10)
	ass.Equal(t, 0, repository.MessageCount(bag))

	// Send a message to the bag.
	var content = doc.ParseSource(`"Hello World!"`)
	repository.SendMessage(bag, content)
	ass.Equal(t, 1, repository.MessageCount(bag))

	// Retrieve a message from the bag.
	contract = repository.RetrieveMessage(bag)
	ass.Equal(t, 0, repository.MessageCount(bag))

	// Reject the message.
	repository.RejectMessage(contract)
	ass.Equal(t, 1, repository.MessageCount(bag))

	// Process the message.
	contract = repository.RetrieveMessage(bag)
	repository.AcceptMessage(contract)
	ass.Equal(t, 0, repository.MessageCount(bag))

	// Delete the bag.
	repository.DeleteBag(bag)

	// Publish an event.
	var kind = "<bali:/events/Example:v3>"
	content = doc.ParseSource(`"Something Happened!"`)
	repository.PublishEvent(kind, content, permissions.AsString())
}
