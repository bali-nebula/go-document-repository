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
		cache_:   doc.Catalog[string, not.ContractLike](),
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

func (v *cachedStorage_) ReadCitation(
	name doc.NameLike,
	version doc.VersionLike,
) (
	citation doc.ResourceLike,
	status rep.Status,
) {
	return v.storage_.ReadCitation(name, version)
}

func (v *cachedStorage_) WriteCitation(
	name doc.NameLike,
	version doc.VersionLike,
	citation doc.ResourceLike,
) (
	status rep.Status,
) {
	return v.storage_.WriteCitation(name, version, citation)
}

func (v *cachedStorage_) DeleteCitation(
	name doc.NameLike,
	version doc.VersionLike,
) (
	status rep.Status,
) {
	return v.storage_.DeleteCitation(name, version)
}

func (v *cachedStorage_) ListCitations(
	path doc.NameLike,
) (
	citations doc.Sequential[doc.ResourceLike],
	status rep.Status,
) {
	return v.storage_.ListCitations(path)
}

func (v *cachedStorage_) ReadDraft(
	citation doc.ResourceLike,
) (
	draft not.Parameterized,
	status rep.Status,
) {
	return v.storage_.ReadDraft(citation)
}

func (v *cachedStorage_) WriteDraft(
	draft not.Parameterized,
) (
	citation doc.ResourceLike,
	status rep.Status,
) {
	return v.storage_.WriteDraft(draft)
}

func (v *cachedStorage_) DeleteDraft(
	citation doc.ResourceLike,
) (
	status rep.Status,
) {
	return v.storage_.DeleteDraft(citation)
}

func (v *cachedStorage_) ReadContract(
	citation doc.ResourceLike,
) (
	contract not.ContractLike,
	status rep.Status,
) {
	contract = v.lookupContract(citation)
	if uti.IsDefined(contract) {
		status = rep.Retrieved
		return
	}
	contract, status = v.storage_.ReadContract(citation)
	if status == rep.Retrieved {
		v.cacheContract(citation, contract)
	}
	return
}

func (v *cachedStorage_) WriteContract(
	contract not.ContractLike,
) (
	citation doc.ResourceLike,
	status rep.Status,
) {
	citation, status = v.storage_.WriteContract(contract)
	if status == rep.Written {
		v.cacheContract(citation, contract)
	}
	return
}

func (v *cachedStorage_) DeleteContract(
	citation doc.ResourceLike,
) (
	status rep.Status,
) {
	status = v.storage_.DeleteContract(citation)
	if status == rep.Deleted {
		v.uncacheContract(citation)
	}
	return
}

// PROTECTED INTERFACE

// Private Methods

func (v *cachedStorage_) cacheContract(
	citation doc.ResourceLike,
	contract not.ContractLike,
) {
	var key = v.getCitationTag(citation) + v.getCitationVersion(citation)
	v.cache_.SetValue(key, contract)
}

func (v *cachedStorage_) lookupContract(
	citation doc.ResourceLike,
) not.ContractLike {
	var key = v.getCitationTag(citation) + v.getCitationVersion(citation)
	var contract = v.cache_.GetValue(key)
	return contract
}

func (v *cachedStorage_) uncacheContract(
	citation doc.ResourceLike,
) {
	var key = v.getCitationTag(citation) + v.getCitationVersion(citation)
	v.cache_.RemoveValue(key)
}

func (v *cachedStorage_) getCitationTag(
	resource doc.ResourceLike,
) string {
	var citation = not.Citation(resource)
	var tag = citation.GetTag()
	return tag.AsString()[1:] // Remove the leading "#" character.
}

func (v *cachedStorage_) getCitationVersion(
	resource doc.ResourceLike,
) string {
	var citation = not.Citation(resource)
	var version = citation.GetVersion()
	return version.AsString()
}

// Instance Structure

type cachedStorage_ struct {
	// Declare the instance attributes.
	cache_   doc.CatalogLike[string, not.ContractLike]
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
