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
		cache_:   fra.Catalog[string, not.Notarized](),
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
	name fra.NameLike,
	version fra.VersionLike,
) fra.ResourceLike {
	return v.storage_.ReadCitation(name, version)
}

func (v *cachedStorage_) WriteCitation(
	name fra.NameLike,
	version fra.VersionLike,
	citation fra.ResourceLike,
) {
	v.storage_.WriteCitation(name, version, citation)
}

func (v *cachedStorage_) DeleteCitation(
	name fra.NameLike,
	version fra.VersionLike,
) {
	v.storage_.DeleteCitation(name, version)
}

func (v *cachedStorage_) ListCitations(
	path fra.NameLike,
) fra.Sequential[fra.ResourceLike] {
	return v.storage_.ListCitations(path)
}

func (v *cachedStorage_) ReadDraft(
	citation fra.ResourceLike,
) not.Parameterized {
	return v.storage_.ReadDraft(citation)
}

func (v *cachedStorage_) WriteDraft(
	draft not.Parameterized,
) fra.ResourceLike {
	return v.storage_.WriteDraft(draft)
}

func (v *cachedStorage_) DeleteDraft(
	citation fra.ResourceLike,
) {
	v.storage_.DeleteDraft(citation)
}

func (v *cachedStorage_) ReadContract(
	citation fra.ResourceLike,
) not.Notarized {
	// Attempt to read the notarized document from cached storage.
	var document = v.lookupContract(citation)
	if uti.IsUndefined(document) {
		// Read the notarized document from persistent storage.
		document = v.storage_.ReadContract(citation)
		if uti.IsDefined(document) {
			v.cacheContract(citation, document)
		}
	}
	return document
}

func (v *cachedStorage_) WriteContract(
	contract not.Notarized,
) fra.ResourceLike {
	var citation = v.storage_.WriteContract(contract)
	v.cacheContract(citation, contract)
	return citation
}

func (v *cachedStorage_) DeleteContract(
	citation fra.ResourceLike,
) {
	// Delete the notarized document from persistent storage.
	v.storage_.DeleteContract(citation)

	// Remove the notarized document from cached storage.
	v.uncacheContract(citation)
}

// PROTECTED INTERFACE

// Private Methods

func (v *cachedStorage_) cacheContract(
	citation fra.ResourceLike,
	contract not.Notarized,
) {
	var key = v.getCitationTag(citation) + v.getCitationVersion(citation)
	v.cache_.SetValue(key, contract)
}

func (v *cachedStorage_) lookupContract(
	citation fra.ResourceLike,
) not.Notarized {
	var key = v.getCitationTag(citation) + v.getCitationVersion(citation)
	return v.cache_.GetValue(key)
}

func (v *cachedStorage_) uncacheContract(
	citation fra.ResourceLike,
) {
	var key = v.getCitationTag(citation) + v.getCitationVersion(citation)
	v.cache_.RemoveValue(key)
}

func (v *cachedStorage_) getCitationTag(
	resource fra.ResourceLike,
) string {
	var citation = not.CitationFromResource(resource)
	var tag = citation.GetTag()
	return tag.AsString()[1:] // Remove the leading "#" character.
}

func (v *cachedStorage_) getCitationVersion(
	resource fra.ResourceLike,
) string {
	var citation = not.CitationFromResource(resource)
	var version = citation.GetVersion()
	return version.AsString()
}

// Instance Structure

type cachedStorage_ struct {
	// Declare the instance attributes.
	cache_   fra.CatalogLike[string, not.Notarized]
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
