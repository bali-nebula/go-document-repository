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
	sha "crypto/sha512"
	doc "github.com/bali-nebula/go-bali-documents/v3"
	not "github.com/bali-nebula/go-digital-notary/v3"
	rep "github.com/bali-nebula/go-document-repository/v3/repository"
	uti "github.com/craterdog/go-missing-utilities/v7"
	sts "strings"
)

// CLASS INTERFACE

// Access Function

func LocalStorageClass() LocalStorageClassLike {
	return localStorageClass()
}

// Constructor Methods

func (c *localStorageClass_) LocalStorage(
	notary not.DigitalNotaryLike,
	directory string,
) LocalStorageLike {
	if uti.IsUndefined(notary) {
		panic("The \"notary\" attribute is required by this class.")
	}
	if uti.IsUndefined(directory) {
		panic("The \"directory\" attribute is required by this class.")
	}
	if !sts.HasSuffix(directory, "/") {
		// Make the directory name canonical.
		directory += "/"
	}
	var instance = &localStorage_{
		// Initialize the instance attributes.
		notary_:    notary,
		directory_: directory,
	}
	return instance
}

// Constant Methods

// Function Methods

// INSTANCE INTERFACE

// Principal Methods

func (v *localStorage_) GetClass() LocalStorageClassLike {
	return localStorageClass()
}

// Attribute Methods

// Persistent Methods

func (v *localStorage_) WriteCitation(
	name doc.NameLike,
	version doc.VersionLike,
	citation not.CitationLike,
) (
	status rep.Status,
) {
	var path = v.getPath(name)
	uti.MakeDirectory(path)
	var filename = v.getFilename(name, version)
	var source = citation.AsString()
	uti.WriteFile(filename, source)
	status = rep.Success
	return
}

func (v *localStorage_) ReadCitation(
	name doc.NameLike,
	version doc.VersionLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	var filename = v.getFilename(name, version)
	var source = uti.ReadFile(filename)
	citation = not.CitationFromString(source)
	status = rep.Success
	return
}

func (v *localStorage_) DeleteCitation(
	name doc.NameLike,
	version doc.VersionLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	// Remove the citation file.
	var filename = v.getFilename(name, version)
	uti.RemovePath(filename)

	// Remove any empty directories in the citation path.
	var path = v.getPath(name)
	for len(path) > 0 {
		if len(uti.ReadDirectory(path)) > 0 {
			// The directory is not empty so we are done.
			break
		}
		uti.RemovePath(path)
		var directories = sts.Split(path, "/")
		directories = directories[:len(directories)-1] // Strip off the last one.
		path = sts.Join(directories, "/")
	}
	status = rep.Success
	return
}

func (v *localStorage_) WriteMessage(
	bag doc.NameLike,
	message not.CitationLike,
) (
	status rep.Status,
) {
	var free = v.getPath(bag) + "/free"
	uti.MakeDirectory(free)
	var name = message.GetTag().AsString()[1:] + ".bali"
	var filename = free + "/" + name
	var source = message.AsString()
	uti.WriteFile(filename, source)
	status = rep.Success
	return
}

func (v *localStorage_) ReadMessage(
	bag doc.NameLike,
) (
	message not.CitationLike,
	status rep.Status,
) {
	var free = v.getPath(bag) + "/free"
	var names = doc.List[string](uti.ReadDirectory(free))
	var index = int(doc.Generator().RandomOrdinal(names.GetSize()))
	var name = names.GetValue(index)
	var filename = free + "/" + name
	var source = uti.ReadFile(filename)
	uti.RemovePath(filename)
	var lent = v.getPath(bag) + "/lent"
	uti.MakeDirectory(lent)
	filename = lent + "/" + name
	uti.WriteFile(filename, source)
	message = not.CitationFromString(source)
	status = rep.Success
	return
}

func (v *localStorage_) UnreadMessage(
	bag doc.NameLike,
	message not.CitationLike,
) (
	status rep.Status,
) {
	var lent = v.getPath(bag) + "/lent"
	var name = message.GetTag().AsString()[1:] + ".bali"
	var filename = lent + "/" + name
	uti.RemovePath(filename)
	var free = v.getPath(bag) + "/free"
	uti.MakeDirectory(free)
	filename = free + "/" + name
	var source = message.AsString()
	uti.WriteFile(filename, source)
	status = rep.Success
	return
}

func (v *localStorage_) DeleteMessage(
	bag doc.NameLike,
	message not.CitationLike,
) (
	status rep.Status,
) {
	var lent = v.getPath(bag) + "/lent"
	var name = message.GetTag().AsString()[1:] + ".bali"
	var filename = lent + "/" + name
	uti.RemovePath(filename)
	status = rep.Success
	return
}

func (v *localStorage_) WriteSubscription(
	bag doc.NameLike,
	type_ doc.ResourceLike,
) (
	status rep.Status,
) {
	var path = v.getPath(doc.Name("/subscriptions"))
	path += "/" + v.hashString(type_.AsString())
	uti.MakeDirectory(path)
	var name = v.hashString(bag.AsString()) + ".bali"
	var filename = path + "/" + name
	var source = bag.AsString()
	uti.WriteFile(filename, source)
	status = rep.Success
	return
}

