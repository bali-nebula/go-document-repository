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
	syn "sync"
	tes "testing"
)

const testDirectory = "./test/"

func TestLocalStorage(t *tes.T) {
	// Initialize the document repository.
	var group rep.Synchronized = new(syn.WaitGroup)
	uti.RemakeDirectory(testDirectory)
	var ssm = not.Ssm(testDirectory)
	var hsm = ssm
	var notary = not.DigitalNotary(testDirectory, ssm, hsm)
	notary.ForgetKey()
	var certificate = notary.GenerateKey()
	var storage rep.Persistent = rep.LocalStorage(notary, testDirectory)
	storage = rep.ValidatedStorage(notary, storage)
	storage = rep.CachedStorage(storage)
	var repository = rep.DocumentRepository(group, notary, storage)

	// Save the certificate.
	var status rep.Status
	var citation not.CitationLike
	citation, status = repository.SaveCertificate(certificate)
	ass.Equal(t, rep.Success, status)
	ass.True(t, notary.CitationMatches(citation, certificate))
	certificate = notary.RefreshKey()
	citation, status = repository.SaveCertificate(certificate)
	ass.Equal(t, rep.Success, status)
	ass.True(t, notary.CitationMatches(citation, certificate))

	// Save a draft document.
	var entity any = doc.Angle("~π")
	var type_ = doc.Resource("<bali:/examples/Pi:v1>")
	var tag = doc.Tag()
	var version = doc.Version("v1.2.3")
	var previous doc.ResourceLike
	var permissions = doc.Resource("<bali:/permissions/Public:v1>")
	var content = not.Content(
		entity,
		type_,
		tag,
		version,
		previous,
		permissions,
	)
	var document = not.Document(content)
	citation, status = repository.SaveDraft(document)
	ass.Equal(t, rep.Success, status)
	ass.False(t, document.HasSeal())
	ass.True(t, notary.CitationMatches(citation, document))

	// Retrieve the draft document.
	var same not.DocumentLike
	same, status = repository.RetrieveDraft(citation)
	ass.Equal(t, rep.Success, status)
	ass.Equal(t, document.AsString(), same.AsString())

	// Create a notarized document.
	var name = doc.Name("/examples/documents/transaction")
	status = repository.NotarizeDocument(name, version, document)
	ass.Equal(t, rep.Success, status)
	ass.True(t, document.HasSeal())
	same, status = repository.RetrieveDocument(name, version)
	ass.Equal(t, rep.Success, status)
	ass.Equal(t, document.AsString(), same.AsString())

	// Checkout a new draft of the document.
	document, status = repository.CheckoutDocument(name, version, uint(2))
	ass.Equal(t, rep.Success, status)
	ass.False(t, document.HasSeal())
	ass.NotEqual(t, document.AsString(), same.AsString())

	// Save the new draft document.
	citation, status = repository.SaveDraft(document)
	ass.Equal(t, rep.Success, status)
	ass.False(t, document.HasSeal())
	ass.True(t, notary.CitationMatches(citation, document))

	// Discard the draft document
	citation = notary.CiteDocument(document)
	same, status = repository.DiscardDraft(citation)
	ass.Equal(t, rep.Success, status)
	ass.Equal(t, document.AsString(), same.AsString())

	// Send a message to a bag.
	var bag = doc.Name("/examples/bag")
	entity = doc.Quote("Hello World!")
	type_ = doc.Resource("<bali:/examples/Message:v1>")
	tag = doc.Tag()
	version = doc.Version()
	permissions = doc.Resource("<bali:/permissions/Public:v1>")
	content = not.Content(
		entity,
		type_,
		tag,
		version,
		previous,
		permissions,
	)

	var message = not.Document(content)
	ass.False(t, message.HasSeal())
	status = repository.SendMessage(bag, message)
	ass.Equal(t, rep.Success, status)
	ass.True(t, message.HasSeal())

	// Retrieve a message from the bag.
	message, status = repository.RetrieveMessage(bag)
	ass.Equal(t, rep.Success, status)
	ass.True(t, message.HasSeal())

	// Reject the message.
	status = repository.RejectMessage(bag, message)
	ass.Equal(t, rep.Success, status)

	// Process the message.
	message, status = repository.RetrieveMessage(bag)
	ass.Equal(t, rep.Success, status)
	ass.True(t, message.HasSeal())

	// Accept the message.
	status = repository.AcceptMessage(bag, message)
	ass.Equal(t, rep.Success, status)

	// Subscribe to events.
	type_ = doc.Resource("<bali:/examples/Event:v1>")
	status = repository.SubscribeEvents(bag, type_)
	ass.Equal(t, rep.Success, status)

	// Publish an event.
	entity = doc.Quote("Something happened...")
	tag = doc.Tag()
	version = doc.Version()
	permissions = doc.Resource("<bali:/permissions/Public:v1>")
	content = not.Content(
		entity,
		type_,
		tag,
		version,
		previous,
		permissions,
	)
	var event = not.Document(content)
	ass.False(t, event.HasSeal())
	status = repository.PublishEvent(event)
	ass.Equal(t, rep.Success, status)
	ass.True(t, event.HasSeal())
	group.Wait()

	// Unsubscribe from events.
	status = repository.UnsubscribeEvents(bag, type_)
	ass.Equal(t, rep.Success, status)

	// Retrieve a message from the bag.
	message, status = repository.RetrieveMessage(bag)
	ass.Equal(t, rep.Success, status)
	ass.True(t, message.HasSeal())

	// Accept the message.
	status = repository.AcceptMessage(bag, message)
	ass.Equal(t, rep.Success, status)
}
