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
	version doc.VersionLike,
	citation not.CitationLike,
) (
	status rep.Status,
) {
	status = v.storage_.WriteCitation(name, version, citation)
	return
}

func (v *validatedStorage_) ReadCitation(
	name doc.NameLike,
	version doc.VersionLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	citation, status = v.storage_.ReadCitation(name, version)
	if status != rep.Success {
		return
	}
	if v.invalidCitation(citation) {
		status = rep.Invalid
	}
	return
}

func (v *validatedStorage_) DeleteCitation(
	name doc.NameLike,
	version doc.VersionLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	citation, status = v.storage_.DeleteCitation(name, version)
	return
}

func (v *validatedStorage_) WriteMessage(
	bag doc.NameLike,
	message not.CitationLike,
) (
	status rep.Status,
) {
	status = v.storage_.WriteMessage(bag, message)
	return
}

func (v *validatedStorage_) ReadMessage(
	bag doc.NameLike,
) (
	message not.CitationLike,
	status rep.Status,
) {
	message, status = v.storage_.ReadMessage(bag)
	if status != rep.Success {
		return
	}
	if v.invalidCitation(message) {
		status = rep.Invalid
	}
	return
}

func (v *validatedStorage_) UnreadMessage(
	bag doc.NameLike,
	message not.CitationLike,
) (
	status rep.Status,
) {
	if v.invalidCitation(message) {
		status = rep.Invalid
		return
	}
	status = v.storage_.UnreadMessage(bag, message)
	return
}

func (v *validatedStorage_) DeleteMessage(
	bag doc.NameLike,
	message not.CitationLike,
) (
	status rep.Status,
) {
	status = v.storage_.DeleteMessage(bag, message)
	return
}

func (v *validatedStorage_) WriteDraft(
	draft not.DocumentLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	if v.invalidDocument(draft) {
		status = rep.Invalid
		return
	}
	citation, status = v.storage_.WriteDraft(draft)
	return
}

func (v *validatedStorage_) ReadDraft(
	citation not.CitationLike,
) (
	draft not.DocumentLike,
	status rep.Status,
) {
	draft, status = v.storage_.ReadDraft(citation)
	if v.invalidDocument(draft) {
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
	if v.invalidDocument(document) {
		status = rep.Invalid
		return
	}
	citation, status = v.storage_.WriteDocument(document)
	return
}

func (v *validatedStorage_) ReadDocument(
	citation not.CitationLike,
) (
	document not.DocumentLike,
	status rep.Status,
) {
	document, status = v.storage_.ReadDocument(citation)
	if v.invalidDocument(document) {
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
	var document, status = v.storage_.ReadDocument(citation)
	if status != rep.Success {
		log.Printf("The citation does not cite a document: %s\n", citation)
		return true
	}
	var matches = v.notary_.CitationMatches(citation, document)
	if !matches {
		log.Printf(
			"The citation digest does not match the document: %s %s\n",
			citation,
			document.AsString(),
		)
	}
	return !matches
}

func (v *validatedStorage_) invalidDocument(
	document not.DocumentLike,
) bool {
	if !document.HasSeal() {
		return false
	}
	// Retrieve the citation to the certificate that signed the document.
	var notary = document.GetNotary()
	var certificate = document
	var status rep.Status
	if uti.IsDefined(notary) {
		// The document is not self-signed, so read the notary certificate.
		certificate, status = v.storage_.ReadDocument(notary)
		if status != rep.Success {
			log.Printf(
				"ValidatedStorage: The cited notary certificate does not exist: %s\n",
				notary,
			)
			return true
		}
	}
	return !v.notary_.SealMatches(document, certificate)
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
