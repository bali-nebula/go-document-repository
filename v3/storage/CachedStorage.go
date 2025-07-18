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
	not "github.com/bali-nebula/go-digital-notary/v3"
	rep "github.com/bali-nebula/go-document-repository/v3/repository"
	fra "github.com/craterdog/go-component-framework/v7"
	uti "github.com/craterdog/go-missing-utilities/v7"
)

// CLASS INTERFACE

// Access Function

func CachedStorageClass() CachedStorageClassLike {
	return cachedStorageClass()
}

// Constructor Methods

func (c *cachedStorageClass_) CachedStorage(
	storage rep.Persistent,
) CachedStorageLike {
	if uti.IsUndefined(storage) {
		panic("The \"storage\" attribute is required by this class.")
	}
	var instance = &cachedStorage_{
		// Initialize the instance attributes.
		cache_:   fra.Catalog[string, not.ContractLike](),
		storage_: storage,
	}
	return instance
}

// Constant Methods

// Function Methods

// INSTANCE INTERFACE

// Principal Methods

func (v *cachedStorage_) GetClass() CachedStorageClassLike {
	return cachedStorageClass()
}

// Attribute Methods

// Persistent Methods

func (v *cachedStorage_) CitationExists(
	name fra.ResourceLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a document citation exists.",
	)

	// Determine whether or not the document citation exists.
	return v.storage_.CitationExists(name)
}

func (v *cachedStorage_) ReadCitation(
	name fra.ResourceLike,
) not.CitationLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a document citation.",
	)

	// Read the document citation from persistent storage.
	return v.storage_.ReadCitation(name)
}

func (v *cachedStorage_) WriteCitation(
	name fra.ResourceLike,
	citation not.CitationLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a document citation.",
	)

	// Write the document citation to persistent storage.
	v.storage_.WriteCitation(name, citation)
}

func (v *cachedStorage_) DeleteCitation(
	name fra.ResourceLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to delete a document citation.",
	)

	// Delete the document citation from persistent storage.
	v.storage_.DeleteCitation(name)
}

func (v *cachedStorage_) DraftExists(
	citation not.CitationLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a draft document exists.",
	)

	// Determine whether or not the draft document exists.
	return v.storage_.DraftExists(citation)
}

func (v *cachedStorage_) ReadDraft(
	citation not.CitationLike,
) not.DraftLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a draft document.",
	)

	// Read the draft document from persistent storage.
	return v.storage_.ReadDraft(citation)
}

func (v *cachedStorage_) WriteDraft(
	draft not.DraftLike,
) not.CitationLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a draft document.",
	)

	// Write the draft document to persistent storage.
	return v.storage_.WriteDraft(draft)
}

func (v *cachedStorage_) DeleteDraft(
	citation not.CitationLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to delete a draft document.",
	)

	// Delete the draft document from persistent storage.
	v.storage_.DeleteDraft(citation)
}

func (v *cachedStorage_) CertificateExists(
	citation not.CitationLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a certificate document exists.",
	)

	// Determine if the certificate document exists in cached storage.
	var certificate = v.lookupContract(citation)
	if uti.IsUndefined(certificate) {
		// Determine if the certificate document exists in persistent storage.
		return v.storage_.CertificateExists(citation)
	}
	return true
}

func (v *cachedStorage_) ReadCertificate(
	citation not.CitationLike,
) not.ContractLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a certificate document.",
	)

	// Attempt to read the certificate document from cached storage.
	var certificate = v.lookupContract(citation)
	if uti.IsUndefined(certificate) {
		// Read the certificate document from persistent storage.
		certificate = v.storage_.ReadCertificate(citation)
		if uti.IsDefined(certificate) {
			v.cacheContract(citation, certificate)
		}
	}
	return certificate
}

func (v *cachedStorage_) WriteCertificate(
	certificate not.ContractLike,
) not.CitationLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a certificate document.",
	)

	// Write the certificate document to persistent storage.
	var citation = v.storage_.WriteCertificate(certificate)
	v.cacheContract(citation, certificate)
	return citation
}

func (v *cachedStorage_) ContractExists(
	citation not.CitationLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a contract document exists.",
	)

	// Determine if the contract document exists in cached storage.
	var contract = v.lookupContract(citation)
	if uti.IsUndefined(contract) {
		// Determine if the contract document exists in persistent storage.
		return v.storage_.ContractExists(citation)
	}
	return true
}

