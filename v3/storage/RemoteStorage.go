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
	not "github.com/bali-nebula/go-digital-notary/v3"
	fra "github.com/craterdog/go-component-framework/v7"
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
	service fra.ResourceLike,
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
	name fra.NameLike,
	version fra.VersionLike,
) fra.ResourceLike {
	var result_ fra.ResourceLike
	// TBD - Add the method implementation.
	return result_
}

func (v *remoteStorage_) WriteCitation(
	name fra.NameLike,
	version fra.VersionLike,
	citation fra.ResourceLike,
) {
	// TBD - Add the method implementation.
}

func (v *remoteStorage_) DeleteCitation(
	name fra.NameLike,
	version fra.VersionLike,
) {
	// TBD - Add the method implementation.
}

func (v *remoteStorage_) ListCitations(
	path fra.NameLike,
) fra.Sequential[fra.ResourceLike] {
	var result_ fra.Sequential[fra.ResourceLike]
	// TBD - Add the method implementation.
	return result_
}

func (v *remoteStorage_) ReadDraft(
	citation fra.ResourceLike,
) not.Parameterized {
	var result_ not.Parameterized
	// TBD - Add the method implementation.
	return result_
}

func (v *remoteStorage_) WriteDraft(
	draft not.Parameterized,
) fra.ResourceLike {
	var result_ fra.ResourceLike
	// TBD - Add the method implementation.
	return result_
}

func (v *remoteStorage_) DeleteDraft(
	citation fra.ResourceLike,
) {
	// TBD - Add the method implementation.
}

func (v *remoteStorage_) ReadContract(
	citation fra.ResourceLike,
) not.Notarized {
	var result_ not.Notarized
	// TBD - Add the method implementation.
	return result_
}

func (v *remoteStorage_) WriteContract(
	contract not.Notarized,
) fra.ResourceLike {
	var result_ fra.ResourceLike
	// TBD - Add the method implementation.
	return result_
}

func (v *remoteStorage_) DeleteContract(
	citation fra.ResourceLike,
) {
	// TBD - Add the method implementation.
}

// PROTECTED INTERFACE

// Private Methods

// Instance Structure

type remoteStorage_ struct {
	// Declare the instance attributes.
	notary_  not.DigitalNotaryLike
	service_ fra.ResourceLike
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
