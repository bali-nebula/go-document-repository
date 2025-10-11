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

package storage

import (
	doc "github.com/bali-nebula/go-bali-documents/v3"
	not "github.com/bali-nebula/go-digital-notary/v3"
	rep "github.com/bali-nebula/go-document-repository/v3/repository"
	uti "github.com/craterdog/go-missing-utilities/v7"
	log "log"
)

// CLASS INTERFACE

// Access Function

func ValidatedStorageClass() ValidatedStorageClassLike {
	return validatedStorageClass()
}

// Constructor Methods

func (c *validatedStorageClass_) ValidatedStorage(
	notary not.DigitalNotaryLike,
	storage rep.Persistent,
) ValidatedStorageLike {
	if uti.IsUndefined(notary) {
		panic("The \"notary\" attribute is required by this class.")
	}
	if uti.IsUndefined(storage) {
		panic("The \"storage\" attribute is required by this class.")
	}
	var instance = &validatedStorage_{
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

func (v *validatedStorage_) GetClass() ValidatedStorageClassLike {
	return validatedStorageClass()
}

// Attribute Methods

// Persistent Methods

func (v *validatedStorage_) WriteCitation(
	name doc.NameLike,
	citation not.CitationLike,
) (
	status rep.Status,
) {
	switch {
	case v.invalidCitation(citation):
		status = rep.Invalid
	default:
		status = v.storage_.WriteCitation(name, citation)
	}
	return
}

func (v *validatedStorage_) ReadCitation(
	name doc.NameLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	citation, status = v.storage_.ReadCitation(name)
	switch {
	case status != rep.Success:
	case v.invalidCitation(citation):
		status = rep.Invalid
	}
	return
}

func (v *validatedStorage_) DeleteCitation(
	name doc.NameLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	citation, status = v.storage_.DeleteCitation(name)
	return
}

func (v *validatedStorage_) WriteMessage(
	bag doc.NameLike,
	message not.DocumentLike,
) (
	status rep.Status,
) {
	switch {
	case v.invalidDocument(message):
		status = rep.Invalid
	default:
		status = v.storage_.WriteMessage(bag, message)
	}
	return
}

func (v *validatedStorage_) ReadMessage(
	bag doc.NameLike,
) (
	message not.DocumentLike,
	status rep.Status,
) {
	message, status = v.storage_.ReadMessage(bag)
	switch {
	case status != rep.Success:
	case v.invalidDocument(message):
		status = rep.Invalid
	}
	return
}

func (v *validatedStorage_) UnreadMessage(
	bag doc.NameLike,
	message not.DocumentLike,
) (
	status rep.Status,
) {
	switch {
	case v.invalidDocument(message):
		status = rep.Invalid
	default:
		status = v.storage_.UnreadMessage(bag, message)
	}
	return
}

func (v *validatedStorage_) DeleteMessage(
	bag doc.NameLike,
	message not.DocumentLike,
) (
	status rep.Status,
) {
	switch {
	case v.invalidDocument(message):
		status = rep.Invalid
	default:
		status = v.storage_.DeleteMessage(bag, message)
	}
	return
}

func (v *validatedStorage_) WriteSubscription(
	bag doc.NameLike,
	type_ doc.NameLike,
) (
	status rep.Status,
) {
	status = v.storage_.WriteSubscription(bag, type_)
	return
}

func (v *validatedStorage_) ReadSubscriptions(
	type_ doc.NameLike,
) (
	bags doc.Sequential[doc.NameLike],
	status rep.Status,
) {
	bags, status = v.storage_.ReadSubscriptions(type_)
	return
}

func (v *validatedStorage_) DeleteSubscription(
	bag doc.NameLike,
	type_ doc.NameLike,
) (
	status rep.Status,
) {
	status = v.storage_.DeleteSubscription(bag, type_)
	return
}

func (v *validatedStorage_) WriteDraft(
	draft not.DocumentLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	var content = draft.GetContent()
	switch {
	case draft.HasSeal():
		status = rep.Invalid
	case v.invalidContent(content):
		status = rep.Invalid
	default:
		citation, status = v.storage_.WriteDraft(draft)
	}
	return
}

func (v *validatedStorage_) ReadDraft(
	citation not.CitationLike,
) (
	draft not.DocumentLike,
	status rep.Status,
) {
	draft, status = v.storage_.ReadDraft(citation)
	switch {
	case status != rep.Success:
	case v.invalidContent(draft.GetContent()):
		status = rep.Invalid
	}
	return
}

func (v *validatedStorage_) DeleteDraft(
	citation not.CitationLike,
) (
	draft not.DocumentLike,
	status rep.Status,
) {
	draft, status = v.storage_.DeleteDraft(citation)
	return
}

func (v *validatedStorage_) WriteDocument(
	document not.DocumentLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	switch {
	case v.invalidContent(document.GetContent()):
		status = rep.Invalid
	case v.invalidDocument(document):
		status = rep.Invalid
	default:
		citation, status = v.storage_.WriteDocument(document)
	}
	return
}

func (v *validatedStorage_) ReadDocument(
	citation not.CitationLike,
) (
	document not.DocumentLike,
	status rep.Status,
) {
	document, status = v.storage_.ReadDocument(citation)
	switch {
	case status != rep.Success:
	case v.invalidDocument(document):
		status = rep.Invalid
	}
	return
}

func (v *validatedStorage_) DeleteDocument(
	citation not.CitationLike,
) (
	document not.DocumentLike,
	status rep.Status,
) {
	document, status = v.storage_.DeleteDocument(citation)
	return
}

// PROTECTED INTERFACE

// Private Methods

func (v *validatedStorage_) invalidCitation(
	citation not.CitationLike,
) bool {
	// Validate that the citation refers to a document.
	var document, status = v.storage_.ReadDocument(citation)
	if status != rep.Success {
		log.Printf(
			"ValidatedStorage: "+
				"The following citation doesn't cite a document: %s\n",
			citation,
		)
		return true
	}

	// Validate the digest of the cited document.
	var doesNotMatch = !v.notary_.CitationMatches(citation, document)
	if doesNotMatch {
		log.Printf(
			"ValidataedStorage: "+
				"The following digest doesn't match the cited document: %s %s\n",
			citation,
			document.AsString(),
		)
	}
	return doesNotMatch
}

func (v *validatedStorage_) invalidContent(
	content not.Parameterized,
) bool {
	// TBD - Validate the citations to the type and permissions documents.

	// Validate the citation to the previous version of the document.
	var previous = content.GetOptionalPrevious()
	if uti.IsDefined(previous) {
		var citation = not.Citation(previous)
		if v.invalidCitation(citation) {
			log.Printf(
				"ValidataedStorage: "+
					"The previous citation doesn't cite an existing document: %s\n",
				content.AsString(),
			)
			return true
		}
	}
	return false
}

func (v *validatedStorage_) invalidDocument(
	document not.DocumentLike,
) bool {
	// Validate the content of the document.
	var content = document.GetContent()
	if v.invalidContent(content) {
		log.Printf(
			"ValidataedStorage: "+
				"The content for the following document isn't valid: %s\n",
			document.AsString(),
		)
		return true
	}

	// Validate the signature on the notarized document.
	if !document.HasSeal() {
		log.Printf(
			"ValidataedStorage: "+
				"The following document is missing a notary seal: %s\n",
			document.AsString(),
		)
		return true
	}
	var notary = document.GetNotary()
	var certificate = document
	var status rep.Status
	if uti.IsDefined(notary) {
		// The document is not self-signed, so read the notary certificate.
		certificate, status = v.storage_.ReadDocument(notary)
		if status != rep.Success {
			log.Printf(
				"ValidatedStorage: "+
					"The cited notary certificate doesn't exist: %s\n",
				notary,
			)
			return true
		}
	}
	var doesNotMatch = !v.notary_.SealMatches(document, certificate)
	if doesNotMatch {
		log.Printf(
			"ValidatedStorage: "+
				"The notary seal doesn't match the notarized document: %s %s\n",
			document.AsString(),
			certificate,
		)
		return true
	}
	return false
}

// Instance Structure

type validatedStorage_ struct {
	// Declare the instance attributes.
	notary_  not.DigitalNotaryLike
	storage_ rep.Persistent
}

// Class Structure

type validatedStorageClass_ struct {
	// Declare the class constants.
}

// Class Reference

func validatedStorageClass() *validatedStorageClass_ {
	return validatedStorageClassReference_
}

var validatedStorageClassReference_ = &validatedStorageClass_{
	// Initialize the class constants.
}
