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
	var citation = notary.CiteDocument(certificate)
	ass.True(t, notary.CitationMatches(citation, certificate))
	citation, status = repository.SaveCertificate(certificate)
	ass.Equal(t, rep.Success, status)
	ass.True(t, notary.CitationMatches(citation, certificate))
	certificate = notary.RefreshKey()
	citation = notary.CiteDocument(certificate)
	ass.True(t, notary.CitationMatches(citation, certificate))
	citation, status = repository.SaveCertificate(certificate)
	ass.Equal(t, rep.Success, status)
	ass.True(t, notary.CitationMatches(citation, certificate))

	// Save a draft document.
	var angle = doc.Angle("~π")
	var component = doc.Component(angle, nil)
	var type_ = doc.Resource("<bali:/examples/Angle:v1>")
	var tag = doc.Tag()
	var version = doc.Version("v1.2.3")
	var previous doc.ResourceLike
	var permissions = doc.Resource("<bali:/permissions/Public:v1>")
	var content = not.Content(
		component,
		type_,
		tag,
		version,
		previous,
		permissions,
	)
	var document = not.Document(content)
	citation, status = repository.SaveDraft(document)
	ass.Equal(t, rep.Success, status)
	ass.True(t, notary.CitationMatches(citation, document))

	// Retrieve the draft document.
	var same not.DocumentLike
	same, status = repository.RetrieveDraft(citation)
	ass.Equal(t, rep.Success, status)
	ass.Equal(t, document.AsString(), same.AsString())

	// Create a notarized document.
	var name = doc.Name("/examples/Document")
	status = repository.NotarizeDocument(name, version, document)
	ass.Equal(t, rep.Success, status)
	same, status = repository.RetrieveDocument(name, version)
	ass.Equal(t, rep.Success, status)
	ass.Equal(t, document.AsString(), same.AsString())

	// Checkout a new draft of the contract document.
	document, status = repository.CheckoutDocument(name, version, uint(2))
	ass.Equal(t, rep.Success, status)
	ass.NotEqual(t, document.AsString(), same.AsString())

	// Discard the draft document
	citation = notary.CiteDocument(document)
	same, status = repository.DiscardDraft(citation)
	ass.Equal(t, rep.Success, status)
	ass.Equal(t, document.AsString(), same.AsString())

	// Send a message to a bag.
	var bag = doc.Name("/examples/bag")
	var quote = doc.Quote("Hello World!")
	component = doc.Component(quote, nil)
	type_ = doc.Resource("<bali:/examples/Message:v1>")
	tag = doc.Tag()
	version = doc.Version()
	permissions = doc.Resource("<bali:/permissions/Public:v1>")
	content = not.Content(
		component,
		type_,
		tag,
		version,
		previous,
		permissions,
	)

	var message = not.Document(content)
	status = repository.PostMessage(bag, message)
	ass.Equal(t, rep.Success, status)

	// Retrieve a message from the bag.
	message, status = repository.RetrieveMessage(bag)
	ass.Equal(t, rep.Success, status)

	// Reject the message.
	status = repository.RejectMessage(bag, message)
	ass.Equal(t, rep.Success, status)

	// Process the message.
	message, status = repository.RetrieveMessage(bag)
	ass.Equal(t, rep.Success, status)

	// Accept the message.
	status = repository.AcceptMessage(bag, message)
	ass.Equal(t, rep.Success, status)
}
