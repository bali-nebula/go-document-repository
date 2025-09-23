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
	fmt "fmt"
	doc "github.com/bali-nebula/go-bali-documents/v3"
	not "github.com/bali-nebula/go-digital-notary/v3"
	rep "github.com/bali-nebula/go-document-repository/v3/repository"
	uti "github.com/craterdog/go-missing-utilities/v7"
)

// CLASS INTERFACE

// Access Function

func SecureStorageClass() SecureStorageClassLike {
	return secureStorageClass()
}

// Constructor Methods

func (c *secureStorageClass_) SecureStorage(
	notary not.DigitalNotaryLike,
	storage rep.Persistent,
) SecureStorageLike {
	if uti.IsUndefined(notary) {
		panic("The \"notary\" attribute is required by this class.")
	}
	if uti.IsUndefined(storage) {
		panic("The \"storage\" attribute is required by this class.")
	}
	var instance = &secureStorage_{
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

func (v *secureStorage_) GetClass() SecureStorageClassLike {
	return secureStorageClass()
}

// Attribute Methods

// Persistent Methods

func (v *secureStorage_) ReadCitation(
	name doc.NameLike,
	version doc.VersionLike,
) doc.ResourceLike {
	return v.storage_.ReadCitation(name, version)
}

func (v *secureStorage_) WriteCitation(
	name doc.NameLike,
	version doc.VersionLike,
	citation doc.ResourceLike,
) {
	v.storage_.WriteCitation(name, version, citation)
}

func (v *secureStorage_) DeleteCitation(
	name doc.NameLike,
	version doc.VersionLike,
) {
	v.storage_.DeleteCitation(name, version)
}

func (v *secureStorage_) ListCitations(
	path doc.NameLike,
) doc.Sequential[doc.ResourceLike] {
	return v.storage_.ListCitations(path)
}

func (v *secureStorage_) ReadDraft(
	citation doc.ResourceLike,
) not.Parameterized {
	var draft = v.storage_.ReadDraft(citation)
	if !v.notary_.CitationMatches(citation, draft) {
		var message = fmt.Sprintf(
			"The citation does not match the cited draft document: %s%s",
			citation.AsString(),
			draft.AsString(),
		)
		panic(message)
	}
	return draft
}

func (v *secureStorage_) WriteDraft(
	draft not.Parameterized,
) doc.ResourceLike {
	return v.storage_.WriteDraft(draft)
}

func (v *secureStorage_) DeleteDraft(
	citation doc.ResourceLike,
) {
	v.storage_.DeleteDraft(citation)
}

func (v *secureStorage_) ReadContract(
	citation doc.ResourceLike,
) not.ContractLike {
	var document = v.storage_.ReadContract(citation)
	v.validateContract(document)
	return document
}

func (v *secureStorage_) WriteContract(
	contract not.ContractLike,
) doc.ResourceLike {
	v.validateContract(contract)
	return v.storage_.WriteContract(contract)
}

func (v *secureStorage_) DeleteContract(
	citation doc.ResourceLike,
) {
	v.storage_.DeleteContract(citation)
}

// PROTECTED INTERFACE

// Private Methods

func (v *secureStorage_) validateContract(
	contract not.ContractLike,
) {
	// Retrieve the citation to the certificate that signed the document.
	var notary = contract.GetNotary()
	var draft = contract.GetContent()
	if !v.notary_.CitationMatches(notary, draft) {
		// Not self-signed, read the certificate that signed the document.
		draft = v.storage_.ReadContract(notary).GetContent()
	}
	var certificate = not.Certificate(draft.AsString())
	if !v.notary_.SealMatches(contract, certificate) {
		var message = fmt.Sprintf(
			"The cited certificate was not used to notarize the document: %s%s",
			contract.AsString(),
			certificate.AsString(),
		)
		panic(message)
	}
}

// Instance Structure

type secureStorage_ struct {
	// Declare the instance attributes.
	notary_  not.DigitalNotaryLike
	storage_ rep.Persistent
}

// Class Structure

type secureStorageClass_ struct {
	// Declare the class constants.
}

// Class Reference

func secureStorageClass() *secureStorageClass_ {
	return secureStorageClassReference_
}

var secureStorageClassReference_ = &secureStorageClass_{
	// Initialize the class constants.
}
