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

func (v *localStorage_) ReadCitation(
	name doc.NameLike,
	version doc.VersionLike,
) doc.ResourceLike {
	var filename = v.getNamePath(name) + v.getVersionFilename(version)
	var source = uti.ReadFile(filename)
	var citation = not.Citation(source)
	return citation.AsResource()
}

func (v *localStorage_) WriteCitation(
	name doc.NameLike,
	version doc.VersionLike,
	citation doc.ResourceLike,
) {
	var path = v.getNamePath(name)
	uti.MakeDirectory(path)
	var filename = path + v.getVersionFilename(version)
	var source = not.Citation(citation).AsString()
	uti.WriteFile(filename, source)
}

func (v *localStorage_) DeleteCitation(
	name doc.NameLike,
	version doc.VersionLike,
) {
	// Remove the citation file.
	var filename = v.getNamePath(name) + v.getVersionFilename(version)
	uti.RemovePath(filename)

	// Remove any empty directories in the citation path.
	var path = v.getNamePath(name)
	for len(path) > 0 {
		if len(uti.ReadDirectory(path)) > 0 {
			// The directory is not empty so we are done.
			return
		}
		uti.RemovePath(path)
		var directories = sts.Split(path, "/")
		directories = directories[:len(directories)-2] // Strip off the last one.
		path = sts.Join(directories, "/")
	}
}

func (v *localStorage_) ListCitations(
	path doc.NameLike,
) doc.Sequential[doc.ResourceLike] {
	var citations = doc.List[doc.ResourceLike]()
	var directory = v.getNamePath(path)
	var filenames = uti.ReadDirectory(directory)
	for _, filename := range filenames {
		var source = uti.ReadFile(directory + filename + "/v1.bali")
		var citation = not.Citation(source)
		citations.AppendValue(citation.AsResource())
	}
	return citations
}

func (v *localStorage_) ReadDraft(
	citation doc.ResourceLike,
) not.Parameterized {
	var path = v.directory_ + "nebula/" + v.getCitationTag(citation) + "/"
	var filename = path + v.getCitationVersion(citation) + ".bali"
	var source = uti.ReadFile(filename)
	var draft = not.Draft(source)
	return draft
}

func (v *localStorage_) WriteDraft(
	draft not.Parameterized,
) doc.ResourceLike {
	var citation = v.notary_.CiteDocument(draft)
	var path = v.directory_ + "nebula/" + v.getCitationTag(citation) + "/"
	uti.MakeDirectory(path)
	var filename = path + v.getCitationVersion(citation) + ".bali"
	var source = draft.AsString()
	uti.WriteFile(filename, source)
	return citation
}

func (v *localStorage_) DeleteDraft(
	citation doc.ResourceLike,
) {
	var path = v.directory_ + "nebula/" + v.getCitationTag(citation) + "/"
	var filename = path + v.getCitationVersion(citation) + ".bali"
	uti.RemovePath(filename)
	var filenames = uti.ReadDirectory(path)
	if len(filenames) == 0 {
		// This was the last version of the draft so delete the directory too.
		uti.RemovePath(path)
	}
}

func (v *localStorage_) ReadContract(
	citation doc.ResourceLike,
) not.ContractLike {
	var path = v.directory_ + "nebula/" + v.getCitationTag(citation) + "/"
	var filename = path + v.getCitationVersion(citation) + ".bali"
	var source = uti.ReadFile(filename)
	var document = not.Contract(source)
	return document
}

func (v *localStorage_) WriteContract(
	contract not.ContractLike,
) doc.ResourceLike {
	var content = contract.GetContent()
	var citation = v.notary_.CiteDocument(content)
	var path = v.directory_ + "nebula/" + v.getCitationTag(citation) + "/"
	uti.MakeDirectory(path)
	var filename = path + v.getCitationVersion(citation) + ".bali"
	var source = contract.AsString()
	uti.WriteFile(filename, source)
	return citation
}

func (v *localStorage_) DeleteContract(
	citation doc.ResourceLike,
) {
	var path = v.directory_ + "nebula/" + v.getCitationTag(citation) + "/"
	var filename = path + v.getCitationVersion(citation) + ".bali"
	uti.RemovePath(filename)
	var filenames = uti.ReadDirectory(path)
	if len(filenames) == 0 {
		// This was the last version of the document so delete the directory too.
		uti.RemovePath(path)
	}
}

// PROTECTED INTERFACE

// Private Methods

func (v *localStorage_) getCitationTag(
	resource doc.ResourceLike,
) string {
	var citation = not.Citation(resource)
	var tag = citation.GetTag()
	return tag.AsString()[1:] // Remove the leading "#" character.
}

func (v *localStorage_) getCitationVersion(
	resource doc.ResourceLike,
) string {
	var citation = not.Citation(resource)
	var version = citation.GetVersion()
	return version.AsString()
}

func (v *localStorage_) getNamePath(
	name doc.NameLike,
) string {
	return v.directory_ + "bali" + name.AsString() + "/"
}

func (v *localStorage_) getVersionFilename(
	version doc.VersionLike,
) string {
	return version.AsString() + ".bali"
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
