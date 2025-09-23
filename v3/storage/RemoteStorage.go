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
	uti "github.com/craterdog/go-missing-utilities/v7"
)

// CLASS INTERFACE

// Access Function

func RemoteStorageClass() RemoteStorageClassLike {
	return remoteStorageClass()
}

// Constructor Methods

func (c *remoteStorageClass_) RemoteStorage(
	notary not.DigitalNotaryLike,
	service doc.ResourceLike,
) RemoteStorageLike {
	if uti.IsUndefined(notary) {
		panic("The \"notary\" attribute is required by this class.")
	}
	if uti.IsUndefined(service) {
		panic("The \"service\" attribute is required by this class.")
	}
	var instance = &remoteStorage_{
		// Initialize the instance attributes.
		notary_:  notary,
		service_: service,
	}
	return instance
}

// Constant Methods

// Function Methods

// INSTANCE INTERFACE

// Principal Methods

func (v *remoteStorage_) GetClass() RemoteStorageClassLike {
	return remoteStorageClass()
}

// Attribute Methods

// Persistent Methods

func (v *remoteStorage_) ReadCitation(
	name doc.NameLike,
	version doc.VersionLike,
) doc.ResourceLike {
	var result_ doc.ResourceLike
	// TBD - Add the method implementation.
	return result_
}

func (v *remoteStorage_) WriteCitation(
	name doc.NameLike,
	version doc.VersionLike,
	citation doc.ResourceLike,
) {
	// TBD - Add the method implementation.
}

func (v *remoteStorage_) DeleteCitation(
	name doc.NameLike,
	version doc.VersionLike,
) {
	// TBD - Add the method implementation.
}

func (v *remoteStorage_) ListCitations(
	path doc.NameLike,
) doc.Sequential[doc.ResourceLike] {
	var result_ doc.Sequential[doc.ResourceLike]
	// TBD - Add the method implementation.
	return result_
}

func (v *remoteStorage_) ReadDraft(
	citation doc.ResourceLike,
) not.Parameterized {
	var result_ not.Parameterized
	// TBD - Add the method implementation.
	return result_
}

func (v *remoteStorage_) WriteDraft(
	draft not.Parameterized,
) doc.ResourceLike {
	var result_ doc.ResourceLike
	// TBD - Add the method implementation.
	return result_
}

func (v *remoteStorage_) DeleteDraft(
	citation doc.ResourceLike,
) {
	// TBD - Add the method implementation.
}

func (v *remoteStorage_) ReadContract(
	citation doc.ResourceLike,
) not.ContractLike {
	var result_ not.ContractLike
	// TBD - Add the method implementation.
	return result_
}

func (v *remoteStorage_) WriteContract(
	contract not.ContractLike,
) doc.ResourceLike {
	var result_ doc.ResourceLike
	// TBD - Add the method implementation.
	return result_
}

func (v *remoteStorage_) DeleteContract(
	citation doc.ResourceLike,
) {
	// TBD - Add the method implementation.
}

// PROTECTED INTERFACE

// Private Methods

// Instance Structure

type remoteStorage_ struct {
	// Declare the instance attributes.
	notary_  not.DigitalNotaryLike
	service_ doc.ResourceLike
}

// Class Structure

type remoteStorageClass_ struct {
	// Declare the class constants.
}

// Class Reference

func remoteStorageClass() *remoteStorageClass_ {
	return remoteStorageClassReference_
}

var remoteStorageClassReference_ = &remoteStorageClass_{
	// Initialize the class constants.
}
