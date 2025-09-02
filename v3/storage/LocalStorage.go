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
	ran "crypto/rand"
	fmt "fmt"
	not "github.com/bali-nebula/go-digital-notary/v3"
	fra "github.com/craterdog/go-component-framework/v7"
	uti "github.com/craterdog/go-missing-utilities/v7"
	big "math/big"
	osx "os"
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

func (v *localStorage_) CitationExists(
	name fra.ResourceLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a document citation exists.",
	)

	// Determine whether or not the document citation exists.
	var path = v.directory_ + "citations"
	path += v.getNamePath(name)
	path += "/" + v.getNameFilename(name)
	return uti.PathExists(path)
}

func (v *localStorage_) ReadCitation(
	name fra.ResourceLike,
) fra.ResourceLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a document citation.",
	)

	// Read the document citation from local storage.
	var filename = v.directory_ + "citations"
	filename += v.getNamePath(name)
	filename += "/" + v.getNameFilename(name)
	var source = uti.ReadFile(filename)
	var citation = fra.ResourceFromString(source)
	return citation
}

func (v *localStorage_) WriteCitation(
	name fra.ResourceLike,
	citation fra.ResourceLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a document citation.",
	)

	// Write the document citation to local storage.
	var path = v.directory_ + "citations"
	path += v.getNamePath(name)
	uti.MakeDirectory(path)
	var filename = path + "/" + v.getNameFilename(name)
	var source = not.CitationFromResource(citation).AsString()
	uti.WriteFile(filename, source)
}

func (v *localStorage_) DeleteCitation(
	name fra.ResourceLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to delete a document citation.",
	)

	// Delete the document citation from local storage.
	var root = v.directory_ + "citations/"
	var path = v.getNamePath(name)
	var filename = root + path + "/" + v.getNameFilename(name)
	uti.RemovePath(filename)
	for len(path) > 0 {
		if len(uti.ReadDirectory(root+path)) > 0 {
			// The directory is not empty so we are done.
			return
		}
		uti.RemovePath(root + path)
		var directories = sts.Split(path, "/")
		directories = directories[:len(directories)-1] // Strip off last one.
		path = sts.Join(directories, "/")
	}
}

func (v *localStorage_) DraftExists(
	citation fra.ResourceLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a draft document exists.",
	)

	// Determine whether or not the draft document exists.
	var path = v.directory_ + "drafts/"
	path += v.getCitationTag(citation) + "/"
	path += v.getCitationVersion(citation) + ".bali"
	return uti.PathExists(path)
}

func (v *localStorage_) ReadDraft(
	citation fra.ResourceLike,
) not.Parameterized {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a draft document.",
	)

	// Read the draft document from local storage.
	var filename = v.directory_ + "drafts/"
	filename += v.getCitationTag(citation) + "/"
	filename += v.getCitationVersion(citation) + ".bali"
	var source = uti.ReadFile(filename)
	var draft = not.DraftFromString(source)
	return draft
}

func (v *localStorage_) WriteDraft(
	draft not.Parameterized,
) fra.ResourceLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a draft document.",
	)

	// Write the draft document to local storage.
	var citation = v.notary_.CiteDraft(draft)
	var path = v.directory_ + "drafts/"
	path += v.getCitationTag(citation) + "/"
	uti.MakeDirectory(path)
	var filename = path + v.getCitationVersion(citation) + ".bali"
	var source = draft.AsString()
	uti.WriteFile(filename, source)
	return citation
}

func (v *localStorage_) DeleteDraft(
	citation fra.ResourceLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to delete a draft document.",
	)

	// Delete the draft document from local storage.
	var path = v.directory_ + "drafts/" + v.getCitationTag(citation)
	var filename = path + "/" + v.getCitationVersion(citation) + ".bali"
	uti.RemovePath(filename)
	var filenames = uti.ReadDirectory(path)
	if len(filenames) == 0 {
		// This was the last version of the draft so delete the directory too.
		uti.RemovePath(path)
	}
}

func (v *localStorage_) DocumentExists(
	citation fra.ResourceLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a notarized document exists.",
	)

	// Determine whether or not the notarized document exists.
	var path = v.directory_ + "contracts/"
	path += v.getCitationTag(citation) + "/"
	path += v.getCitationVersion(citation) + ".bali"
	return uti.PathExists(path)
}

func (v *localStorage_) ReadDocument(
	citation fra.ResourceLike,
) not.Notarized {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a notarized document.",
	)

	// Read the notarized document from local storage.
	var filename = v.directory_ + "contracts/"
	filename += v.getCitationTag(citation) + "/"
	filename += v.getCitationVersion(citation) + ".bali"
	var source = uti.ReadFile(filename)
	var document = not.ContractFromString(source)
	return document
}

func (v *localStorage_) WriteDocument(
	document not.Notarized,
) fra.ResourceLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a notarized document.",
	)

	// Write the notarized document to local storage.
	var content = document.GetContent()
	var citation = v.notary_.CiteDraft(content)
	var path = v.directory_ + "contracts/"
	path += v.getCitationTag(citation) + "/"
	uti.MakeDirectory(path)
	var filename = path + v.getCitationVersion(citation) + ".bali"
	var source = document.AsString()
	uti.WriteFile(filename, source)
	return citation
}

func (v *localStorage_) BagExists(
	citation fra.ResourceLike,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while checking to see if a message bag exists.",
	)

	// Determine whether or not the message bag exists.
	var path = v.directory_ + "bags/"
	path += v.getCitationTag(citation) + ".bali"
	return uti.PathExists(path)
}

