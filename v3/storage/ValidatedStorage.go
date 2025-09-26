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

func (v *validatedStorage_) ReadCitation(
	name doc.NameLike,
	version doc.VersionLike,
) (
	citation doc.ResourceLike,
	status rep.Status,
) {
	citation, status = v.storage_.ReadCitation(name, version)
	if status != rep.Retrieved {
		return
	}
	if v.invalidCitation(citation) {
		status = rep.Invalid
		return
	}
	return
}

func (v *validatedStorage_) WriteCitation(
	name doc.NameLike,
	version doc.VersionLike,
	citation doc.ResourceLike,
) (
	status rep.Status,
) {
	if v.invalidCitation(citation) {
		status = rep.Invalid
		return
	}
	status = v.storage_.WriteCitation(name, version, citation)
	return
}

func (v *validatedStorage_) DeleteCitation(
	name doc.NameLike,
	version doc.VersionLike,
) (
	status rep.Status,
) {
	status = v.storage_.DeleteCitation(name, version)
	return
}

func (v *validatedStorage_) ListCitations(
	path doc.NameLike,
) (
	citations doc.Sequential[doc.ResourceLike],
	status rep.Status,
) {
	citations, status = v.storage_.ListCitations(path)
	return
}

func (v *validatedStorage_) ReadDraft(
	citation doc.ResourceLike,
) (
	draft not.Parameterized,
	status rep.Status,
) {
	if v.invalidCitation(citation) {
		status = rep.Invalid
		return
	}
	draft, status = v.storage_.ReadDraft(citation)
	return
}

func (v *validatedStorage_) WriteDraft(
	draft not.Parameterized,
) (
	citation doc.ResourceLike,
	status rep.Status,
) {
	citation, status = v.storage_.WriteDraft(draft)
	if status != rep.Written {
		return
	}
	if v.invalidCitation(citation) {
		status = rep.Invalid
	}
	return
}

func (v *validatedStorage_) DeleteDraft(
	citation doc.ResourceLike,
) (
	status rep.Status,
) {
	if v.invalidCitation(citation) {
		status = rep.Invalid
		return
	}
	status = v.storage_.DeleteDraft(citation)
	return
}

func (v *validatedStorage_) ReadContract(
	citation doc.ResourceLike,
) (
	contract not.ContractLike,
	status rep.Status,
) {
	if v.invalidCitation(citation) {
		status = rep.Invalid
		return
	}
	contract, status = v.storage_.ReadContract(citation)
	if v.invalidContract(contract) {
		status = rep.Invalid
	}
	return
}

func (v *validatedStorage_) WriteContract(
	contract not.ContractLike,
) (
	citation doc.ResourceLike,
	status rep.Status,
) {
	if v.invalidContract(contract) {
		status = rep.Invalid
		return
	}
	citation, status = v.storage_.WriteContract(contract)
	if status != rep.Written {
		return
	}
	if v.invalidCitation(citation) {
		status = rep.Invalid
	}
	return
}

func (v *validatedStorage_) DeleteContract(
	citation doc.ResourceLike,
) (
	status rep.Status,
) {
	if v.invalidCitation(citation) {
		status = rep.Invalid
		return
	}
	status = v.storage_.DeleteContract(citation)
	return
}

// PROTECTED INTERFACE

// Private Methods

func (v *validatedStorage_) invalidCitation(
	citation doc.ResourceLike,
) bool {
	var contract, status = v.storage_.ReadContract(citation)
	if status != rep.Retrieved {
		return true
	}
	var draft = contract.GetContent()
	return !v.notary_.CitationMatches(citation, draft)
}

func (v *validatedStorage_) invalidContract(
	contract not.ContractLike,
) bool {
	// Retrieve the citation to the certificate that signed the contract.
	var notary = contract.GetNotary()
	var draft = contract.GetContent()
	if !v.notary_.CitationMatches(notary, draft) {
		// Not self-signed, read the certificate that signed the contract.
		var document, status = v.storage_.ReadContract(notary)
		if status != rep.Retrieved {
			return true
		}
		draft = document.GetContent()
	}
	var certificate = not.Certificate(draft.AsString())
	return !v.notary_.SealMatches(contract, certificate)
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
