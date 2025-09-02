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
) fra.ResourceLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a document citation.",
	)

	// Read the document citation from persistent storage.
	return v.storage_.ReadCitation(name)
}

func (v *validatedStorage_) WriteCitation(
	name fra.ResourceLike,
	citation fra.ResourceLike,
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
	citation fra.ResourceLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a draft document exists.",
	)

	// Determine whether or not the draft document exists.
	return v.storage_.DraftExists(citation)
}

func (v *validatedStorage_) ReadDraft(
	citation fra.ResourceLike,
) not.Parameterized {
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
	draft not.Parameterized,
) fra.ResourceLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a draft document.",
	)

	// Write the draft document to persistent storage.
	return v.storage_.WriteDraft(draft)
}

func (v *validatedStorage_) DeleteDraft(
	citation fra.ResourceLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to delete a draft document.",
	)

	// Delete the draft document from persistent storage.
	v.storage_.DeleteDraft(citation)
}

func (v *validatedStorage_) DocumentExists(
	citation fra.ResourceLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a notarized document exists.",
	)

	// Determine if the notarized document exists in validated storage.
	return v.storage_.DocumentExists(citation)
}

func (v *validatedStorage_) ReadDocument(
	citation fra.ResourceLike,
) not.Notarized {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a notarized document.",
	)

	// Attempt to read the notarized document from validated storage.
	var document = v.storage_.ReadDocument(citation)
	var draft = v.storage_.ReadCertificate(document.GetCertificate()).GetDraft()
	var certificate = not.CertificateFromString(draft.AsString())
	if !v.notary_.SignatureMatches(document, certificate) {
		var message = fmt.Sprintf(
			"The certificate does not match the cited notarized document: %s%s",
			certificate.AsString(),
			document.AsString(),
		)
		panic(message)
	}
	return document
}

func (v *validatedStorage_) WriteDocument(
	document not.Notarized,
) fra.ResourceLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a notarized document.",
	)

	// Write the notarized document to persistent storage.
	var draft = v.storage_.ReadCertificate(document.GetCertificate()).GetDraft()
	var certificate = not.CertificateFromString(draft.AsString())
	if !v.notary_.SignatureMatches(document, certificate) {
		var message = fmt.Sprintf(
			"The certificate does not match the notarized document: %s%s",
			certificate.AsString(),
			document.AsString(),
		)
		panic(message)
	}
	return v.storage_.WriteDocument(document)
}

func (v *validatedStorage_) BagExists(
	citation fra.ResourceLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a message bag exists.",
	)

	// Determine whether or not the message bag exists.
	return v.storage_.BagExists(citation)
}

func (v *validatedStorage_) ReadBag(
	citation fra.ResourceLike,
) not.Notarized {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a message bag.",
	)

	// Read the message bag from persistent storage.
	var bag = v.storage_.ReadDocument(citation)
	var draft = v.storage_.ReadCertificate(bag.GetCertificate()).GetDraft()
	var certificate = not.CertificateFromString(draft.AsString())
	if !v.notary_.SignatureMatches(bag, certificate) {
		var message = fmt.Sprintf(
			"The certificate does not match the cited bag: %s%s",
			certificate.AsString(),
			bag.AsString(),
		)
		panic(message)
	}
	return bag
}

func (v *validatedStorage_) WriteBag(
	bag not.Notarized,
) fra.ResourceLike {
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
	citation fra.ResourceLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to delete a message bag.",
	)

	// Delete the bag and any remaining messages.
	v.storage_.DeleteBag(citation)
}

func (v *validatedStorage_) MessageCount(
	bag fra.ResourceLike,
) uti.Cardinal {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while counting the messages in a message bag.",
	)

	// Determine the number of messages currently available in the bag.
	return v.storage_.MessageCount(bag)
}

func (v *validatedStorage_) ReadMessage(
	bag fra.ResourceLike,
) not.Notarized {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a message from a message bag.",
	)

	// Read a random message from persistent storage.
	var message = v.storage_.ReadMessage(bag)
	var draft = v.storage_.ReadCertificate(message.GetCertificate()).GetDraft()
	var certificate = not.CertificateFromString(draft.AsString())
	if !v.notary_.SignatureMatches(message, certificate) {
		var message = fmt.Sprintf(
			"The certificate does not match the returned message: %s%s",
			certificate.AsString(),
			message.AsString(),
		)
		panic(message)
	}
	return message
}

func (v *validatedStorage_) WriteMessage(
	bag fra.ResourceLike,
	message not.Notarized,
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
	bag fra.ResourceLike,
	citation fra.ResourceLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to delete a message from a message bag.",
	)

	// Delete the message from the message bag in persistent storage.
	v.storage_.DeleteMessage(bag, citation)
}

func (v *validatedStorage_) ReleaseMessage(
	bag fra.ResourceLike,
	citation fra.ResourceLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to reset the lease on a message.",
	)

	// Reset the message lease for the message in persistent storage.
	v.storage_.ReleaseMessage(bag, citation)
}

func (v *validatedStorage_) WriteEvent(
	event not.Notarized,
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
