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
	doc "github.com/bali-nebula/go-bali-documents/v3"
	not "github.com/bali-nebula/go-digital-notary/v3"
	rep "github.com/bali-nebula/go-document-repository/v3"
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
	var status rep.Status
	var content = certificate.GetContent()
	var citation = notary.CiteDocument(content)
	ass.True(t, notary.CitationMatches(citation, content))
	citation, status = repository.SaveCertificate(certificate)
	ass.Equal(t, rep.Written, status)
	ass.True(t, notary.CitationMatches(citation, content))
	certificate = notary.RefreshKey()
	content = certificate.GetContent()
	citation = notary.CiteDocument(content)
	ass.True(t, notary.CitationMatches(citation, content))
	citation, status = repository.SaveCertificate(certificate)
	ass.Equal(t, rep.Written, status)
	ass.True(t, notary.CitationMatches(citation, content))

	// Save a draft document.
	var component = doc.Angle("~π")
	var type_ = doc.Resource("<bali:/examples/Angle:v1>")
	var tag = doc.Tag()
	var version = doc.Version("v1.2.3")
	var permissions = doc.Resource("<bali:/permissions/Public:v1>")
	var previous doc.ResourceLike
	var draft not.Parameterized = not.Draft(
		component,
		type_,
		tag,
		version,
		permissions,
		previous,
	)
	citation, status = repository.SaveDraft(draft)
	ass.Equal(t, rep.Written, status)

	// Retrieve the draft document.
	var same not.Parameterized
	same, status = repository.RetrieveDraft(citation)
	ass.Equal(t, rep.Retrieved, status)
	ass.Equal(t, draft.AsString(), same.AsString())

	// Create a notarized contract document.
	var name = doc.Name("/examples/Contract")
	var contract, same2, same3 not.ContractLike
	contract, status = repository.NotarizeDocument(name, version, draft)
	ass.Equal(t, rep.Written, status)
	same2, status = repository.RetrieveDocument(name, version)
	ass.Equal(t, rep.Retrieved, status)
	ass.Equal(t, contract.AsString(), same2.AsString())
	same3, status = repository.RetrieveDocument(name, version)
	ass.Equal(t, rep.Retrieved, status)
	ass.Equal(t, same2.AsString(), same3.AsString())

	// Checkout a new draft of the contract document.
	draft, status = repository.CheckoutDocument(name, version, uint(2))
	ass.Equal(t, rep.Retrieved, status)
	ass.NotEqual(t, draft.AsString(), same.AsString())

	// Discard the draft document
	status = repository.DiscardDraft(citation)
	ass.Equal(t, rep.Deleted, status)

	// Create a new message bag.
	var bag = doc.Name("/examples/Bag")
	status = repository.CreateBag(bag, 8, 10, permissions)
	ass.Equal(t, rep.Written, status)

	// Send a message to the bag.
	var entity = doc.ParseSource(`[
    $quote: "Hello World!"
]`).GetEntity()
	var message = rep.Message(
		entity,
		doc.Resource("<bali:/examples/Message:v1>"),
		permissions,
	)
	status = repository.PostMessage(bag, message)
	ass.Equal(t, rep.Written, status)

	// Retrieve a message from the bag.
	contract, status = repository.RetrieveMessage(bag)
	ass.Equal(t, rep.Retrieved, status)

	// Reject the message.
	status = repository.RejectMessage(contract)
	ass.Equal(t, rep.Deleted, status)

	// Process the message.
	contract, status = repository.RetrieveMessage(bag)
	ass.Equal(t, rep.Retrieved, status)
	status = repository.AcceptMessage(contract)
	ass.Equal(t, rep.Deleted, status)

	// Remove the bag.
	status = repository.RemoveBag(bag)
	ass.Equal(t, rep.Deleted, status)
}
