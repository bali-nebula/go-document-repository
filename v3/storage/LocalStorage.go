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
	doc "github.com/bali-nebula/go-bali-documents/v3"
	not "github.com/bali-nebula/go-digital-notary/v3"
	rep "github.com/bali-nebula/go-document-repository/v3/repository"
	uti "github.com/craterdog/go-missing-utilities/v7"
	log "log"
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
	defer v.errorCheck(fmt.Sprintf(
		"An error occurred while attempting to write a citation: %s:%s:%s",
		name,
		version,
		citation,
	), &status)
	var path = v.getNamePath(name)
	uti.MakeDirectory(path)
	var filename = path + v.getVersionFilename(version)
	var source = not.Citation(citation).AsString()
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
	defer v.errorCheck(fmt.Sprintf(
		"An error occurred while attempting to read a citation: %s:%s",
		name,
		version,
	), &status)
	var filename = v.getNamePath(name) + v.getVersionFilename(version)
	var source = uti.ReadFile(filename)
	citation = not.Citation(source).AsResource()
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
	defer v.errorCheck(fmt.Sprintf(
		"An error occurred while attempting to delete a citation: %s:%s",
		name,
		version,
	), &status)
	// Remove the citation file.
	var filename = v.getNamePath(name) + v.getVersionFilename(version)
	uti.RemovePath(filename)
	status = rep.Success

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
	return
}

func (v *localStorage_) ListCitations(
	path doc.NameLike,
) (
	citations doc.Sequential[not.CitationLike],
	status rep.Status,
) {
	defer v.errorCheck(fmt.Sprintf(
		"An error occurred while attempting to list citations: %s",
		path,
	), &status)
	var list = doc.List[not.CitationLike]()
	var directory = v.getNamePath(path)
	var filenames = uti.ReadDirectory(directory)
	for _, filename := range filenames {
		var source = uti.ReadFile(directory + filename + "/v1.bali")
		var citation = not.Citation(source)
		list.AppendValue(citation.AsResource())
	}
	citations = list
	status = rep.Success
	return
}

func (v *localStorage_) WriteDocument(
	document not.DocumentLike,
) (
	citation not.CitationLike,
	status rep.Status,
) {
	defer v.errorCheck(fmt.Sprintf(
		"An error occurred while attempting to write a document: %s",
		document,
	), &status)
	var content = document.GetContent()
	citation = v.notary_.CiteDocument(content)
	var path = v.directory_ + "nebula/" + v.getCitationTag(citation) + "/"
	uti.MakeDirectory(path)
	var filename = path + v.getCitationVersion(citation) + ".bali"
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
	defer v.errorCheck(fmt.Sprintf(
		"An error occurred while attempting to read a document: %s",
		citation,
	), &status)
	var path = v.directory_ + "nebula/" + v.getCitationTag(citation) + "/"
	var filename = path + v.getCitationVersion(citation) + ".bali"
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
	defer v.errorCheck(fmt.Sprintf(
		"An error occurred while attempting to delete a document: %s",
		citation,
	), &status)
	var path = v.directory_ + "nebula/" + v.getCitationTag(citation) + "/"
	var filename = path + v.getCitationVersion(citation) + ".bali"
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

func (v *localStorage_) errorCheck(
	message string,
	pstatus *rep.Status,
) {
	if e := recover(); e != nil {
		log.Printf(
			"LocalStorage: %s:\n    %s\n",
			message,
			e,
		)
		*pstatus = rep.Unavailable
	}
}

func (v *localStorage_) getCitationTag(
	resource not.CitationLike,
) string {
	var citation = not.Citation(resource)
	var tag = citation.GetTag()
	return tag.AsString()[1:] // Remove the leading "#" character.
}

func (v *localStorage_) getCitationVersion(
	resource not.CitationLike,
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
