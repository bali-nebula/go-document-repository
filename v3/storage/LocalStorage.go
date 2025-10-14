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
	citation not.CitationLike,
) (
	status rep.Status,
) {
	var filename = v.directory_ + "citations" + name.AsSource() + ".bali"
	var path = v.extractPath(filename)
	uti.MakeDirectory(path)
	if uti.PathExists(filename) {
		status = rep.Existed
		return
	}
	var source = citation.AsSource()
	uti.WriteFile(filename, source)
	status = rep.Success
	return
}

func (v *localStorage_) ReadCitation(
	name doc.NameLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	var filename = v.directory_ + "citations" + name.AsSource() + ".bali"
	if !uti.PathExists(filename) {
		status = rep.Missing
		return
	}
	var source = uti.ReadFile(filename)
	citation = not.Citation(source)
	status = rep.Success
	return
}

func (v *localStorage_) DeleteCitation(
	name doc.NameLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	// Remove the citation file.
	var filename = v.directory_ + "citations" + name.AsSource() + ".bali"
	var path = v.extractPath(filename)
	uti.RemovePath(filename)

	// Remove any empty directories in the citation path.
	for len(path) > 0 {
		if len(uti.ReadDirectory(path)) > 0 {
			// The directory is not empty so we are done.
			break
		}
		uti.RemovePath(path)
		path = v.extractPath(path)
	}
	status = rep.Success
	return
}

func (v *localStorage_) WriteMessage(
	bag doc.NameLike,
	message not.DocumentLike,
) (
	status rep.Status,
) {
	// Generate the message citation.
	var citation = v.notary_.CiteDocument(message)

	// Write the message citation to the free directory in local storage.
	var path = v.directory_ + "citations" + bag.AsSource() + "/free"
	var name = citation.GetTag().AsSource()[1:]
	var filename = path + "/" + name + ".bali"
	uti.MakeDirectory(path)
	if uti.PathExists(filename) {
		status = rep.Existed
		return
	}
	var source = citation.AsSource()
	uti.WriteFile(filename, source)

	// Write the message to the messages directory in local storage.
	path = v.directory_ + "messages"
	uti.MakeDirectory(path)
	filename = path + "/" + name + ".bali"
	source = message.AsSource()
	uti.WriteFile(filename, source)

	status = rep.Success
	return
}

func (v *localStorage_) ReadMessage(
	bag doc.NameLike,
) (
	message not.DocumentLike,
	status rep.Status,
) {
	// Select a random message citation from the free directory in local storage.
	var path = v.directory_ + "citations" + bag.AsSource() + "/free"
	var names = doc.List[string](uti.ReadDirectory(path))
	if names.IsEmpty() {
		status = rep.Missing
		return
	}
	var index = int(doc.Generator().RandomOrdinal(names.GetSize()))
	var name = names.GetValue(index)

	// Delete the selected message citation from the free directory in local storage.
	var filename = path + "/" + name
	if !uti.PathExists(filename) {
		status = rep.Missing
		return
	}
	var source = uti.ReadFile(filename)
	uti.RemovePath(filename)

	// Write the selected message citation to the lent directory in local storage.
	path = v.directory_ + "citations" + bag.AsSource() + "/lent"
	uti.MakeDirectory(path)
	filename = path + "/" + name
	uti.WriteFile(filename, source)

	// Read the message from the messages directory in local storage.
	var citation = not.Citation(source)
	name = citation.GetTag().AsSource()[1:]
	path = v.directory_ + "messages"
	uti.MakeDirectory(path)
	filename = path + "/" + name + ".bali"
	if !uti.PathExists(filename) {
		status = rep.Missing
		return
	}
	source = uti.ReadFile(filename)
	message = not.Document(source)

	status = rep.Success
	return
}

func (v *localStorage_) UnreadMessage(
	bag doc.NameLike,
	message not.DocumentLike,
) (
	status rep.Status,
) {
	// Generate the message citation.
	var citation = v.notary_.CiteDocument(message)
	var name = citation.GetTag().AsSource()[1:]

	// Move the message citation in local storage from lent to free.
	var path = v.directory_ + "citations" + bag.AsSource() + "/lent"
	var filename = path + "/" + name + ".bali"
	uti.RemovePath(filename)
	path = v.directory_ + "citations" + bag.AsSource() + "/free"
	uti.MakeDirectory(path)
	filename = path + "/" + name + ".bali"
	if uti.PathExists(filename) {
		status = rep.Missing
		return
	}
	var source = citation.AsSource()
	uti.WriteFile(filename, source)

	status = rep.Success
	return
}

func (v *localStorage_) DeleteMessage(
	bag doc.NameLike,
	message not.DocumentLike,
) (
	status rep.Status,
) {
	// Generate the message citation.
	var citation = v.notary_.CiteDocument(message)
	var name = citation.GetTag().AsSource()[1:]

	// Delete the message from local storage.
	var path = v.directory_ + "messages"
	var filename = path + "/" + name + ".bali"
	uti.RemovePath(filename)

	// Delete the message citation from local storage.
	path = v.directory_ + "citations" + bag.AsSource() + "/lent"
	filename = path + "/" + name + ".bali"
	uti.RemovePath(filename)

	status = rep.Success
	return
}