func (v *localStorage_) ReadSubscriptions(
	type_ doc.ResourceLike,
) (
	bags doc.Sequential[doc.NameLike],
	status rep.Status,
) {
	var directory = v.hashString(type_.AsString())
	var path = v.getPath(doc.Name("/subscriptions")) + "/" + directory
	uti.MakeDirectory(path)
	var list = doc.List[doc.NameLike]()
	var files = uti.ReadDirectory(path)
	for _, file := range files {
		var name = path + "/" + file
		var source = uti.ReadFile(name)
		var bag = doc.Name(source)
		list.AppendValue(bag)
	}
	bags = list
	status = rep.Success
	return
}

func (v *localStorage_) DeleteSubscription(
	bag doc.NameLike,
	type_ doc.ResourceLike,
) (
	status rep.Status,
) {
	var path = v.getPath(doc.Name("/subscriptions"))
	path += "/" + v.hashString(type_.AsString())
	var name = v.hashString(bag.AsString()) + ".bali"
	var filename = path + "/" + name
	uti.RemovePath(filename)

	// Remove any empty directories in the citation path.
	for len(path) > 0 {
		if len(uti.ReadDirectory(path)) > 0 {
			// The directory is not empty so we are done.
			break
		}
		uti.RemovePath(path)
		var directories = sts.Split(path, "/")
		directories = directories[:len(directories)-1] // Strip off the last one.
		path = sts.Join(directories, "/")
	}
	status = rep.Success
	return
}

func (v *localStorage_) WriteDraft(
	draft not.DocumentLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	citation = v.notary_.CiteDocument(draft)
	var path = v.directory_ + "nebula/" + citation.GetTag().AsString()[1:] + "/"
	uti.MakeDirectory(path)
	var filename = path + citation.GetVersion().AsString() + ".bali"
	var source = draft.AsString()
	uti.WriteFile(filename, source)
	status = rep.Success
	return
}

func (v *localStorage_) ReadDraft(
	citation not.CitationLike,
) (
	draft not.DocumentLike,
	status rep.Status,
) {
	var path = v.directory_ + "nebula/" + citation.GetTag().AsString()[1:] + "/"
	var filename = path + citation.GetVersion().AsString() + ".bali"
	var source = uti.ReadFile(filename)
	draft = not.DocumentFromString(source)
	status = rep.Success
	return
}

func (v *localStorage_) DeleteDraft(
	citation not.CitationLike,
) (
	draft not.DocumentLike,
	status rep.Status,
) {
	var path = v.directory_ + "nebula/" + citation.GetTag().AsString()[1:] + "/"
	var filename = path + citation.GetVersion().AsString() + ".bali"
	var source = uti.ReadFile(filename)
	draft = not.DocumentFromString(source)
	uti.RemovePath(filename)
	var filenames = uti.ReadDirectory(path)
	if len(filenames) == 0 {
		// This was the last version of the document so delete the directory too.
		uti.RemovePath(path)
	}
	status = rep.Success
	return
}

func (v *localStorage_) WriteDocument(
	document not.DocumentLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	citation = v.notary_.CiteDocument(document)
	var path = v.directory_ + "nebula/" + citation.GetTag().AsString()[1:] + "/"
	uti.MakeDirectory(path)
	var filename = path + citation.GetVersion().AsString() + ".bali"
	var source = document.AsString()
	uti.WriteFile(filename, source)
	status = rep.Success
	return
}

func (v *localStorage_) ReadDocument(
	citation not.CitationLike,
) (
	document not.DocumentLike,
	status rep.Status,
) {
	var path = v.directory_ + "nebula/" + citation.GetTag().AsString()[1:] + "/"
	var filename = path + citation.GetVersion().AsString() + ".bali"
	var source = uti.ReadFile(filename)
	document = not.DocumentFromString(source)
	status = rep.Success
	return
}

func (v *localStorage_) DeleteDocument(
	citation not.CitationLike,
) (
	document not.DocumentLike,
	status rep.Status,
) {
	var path = v.directory_ + "nebula/" + citation.GetTag().AsString()[1:] + "/"
	var filename = path + citation.GetVersion().AsString() + ".bali"
	var source = uti.ReadFile(filename)
	document = not.DocumentFromString(source)
	uti.RemovePath(filename)
	var filenames = uti.ReadDirectory(path)
	if len(filenames) == 0 {
		// This was the last version of the document so delete the directory too.
		uti.RemovePath(path)
	}
	status = rep.Success
	return
}

// PROTECTED INTERFACE

// Private Methods

func (v *localStorage_) getPath(
	name doc.NameLike,
) string {
	return v.directory_ + "bali" + name.AsString()
}

func (v *localStorage_) getFilename(
	name doc.NameLike,
	version doc.VersionLike,
) string {
	return v.getPath(name) + "/" + version.AsString() + ".bali"
}

func (v *localStorage_) hashString(
	string_ string,
) string {
	var array = sha.Sum512([]byte(string_))
	var bytes = array[:20] // Convert the array to a slice of the first 20 bytes.
	var hash = doc.Tag(bytes).AsString()[1:]
	return hash
}

// Instance Structure

type localStorage_ struct {
	// Declare the instance attributes.
	notary_    not.DigitalNotaryLike
	directory_ string
}

// Class Structure

type localStorageClass_ struct {
	// Declare the class constants.
}

// Class Reference

func localStorageClass() *localStorageClass_ {
	return localStorageClassReference_
}

var localStorageClassReference_ = &localStorageClass_{
	// Initialize the class constants.
}