func (v *cachedStorage_) ReadContract(
	citation not.CitationLike,
) not.ContractLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a contract document.",
	)

	// Attempt to read the contract document from cached storage.
	var contract = v.lookupContract(citation)
	if uti.IsUndefined(contract) {
		// Read the contract document from persistent storage.
		contract = v.storage_.ReadContract(citation)
		if uti.IsDefined(contract) {
			v.cacheContract(citation, contract)
		}
	}
	return contract
}

func (v *cachedStorage_) WriteContract(
	contract not.ContractLike,
) not.CitationLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a contract document.",
	)

	// Write the contract document to persistent storage.
	var citation = v.storage_.WriteContract(contract)
	v.cacheContract(citation, contract)
	return citation
}

func (v *cachedStorage_) BagExists(
	bag not.CitationLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a message bag exists.",
	)

	// Determine whether or not the message bag exists.
	return v.storage_.BagExists(bag)
}

func (v *cachedStorage_) ReadBag(
	bag not.CitationLike,
) not.ContractLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a message bag.",
	)

	// Read the message bag from persistent storage.
	return v.storage_.ReadBag(bag)
}

func (v *cachedStorage_) WriteBag(
	bag not.ContractLike,
) not.CitationLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a message bag.",
	)

	// Create the new bag.
	return v.storage_.WriteBag(bag)
}

func (v *cachedStorage_) DeleteBag(
	bag not.CitationLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to delete a message bag.",
	)

	// Delete the bag and any remaining messages.
	v.storage_.DeleteBag(bag)
}

func (v *cachedStorage_) MessageCount(
	bag not.CitationLike,
) int {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while counting the messages in a message bag.",
	)

	// Determine the number of messages currently available in the bag.
	return v.storage_.MessageCount(bag)
}

func (v *cachedStorage_) ReadMessage(
	bag not.CitationLike,
) not.ContractLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a message from a message bag.",
	)

	// Read a random message from persistent storage.
	return v.storage_.ReadMessage(bag)
}

func (v *cachedStorage_) WriteMessage(
	bag not.CitationLike,
	message not.ContractLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a message to a message bag.",
	)

	// Write the message to the message bag in persistent storage.
	v.storage_.WriteMessage(bag, message)
}

func (v *cachedStorage_) DeleteMessage(
	bag not.CitationLike,
	message not.CitationLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to delete a message from a message bag.",
	)

	// Delete the message from the message bag in persistent storage.
	v.storage_.DeleteMessage(bag, message)
}

func (v *cachedStorage_) ReleaseMessage(
	bag not.CitationLike,
	message not.CitationLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to reset the lease on a message.",
	)

	// Reset the message lease for the message in persistent storage.
	v.storage_.ReleaseMessage(bag, message)
}

func (v *cachedStorage_) WriteEvent(
	event not.ContractLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write an event.",
	)

	// Write the event to the notification queue in persistent storage.
	v.storage_.WriteEvent(event)
}

// PROTECTED INTERFACE

// Private Methods

func (v *cachedStorage_) cacheContract(
	citation not.CitationLike,
	contract not.ContractLike,
) {
	var key = v.getCitationTag(citation) + v.getCitationVersion(citation)
	v.cache_.SetValue(key, contract)
}

func (v *cachedStorage_) errorCheck(
	message string,
) {
	if e := recover(); e != nil {
		message = fmt.Sprintf(
			"CachedStorage: %s:\n    %v",
			message,
			e,
		)
		panic(message)
	}
}

func (v *cachedStorage_) getCitationTag(
	citation not.CitationLike,
) string {
	var tag = citation.GetTag()
	return tag.AsString()[1:] // Remove the leading "#" character.
}

func (v *cachedStorage_) getCitationVersion(
	citation not.CitationLike,
) string {
	var version = citation.GetVersion()
	return version.AsString()
}

func (v *cachedStorage_) lookupContract(
	citation not.CitationLike,
) not.ContractLike {
	var key = v.getCitationTag(citation) + v.getCitationVersion(citation)
	return v.cache_.GetValue(key)
}

// Instance Structure

type cachedStorage_ struct {
	// Declare the instance attributes.
	cache_   fra.CatalogLike[string, not.ContractLike]
	storage_ rep.Persistent
}

// Class Structure

type cachedStorageClass_ struct {
	// Declare the class constants.
}

// Class Reference

func cachedStorageClass() *cachedStorageClass_ {
	return cachedStorageClassReference_
}

var cachedStorageClassReference_ = &cachedStorageClass_{
	// Initialize the class constants.
}