func (v *localStorage_) ReadBag(
	citation fra.ResourceLike,
) not.Notarized {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a message bag.",
	)

	// Read the message bag from local storage.
	var filename = v.directory_ + "bags/"
	filename += v.getCitationTag(citation) + ".bali"
	var source = uti.ReadFile(filename)
	var bag = not.ContractFromString(source)
	return bag
}

func (v *localStorage_) WriteBag(
	bag not.Notarized,
) fra.ResourceLike {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a message bag.",
	)

	// Create the bags directory.
	var path = v.directory_ + "bags/"
	uti.MakeDirectory(path)

	// Save the bag configuration file.
	var content = bag.GetContent()
	var citation = v.notary_.CiteDraft(content)
	var tag = v.getCitationTag(citation)
	var filename = path + tag + ".bali"
	var source = bag.AsString()
	uti.WriteFile(filename, source)

	// Create the messages directories for the bag.
	path = v.directory_ + "messages/" + tag
	uti.MakeDirectory(path + "/available")
	uti.MakeDirectory(path + "/processing")
	return citation
}

func (v *localStorage_) DeleteBag(
	citation fra.ResourceLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to delete a message bag.",
	)

	// Delete the messages directory for the bag.
	var tag = v.getCitationTag(citation)
	var path = v.directory_ + "messages/" + tag
	uti.RemovePath(path)

	// Delete the bag configuration file.
	path = v.directory_ + "bags/"
	var filename = path + "/" + tag + ".bali"
	uti.RemovePath(filename)
}

func (v *localStorage_) MessageCount(
	bag fra.ResourceLike,
) uti.Cardinal {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while counting the messages in a message bag.",
	)

	// Determine the number of messages currently available in the bag.
	var path = v.directory_ + "messages/" + v.getCitationTag(bag) + "/available"
	var messages = uti.ReadDirectory(path)
	return uti.Cardinal(len(messages))
}

func (v *localStorage_) ReadMessage(
	bag fra.ResourceLike,
) not.Notarized {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read a message from a message bag.",
	)

	// Read a random message from local storage.
	var message not.Notarized
	var path = v.directory_ + "messages/" + v.getCitationTag(bag)
	var available = path + "/available/"
	var processing = path + "/processing/"
	for tries := 5; tries > 0; tries-- {
		var filenames = uti.ReadDirectory(available)
		var count = len(filenames)
		if count == 0 {
			// No messages in the bag, try again...
			continue
		}
		var index, _ = ran.Int(ran.Reader, big.NewInt(int64(count)))
		var filename = "/" + filenames[index.Int64()]
		var err = osx.Rename(available+filename, processing+filename)
		if err != nil {
			// Someone got there first, try again...
			continue
		}
		var source = uti.ReadFile(processing + filename)
		message = not.ContractFromString(source)
		break
	}
	return message
}

func (v *localStorage_) WriteMessage(
	bag fra.ResourceLike,
	message not.Notarized,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write a message to a message bag.",
	)

	// Write the message to the message bag in local storage.
	var path = v.directory_ + "messages/" + v.getCitationTag(bag) + "/available"
	uti.MakeDirectory(path)
	var content = message.GetContent()
	var citation = v.notary_.CiteDraft(content)
	var filename = path + "/" + v.getCitationTag(citation) + ".bali"
	var source = message.AsString()
	uti.WriteFile(filename, source)
}

func (v *localStorage_) DeleteMessage(
	bag fra.ResourceLike,
	citation fra.ResourceLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to delete a message from a message bag.",
	)

	// Delete the message from the message bag in local storage.
	var path = v.directory_ + "messages/" + v.getCitationTag(bag) + "/processing"
	var filename = path + "/" + v.getCitationTag(citation) + ".bali"
	uti.RemovePath(filename)
}

func (v *localStorage_) ReleaseMessage(
	bag fra.ResourceLike,
	citation fra.ResourceLike,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to reset the lease on a message.",
	)

	// Reset the message lease for the message in local storage.
	var path = v.directory_ + "messages/" + v.getCitationTag(bag)
	var available = path + "/available/"
	var processing = path + "/processing/"
	var filename = v.getCitationTag(citation) + ".bali"
	// If Rename() fails the lease has already expired so nothing else to do.
	var _ = osx.Rename(processing+filename, available+filename)
}

func (v *localStorage_) WriteEvent(
	event not.Notarized,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write an event.",
	)

	// Write the event to the notification queue in local storage.
	var path = v.directory_ + "events/"
	uti.MakeDirectory(path)
	var content = event.GetContent()
	var citation = v.notary_.CiteDraft(content)
	var filename = path + v.getCitationTag(citation) + ".bali"
	var source = event.AsString()
	uti.WriteFile(filename, source)
}

// PROTECTED INTERFACE

// Private Methods

func (v *localStorage_) errorCheck(
	message string,
) {
	if e := recover(); e != nil {
		message = fmt.Sprintf(
			"LocalStorage: %s:\n    %v",
			message,
			e,
		)
		panic(message)
	}
}

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
	resource fra.ResourceLike,
) string {
	return sts.Split(resource.GetPath(), ":")[0]
}

func (v *localStorage_) getNameFilename(
	resource fra.ResourceLike,
) string {
	return sts.Split(resource.GetPath(), ":")[1] + ".bali"
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
