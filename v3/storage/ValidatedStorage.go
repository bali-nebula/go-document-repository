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

func (v *validatedStorage_) CitationExists(
	name fra.ResourceLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a document citation exists.",
	)

	// Determine whether or not the document citation exists.
	return v.storage_.CitationExists(name)
}

func (v *validatedStorage_) ReadCitation(
	name fra.ResourceLike,
) not.CitationLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a document citation.",
	)

	// Read the document citation from persistent storage.
	return v.storage_.ReadCitation(name)
}

func (v *validatedStorage_) WriteCitation(
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

func (v *validatedStorage_) DeleteCitation(
	name fra.ResourceLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to delete a document citation.",
	)

	// Delete the document citation from persistent storage.
	v.storage_.DeleteCitation(name)
}

func (v *validatedStorage_) DraftExists(
	citation not.CitationLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a draft document exists.",
	)

	// Determine whether or not the draft document exists.
	return v.storage_.DraftExists(citation)
}

func (v *validatedStorage_) ReadDraft(
	citation not.CitationLike,
) not.DraftLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a draft document.",
	)

	// Read the draft document from persistent storage.
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

func (v *validatedStorage_) WriteDraft(
	draft not.DraftLike,
) not.CitationLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a draft document.",
	)

	// Write the draft document to persistent storage.
	return v.storage_.WriteDraft(draft)
}

func (v *validatedStorage_) DeleteDraft(
	citation not.CitationLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to delete a draft document.",
	)

	// Delete the draft document from persistent storage.
	v.storage_.DeleteDraft(citation)
}

func (v *validatedStorage_) CertificateExists(
	citation not.CitationLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a certificate document exists.",
	)

	// Determine if the certificate document exists in validated storage.
	return v.storage_.CertificateExists(citation)
}

func (v *validatedStorage_) ReadCertificate(
	citation not.CitationLike,
) not.ContractLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a certificate document.",
	)

	// Attempt to read the certificate document from validated storage.
	var contract = v.storage_.ReadCertificate(citation)
	var draft = contract.GetDraft()
	var certificate = not.CertificateFromString(draft.AsString())
	if !v.notary_.SignatureMatches(contract, certificate) {
		var message = fmt.Sprintf(
			"The certificate does not match the cited certificate document: %s%s",
			certificate.AsString(),
			contract.AsString(),
		)
		panic(message)
	}
	return contract
}

func (v *validatedStorage_) WriteCertificate(
	contract not.ContractLike,
) not.CitationLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a certificate document.",
	)

	// Write the certificate document to persistent storage.
	var draft = contract.GetDraft()
	var certificate = not.CertificateFromString(draft.AsString())
	var previous = certificate.GetOptionalPrevious()
	if uti.IsDefined(previous) {
		draft = v.storage_.ReadCertificate(previous).GetDraft()
		certificate = not.CertificateFromString(draft.AsString())
	}
	if !v.notary_.SignatureMatches(contract, certificate) {
		var message = fmt.Sprintf(
			"The certificate does not match the certificate document: %s%s",
			certificate.AsString(),
			contract.AsString(),
		)
		panic(message)
	}
	return v.storage_.WriteCertificate(contract)
}

func (v *validatedStorage_) ContractExists(
	citation not.CitationLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a contract document exists.",
	)

	// Determine if the contract document exists in validated storage.
	return v.storage_.ContractExists(citation)
}

func (v *validatedStorage_) ReadContract(
	citation not.CitationLike,
) not.ContractLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a contract document.",
	)

	// Attempt to read the contract document from validated storage.
	var contract = v.storage_.ReadContract(citation)
	var draft = v.storage_.ReadCertificate(contract.GetCertificate()).GetDraft()
	var certificate = not.CertificateFromString(draft.AsString())
	if !v.notary_.SignatureMatches(contract, certificate) {
		var message = fmt.Sprintf(
			"The certificate does not match the cited contract document: %s%s",
			certificate.AsString(),
			contract.AsString(),
		)
		panic(message)
	}
	return contract
}

func (v *validatedStorage_) WriteContract(
	contract not.ContractLike,
) not.CitationLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a contract document.",
	)

	// Write the contract document to persistent storage.
	var draft = v.storage_.ReadCertificate(contract.GetCertificate()).GetDraft()
	var certificate = not.CertificateFromString(draft.AsString())
	if !v.notary_.SignatureMatches(contract, certificate) {
		var message = fmt.Sprintf(
			"The certificate does not match the contract document: %s%s",
			certificate.AsString(),
			contract.AsString(),
		)
		panic(message)
	}
	return v.storage_.WriteContract(contract)
}

