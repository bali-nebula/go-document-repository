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
		cache_:   doc.Catalog[string, not.DocumentLike](),
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

func (v *cachedStorage_) WriteCitation(
	name doc.NameLike,
	version doc.VersionLike,
	citation not.CitationLike,
) (
	status rep.Status,
) {
	status = v.storage_.WriteCitation(name, version, citation)
	return
}

func (v *cachedStorage_) ReadCitation(
	name doc.NameLike,
	version doc.VersionLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	citation, status = v.storage_.ReadCitation(name, version)
	return
}

func (v *cachedStorage_) DeleteCitation(
	name doc.NameLike,
	version doc.VersionLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	citation, status = v.storage_.DeleteCitation(name, version)
	return
}

func (v *cachedStorage_) WriteMessage(
	bag doc.NameLike,
	message not.CitationLike,
) (
	status rep.Status,
) {
	status = v.storage_.WriteMessage(bag, message)
	return
}

func (v *cachedStorage_) ReadMessage(
	bag doc.NameLike,
) (
	message not.CitationLike,
	status rep.Status,
) {
	message, status = v.storage_.ReadMessage(bag)
	return
}

func (v *cachedStorage_) UnreadMessage(
	bag doc.NameLike,
	message not.CitationLike,
) (
	status rep.Status,
) {
	status = v.storage_.UnreadMessage(bag, message)
	return
}

func (v *cachedStorage_) DeleteMessage(
	bag doc.NameLike,
	message not.CitationLike,
) (
	status rep.Status,
) {
	status = v.storage_.DeleteMessage(bag, message)
	return
}

func (v *cachedStorage_) WriteSubscription(
	bag doc.NameLike,
	type_ doc.ResourceLike,
) (
	status rep.Status,
) {
	status = v.storage_.WriteSubscription(bag, type_)
	return
}

func (v *cachedStorage_) ReadSubscriptions(
	type_ doc.ResourceLike,
) (
	bags doc.Sequential[doc.NameLike],
	status rep.Status,
) {
	bags, status = v.storage_.ReadSubscriptions(type_)
	return
}

func (v *cachedStorage_) DeleteSubscription(
	bag doc.NameLike,
	type_ doc.ResourceLike,
) (
	status rep.Status,
) {
	status = v.storage_.DeleteSubscription(bag, type_)
	return
}

func (v *cachedStorage_) WriteDraft(
	draft not.DocumentLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	citation, status = v.storage_.WriteDraft(draft)
	return
}

func (v *cachedStorage_) ReadDraft(
	citation not.CitationLike,
) (
	draft not.DocumentLike,
	status rep.Status,
) {
	draft, status = v.storage_.ReadDraft(citation)
	return
}

func (v *cachedStorage_) DeleteDraft(
	citation not.CitationLike,
) (
	draft not.DocumentLike,
	status rep.Status,
) {
	draft, status = v.storage_.DeleteDraft(citation)
	return
}

func (v *cachedStorage_) WriteDocument(
	document not.DocumentLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	citation, status = v.storage_.WriteDocument(document)
	if status == rep.Success {
		v.cacheDocument(citation, document)
	}
	return
}

func (v *cachedStorage_) ReadDocument(
	citation not.CitationLike,
) (
	document not.DocumentLike,
	status rep.Status,
) {
	document = v.lookupDocument(citation)
	if uti.IsDefined(document) {
		status = rep.Success
		return
	}
	document, status = v.storage_.ReadDocument(citation)
	if status == rep.Success {
		v.cacheDocument(citation, document)
	}
	return
}

func (v *cachedStorage_) DeleteDocument(
	citation not.CitationLike,
) (
	document not.DocumentLike,
	status rep.Status,
) {
	document, status = v.storage_.DeleteDocument(citation)
	if status == rep.Success {
		v.uncacheDocument(citation)
	}
	return
}

// PROTECTED INTERFACE

// Private Methods

func (v *cachedStorage_) cacheDocument(
	citation not.CitationLike,
	document not.DocumentLike,
) {
	var key = v.getKey(citation)
	v.cache_.SetValue(key, document)
}

func (v *cachedStorage_) lookupDocument(
	citation not.CitationLike,
) not.DocumentLike {
	var key = v.getKey(citation)
	var document = v.cache_.GetValue(key)
	return document
}

func (v *cachedStorage_) uncacheDocument(
	citation not.CitationLike,
) {
	var key = v.getKey(citation)
	v.cache_.RemoveValue(key)
}

func (v *cachedStorage_) getKey(
	citation not.CitationLike,
) string {
	var tag = citation.GetTag().AsString()[1:] // Remove the leading "#" character.
	var version = citation.GetVersion().AsString()
	return tag + ":" + version
}

// Instance Structure

type cachedStorage_ struct {
	// Declare the instance attributes.
	cache_   doc.CatalogLike[string, not.DocumentLike]
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
