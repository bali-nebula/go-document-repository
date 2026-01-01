/*
................................................................................
.    Copyright (c) 2009-2026 Crater Dog Technologies™.  All Rights Reserved.   .
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
	uti "github.com/craterdog/go-essential-utilities/v8"
	ass "github.com/stretchr/testify/assert"
	syn "sync"
	tes "testing"
)

const testDirectory = "./test/"

var identity = not.Identity(`[
    $algorithm: "ED25519"
    $key: '>
        uR+AQ8Gs45g9hHcPUWMu7VXzgadQdSubVnssoE16YrA
    <'
    $attributes: [
        $surname: "Norton"
        $birthname: "Derk David"
        $birthdate: <1966-04-04>
        $birthplace: "Boulder, Colorado, USA"
        $birthsex: $male
        $nationality: "USA"
        $address: ">
            123 Main Street
            Louisville, Colorado, 80027
        <"
        $mobile: "303-555-1212"
        $email: "craterdog@gmail.com"
        $mugshot: '>
            oVVGU2Wa/n+kdOHtZ8Zidq5jD9UZ3G60QOXMdAh6cqg
        <'
    ]
](
    $type: /bali/types/notary/Identity/v3
    $tag: #BDBK83JS4YDAZJKAT9D646Z3PAXY4SXJ
    $version: v1
    $permissions: /bali/permissions/Public/v3
    $previous: none
)`)

func TestLocalStorage(t *tes.T) {
	// Initialize the document repository.
	var group doc.Synchronized = new(syn.WaitGroup)
	uti.RemakeDirectory(testDirectory)
	uti.MakeDirectory(testDirectory + "hsmEd25519/")
	var ssm = not.SsmSha512()
	var device = "dummy"
	var secret = "#ACH22TPZL7QSSFFH6GGG8D21N3S6Y5RQ"
	var hsm = HsmEd25519TestClass().HsmEd25519(device, secret)
	var notary = not.DigitalNotary(ssm, hsm)
	notary.ForgetKey()
	var attributes = identity.GetAttributes()
	var certificate = notary.GenerateKey(attributes)
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
	_, status = repository.SaveCertificate(certificate)
	ass.Equal(t, rep.Existed, status)
	certificate = notary.RefreshKey()
	citation, status = repository.SaveCertificate(certificate)
	ass.Equal(t, rep.Success, status)
	ass.True(t, notary.CitationMatches(citation, certificate))

	// Save a draft document.
	var entity any = doc.Angle("~π")
	var type_ = doc.Name("/bali/examples/Pi/v1")
	var tag = doc.Tag()
	var version = doc.Version("v1.2.3")
	var permissions = doc.Name("/bali/permissions/Public/v1")
	var previous doc.ResourceLike
	var content = not.Content(
		entity,
		type_,
		tag,
		version,
		permissions,
		previous,
	)
	var document = not.Document(content)
	_, status = repository.SaveDraft(document)
	ass.Equal(t, rep.Success, status)
	ass.False(t, document.IsNotarized())
	citation, status = repository.SaveDraft(document)
	ass.Equal(t, rep.Success, status)
	ass.True(t, notary.CitationMatches(citation, document))

	// Retrieve the draft document.
	var same not.DocumentLike
	same, status = repository.RetrieveDraft(citation)
	ass.Equal(t, rep.Success, status)
	ass.Equal(t, document.AsSource(), same.AsSource())

	// Create a notarized document.
	var name = doc.Name("/examples/documents/transaction/v1.2.3")
	status = repository.NotarizeDocument(name, document)
	ass.Equal(t, rep.Success, status)
	ass.True(t, document.IsNotarized())
	same, status = repository.RetrieveDocument(name)
	ass.Equal(t, rep.Success, status)
	ass.Equal(t, document.AsSource(), same.AsSource())
	_, status = repository.RetrieveDraft(citation)
	ass.Equal(t, rep.Missing, status)

	// Checkout a new draft of the document.
	document, status = repository.CheckoutDocument(name, uint(2))
	ass.Equal(t, rep.Success, status)
	ass.False(t, document.IsNotarized())
	ass.NotEqual(t, document.AsSource(), same.AsSource())
	document, status = repository.CheckoutDocument(name, uint(2))
	ass.Equal(t, rep.Success, status)
	ass.False(t, document.IsNotarized())
	ass.NotEqual(t, document.AsSource(), same.AsSource())

	// Save the new draft document.
	citation, status = repository.SaveDraft(document)
	ass.Equal(t, rep.Success, status)
	ass.False(t, document.IsNotarized())
	ass.True(t, notary.CitationMatches(citation, document))

	// Discard the draft document
	citation = notary.CiteDocument(document)
	same, status = repository.DiscardDraft(citation)
	ass.Equal(t, rep.Success, status)
	ass.Equal(t, document.AsSource(), same.AsSource())
	_, status = repository.RetrieveDraft(citation)
	ass.Equal(t, rep.Missing, status)

	// Send a message to a bag.
	var bag = doc.Name("/examples/bag")
	entity = doc.Quote("Hello World!")
	type_ = doc.Name("/bali/examples/Message/v1")
	tag = doc.Tag()
	version = doc.Version()
	permissions = doc.Name("/bali/permissions/Public/v1")
	content = not.Content(
		entity,
		type_,
		tag,
		version,
		permissions,
		previous,
	)
	var message = not.Document(content)
	ass.False(t, message.IsNotarized())
	status = repository.SendMessage(bag, message)
	ass.Equal(t, rep.Success, status)
	ass.True(t, message.IsNotarized())

	// Send another message to a bag.
	entity = doc.Quote("Hello Again...")
	tag = doc.Tag()
	content = not.Content(
		entity,
		type_,
		tag,
		version,
		permissions,
		previous,
	)
	message = not.Document(content)
	ass.False(t, message.IsNotarized())
	status = repository.SendMessage(bag, message)
	ass.Equal(t, rep.Success, status)
	ass.True(t, message.IsNotarized())

	// Subscribe to events.
	type_ = doc.Name("/bali/examples/Event/v1")
	status = repository.SubscribeEvents(bag, type_)
	ass.Equal(t, rep.Success, status)

	// Publish an event.
	entity = doc.Quote("Something happened...")
	tag = doc.Tag()
	content = not.Content(
		entity,
		type_,
		tag,
		version,
		permissions,
		previous,
	)
	var event = not.Document(content)
	ass.False(t, event.IsNotarized())
	status = repository.PublishEvent(event)
	ass.Equal(t, rep.Success, status)
	ass.True(t, event.IsNotarized())
	group.Wait()

	// Unsubscribe from events.
	status = repository.UnsubscribeEvents(bag, type_)
	ass.Equal(t, rep.Success, status)

	// Retrieve a message from the bag.
	message, status = repository.ReceiveMessage(bag)
	ass.Equal(t, rep.Success, status)
	ass.True(t, message.IsNotarized())

	// Reject the message.
	status = repository.RejectMessage(bag, message)
	ass.Equal(t, rep.Success, status)

	// Process the message.
	message, status = repository.ReceiveMessage(bag)
	ass.Equal(t, rep.Success, status)
	ass.True(t, message.IsNotarized())

	// Accept the message.
	status = repository.AcceptMessage(bag, message)
	ass.Equal(t, rep.Success, status)
}
