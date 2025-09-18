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
	name fra.NameLike,
	version fra.VersionLike,
) fra.ResourceLike {
	var filename = v.getNamePath(name) + v.getVersionFilename(version)
	var source = uti.ReadFile(filename)
	var citation = not.CitationFromString(source)
	return citation.AsResource()
}

func (v *localStorage_) WriteCitation(
	name fra.NameLike,
	version fra.VersionLike,
	citation fra.ResourceLike,
) {
	var path = v.getNamePath(name)
	uti.MakeDirectory(path)
	var filename = path + v.getVersionFilename(version)
	var source = not.CitationFromResource(citation).AsString()
	uti.WriteFile(filename, source)
}

func (v *localStorage_) DeleteCitation(
	name fra.NameLike,
	version fra.VersionLike,
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
		directories = directories[:len(directories)-1] // Strip off the last one.
		path = sts.Join(directories, "/")
	}
}

func (v *localStorage_) ListCitations(
	path fra.NameLike,
) fra.Sequential[fra.ResourceLike] {
	var citations = fra.List[fra.ResourceLike]()
	var directory = v.getNamePath(path)
	var filenames = uti.ReadDirectory(directory)
	for _, filename := range filenames {
		var source = uti.ReadFile(directory + filename + "/v1.bali")
		var citation = not.CitationFromString(source)
		citations.AppendValue(citation.AsResource())
	}
	return citations
}

func (v *localStorage_) ReadDraft(
	citation fra.ResourceLike,
) not.Parameterized {
	var path = v.directory_ + "nebula/" + v.getCitationTag(citation) + "/"
	var filename = path + v.getCitationVersion(citation) + ".bali"
	var source = uti.ReadFile(filename)
	var draft = not.DraftFromString(source)
	return draft
}

func (v *localStorage_) WriteDraft(
	draft not.Parameterized,
) fra.ResourceLike {
	var citation = v.notary_.CiteDraft(draft)
	var path = v.directory_ + "nebula/" + v.getCitationTag(citation) + "/"
	uti.MakeDirectory(path)
	var filename = path + v.getCitationVersion(citation) + ".bali"
	var source = draft.AsString()
	uti.WriteFile(filename, source)
	return citation
}

func (v *localStorage_) DeleteDraft(
	citation fra.ResourceLike,
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
	citation fra.ResourceLike,
) not.Notarized {
	var path = v.directory_ + "nebula/" + v.getCitationTag(citation) + "/"
	var filename = path + v.getCitationVersion(citation) + ".bali"
	var source = uti.ReadFile(filename)
	var document = not.ContractFromString(source)
	return document
}

func (v *localStorage_) WriteContract(
	contract not.Notarized,
) fra.ResourceLike {
	var content = contract.GetContent()
	var citation = v.notary_.CiteDraft(content)
	var path = v.directory_ + "nebula/" + v.getCitationTag(citation) + "/"
	uti.MakeDirectory(path)
	var filename = path + v.getCitationVersion(citation) + ".bali"
	var source = contract.AsString()
	uti.WriteFile(filename, source)
	return citation
}

func (v *localStorage_) DeleteContract(
	citation fra.ResourceLike,
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
	resource fra.ResourceLike,
) string {
	var citation = not.CitationFromResource(resource)
	var tag = citation.GetTag()
	return tag.AsString()[1:] // Remove the leading "#" character.
}

func (v *localStorage_) getCitationVersion(
	resource fra.ResourceLike,
) string {
	var citation = not.CitationFromResource(resource)
	var version = citation.GetVersion()
	return version.AsString()
}

func (v *localStorage_) getNamePath(
	name fra.NameLike,
) string {
	return v.directory_ + "bali" + name.AsString() + "/"
}

func (v *localStorage_) getVersionFilename(
	version fra.VersionLike,
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
