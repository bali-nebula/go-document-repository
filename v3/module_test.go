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

package module_test

import (
	not "github.com/bali-nebula/go-digital-notary/v3"
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
	storage = rep.SecureStorage(notary, storage)
	storage = rep.CachedStorage(storage)
	var repository = rep.DocumentRepository(notary, storage)

	// Save the certificate.
	var citation = repository.SaveCertificate(certificate)
	ass.Equal(t, citation.AsString(), notary.GetCitation().AsString())
	certificate = notary.RefreshKey()
	citation = repository.SaveCertificate(certificate)
	ass.Equal(t, citation.AsString(), notary.GetCitation().AsString())

	// Save a draft document.
	var component = fra.AngleFromString("~pi")
	var type_ = fra.ResourceFromString("<bali:/examples/Angle:v1>")
	var tag = fra.TagWithSize(20)
	var version = fra.VersionFromString("v1.2.3")
	var permissions = fra.ResourceFromString("<bali:/permissions/Public:v1>")
	var previous fra.ResourceLike
	var draft not.Parameterized = not.Draft(
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
	var name = fra.NameFromString("/examples/Contract")
	var contract = repository.NotarizeDocument(name, version, draft)
	var same2 = repository.RetrieveDocument(name, version)
	ass.Equal(t, contract.AsString(), same2.AsString())
	var same3 = repository.RetrieveDocument(name, version)
	ass.Equal(t, same2.AsString(), same3.AsString())

	// Checkout a new draft of the contract document.
	draft = repository.CheckoutDocument(name, version, uti.Cardinal(2))
	ass.NotEqual(t, draft.AsString(), same.AsString())

	// Create a new message bag.
	var bag = fra.NameFromString("/examples/Bag")
	repository.CreateBag(bag, 8, 10, permissions)

	// Send a message to the bag.
	var message = rep.Message(
		fra.QuoteFromString(`"Hello World!"`),
		fra.ResourceFromString("<bali:/examples/Message:v1>"),
		permissions,
	)
	repository.PostMessage(bag, message)

	// Retrieve a message from the bag.
	contract = repository.RetrieveMessage(bag)

	// Reject the message.
	repository.RejectMessage(contract)

	// Process the message.
	contract = repository.RetrieveMessage(bag)
	repository.AcceptMessage(contract)

	// Remove the bag.
	repository.RemoveBag(bag)
}