func (v *validatedStorage_) BagExists(
	bag not.CitationLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a message bag exists.",
	)

	// Determine whether or not the message bag exists.
	return v.storage_.BagExists(bag)
}

func (v *validatedStorage_) ReadBag(
	bag not.CitationLike,
) not.ContractLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a message bag.",
	)

	// Read the message bag from persistent storage.
	var contract = v.storage_.ReadContract(bag)
	var draft = v.storage_.ReadCertificate(contract.GetCertificate()).GetDraft()
	var certificate = not.CertificateFromString(draft.AsString())
	if !v.notary_.SignatureMatches(contract, certificate) {
		var message = fmt.Sprintf(
			"The certificate does not match the cited bag: %s%s",
			certificate.AsString(),
			contract.AsString(),
		)
		panic(message)
	}
	return contract
}

func (v *validatedStorage_) WriteBag(
	bag not.ContractLike,
) not.CitationLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a message bag.",
	)

	// Create the new bag.
	var draft = v.storage_.ReadCertificate(bag.GetCertificate()).GetDraft()
	var certificate = not.CertificateFromString(draft.AsString())
	if !v.notary_.SignatureMatches(bag, certificate) {
		var message = fmt.Sprintf(
			"The certificate does not match the message bag: %s%s",
			certificate.AsString(),
			bag.AsString(),
		)
		panic(message)
	}
	return v.storage_.WriteBag(bag)
}

func (v *validatedStorage_) DeleteBag(
	bag not.CitationLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to delete a message bag.",
	)

	// Delete the bag and any remaining messages.
	v.storage_.DeleteBag(bag)
}

func (v *validatedStorage_) MessageCount(
	bag not.CitationLike,
) int {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while counting the messages in a message bag.",
	)

	// Determine the number of messages currently available in the bag.
	return v.storage_.MessageCount(bag)
}

func (v *validatedStorage_) ReadMessage(
	bag not.CitationLike,
) not.ContractLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a message from a message bag.",
	)

	// Read a random message from persistent storage.
	var contract = v.storage_.ReadMessage(bag)
	var draft = v.storage_.ReadCertificate(contract.GetCertificate()).GetDraft()
	var certificate = not.CertificateFromString(draft.AsString())
	if !v.notary_.SignatureMatches(contract, certificate) {
		var message = fmt.Sprintf(
			"The certificate does not match the returned message: %s%s",
			certificate.AsString(),
			contract.AsString(),
		)
		panic(message)
	}
	return contract
}

func (v *validatedStorage_) WriteMessage(
	bag not.CitationLike,
	message not.ContractLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a message to a message bag.",
	)

	// Write the message to the message bag in persistent storage.
	var draft = v.storage_.ReadCertificate(message.GetCertificate()).GetDraft()
	var certificate = not.CertificateFromString(draft.AsString())
	if !v.notary_.SignatureMatches(message, certificate) {
		var message = fmt.Sprintf(
			"The certificate does not match the message bag: %s%s",
			certificate.AsString(),
			message.AsString(),
		)
		panic(message)
	}
	v.storage_.WriteMessage(bag, message)
}

func (v *validatedStorage_) DeleteMessage(
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

func (v *validatedStorage_) ReleaseMessage(
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

func (v *validatedStorage_) WriteEvent(
	event not.ContractLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write an event.",
	)

	// Write the event to the notification queue in persistent storage.
	var draft = v.storage_.ReadCertificate(event.GetCertificate()).GetDraft()
	var certificate = not.CertificateFromString(draft.AsString())
	if !v.notary_.SignatureMatches(event, certificate) {
		var message = fmt.Sprintf(
			"The certificate does not match the event: %s%s",
			certificate.AsString(),
			event.AsString(),
		)
		panic(message)
	}
	v.storage_.WriteEvent(event)
}

// PROTECTED INTERFACE

// Private Methods

func (v *validatedStorage_) errorCheck(
	message string,
) {
	if e := recover(); e != nil {
		message = fmt.Sprintf(
			"ValidatedStorage: %s:\n    %v",
			message,
			e,
		)
		panic(message)
	}
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