func (v *localStorage_) WriteSubscription(
	bag doc.NameLike,
	type_ doc.NameLike,
) (
	status rep.Status,
) {
	var path = v.directory_ + "subscriptions/" + v.hashString(type_.AsSource())
	uti.MakeDirectory(path)
	var name = v.hashString(bag.AsSource()) + ".bali"
	var filename = path + "/" + name
	var source = bag.AsSource()
	uti.WriteFile(filename, source)
	status = rep.Success
	return
}

func (v *localStorage_) ReadSubscriptions(
	type_ doc.NameLike,
) (
	bags doc.Sequential[doc.NameLike],
	status rep.Status,
) {
	var path = v.directory_ + "subscriptions/" + v.hashString(type_.AsSource())
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
	type_ doc.NameLike,
) (
	status rep.Status,
) {
	var path = v.directory_ + "subscriptions/" + v.hashString(type_.AsSource())
	var filename = path + "/" + v.hashString(bag.AsSource()) + ".bali"
	uti.RemovePath(filename)
	var filenames = uti.ReadDirectory(path)
	if len(filenames) == 0 {
		uti.RemovePath(path)
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
	// Generate the draft document citation.
	citation = v.notary_.CiteDocument(draft)
	var name = citation.GetTag().AsSource()[1:]
	var version = citation.GetVersion().AsSource()

	// Write the draft document to local storage.
	var path = v.directory_ + "drafts/" + name
	uti.MakeDirectory(path)
	var filename = path + "/" + version + ".bali"
	var source = draft.AsSource()
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
	var name = citation.GetTag().AsSource()[1:]
	var version = citation.GetVersion().AsSource()
	var path = v.directory_ + "drafts/" + name
	var filename = path + "/" + version + ".bali"
	if !uti.PathExists(filename) {
		status = rep.Missing
		return
	}
	var source = uti.ReadFile(filename)
	draft = not.Document(source)
	status = rep.Success
	return
}

func (v *localStorage_) DeleteDraft(
	citation not.CitationLike,
) (
	draft not.DocumentLike,
	status rep.Status,
) {
	var name = citation.GetTag().AsSource()[1:]
	var version = citation.GetVersion().AsSource()
	var path = v.directory_ + "drafts/" + name
	var filename = path + "/" + version + ".bali"
	if !uti.PathExists(filename) {
		status = rep.Missing
		return
	}
	var source = uti.ReadFile(filename)
	draft = not.Document(source)
	uti.RemovePath(filename)
	var filenames = uti.ReadDirectory(path)
	if len(filenames) == 0 {
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
	// Generate the document citation.
	citation = v.notary_.CiteDocument(document)
	var name = citation.GetTag().AsSource()[1:]
	var version = citation.GetVersion().AsSource()

	// Write the document to local storage.
	var path = v.directory_ + "documents/" + name
	uti.MakeDirectory(path)
	var filename = path + "/" + version + ".bali"
	if uti.PathExists(filename) {
		status = rep.Existed
		return
	}
	var source = document.AsSource()
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
	var name = citation.GetTag().AsSource()[1:]
	var version = citation.GetVersion().AsSource()
	var path = v.directory_ + "documents/" + name
	var filename = path + "/" + version + ".bali"
	if !uti.PathExists(filename) {
		status = rep.Missing
		return
	}
	var source = uti.ReadFile(filename)
	document = not.Document(source)
	status = rep.Success
	return
}

func (v *localStorage_) DeleteDocument(
	citation not.CitationLike,
) (
	document not.DocumentLike,
	status rep.Status,
) {
	var name = citation.GetTag().AsSource()[1:]
	var version = citation.GetVersion().AsSource()
	var path = v.directory_ + "documents/" + name
	var filename = path + "/" + version + ".bali"
	if !uti.PathExists(filename) {
		status = rep.Missing
		return
	}
	var source = uti.ReadFile(filename)
	document = not.Document(source)
	uti.RemovePath(filename)
	var filenames = uti.ReadDirectory(path)
	if len(filenames) == 0 {
		uti.RemovePath(path)
	}
	status = rep.Success
	return
}

// PROTECTED INTERFACE

// Private Methods

func (v *localStorage_) extractPath(
	name string,
) string {
	var folders = sts.Split(name, "/")
	folders = folders[:len(folders)-1] // Strip off the last one.
	var path = sts.Join(folders, "/")
	return path
}

func (v *localStorage_) hashString(
	string_ string,
) string {
	var array = sha.Sum512([]byte(string_))
	var bytes = array[:20] // Convert the array to a slice of the first 20 bytes.
	var hash = doc.Tag(bytes).AsSource()[1:]
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
