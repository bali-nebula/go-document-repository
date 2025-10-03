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

/*
Package "repository" provides an implementation of a document repository for
managing documents formatted using Bali Document Notation™ (Bali).

For detailed documentation on this package refer to the wiki:
  - https://github.com/bali-nebula/go-document-repository/wiki

Detailed information on Bali Document Notation™ can be found here:
  - https://github.com/bali-nebula/go-document-notation/wiki

This package follows the Crater Dog Technologies™ Go Coding Conventions located
here:
  - https://github.com/craterdog/go-development-tools/wiki/Coding-Conventions

Additional concrete implementations of the classes declared by this package can
be developed and used seamlessly since the interface declarations only depend on
other interfaces and intrinsic types—and the class implementations only depend
on interfaces, not on each other.
*/
package repository

import (
	doc "github.com/bali-nebula/go-bali-documents/v3"
	not "github.com/bali-nebula/go-digital-notary/v3"
)

// TYPE DECLARATIONS

/*
Status is a constrained type specifying the result of a storage operation.
*/
type Status uint8

const (
	Problem Status = iota
	Success
	Missing
	Existed
	Illegal
	Invalid
)

// FUNCTIONAL DECLARATIONS

// CLASS DECLARATIONS

/*
DocumentRepositoryClassLike is a class interface that declares the complete set of
class constructors, constants and functions that must be supported by each
concrete document-repository-like class.
*/
type DocumentRepositoryClassLike interface {
	// Constructor Methods
	DocumentRepository(
		notary not.DigitalNotaryLike,
		storage Persistent,
	) DocumentRepositoryLike
}

// INSTANCE DECLARATIONS

/*
DocumentRepositoryLike is an instance interface that declares the complete set of
principal, attribute and aspect methods that must be supported by each instance
of a concrete document-repository-like class.
*/
type DocumentRepositoryLike interface {
	// Principal Methods
	GetClass() DocumentRepositoryClassLike
	SaveCertificate(
		certificate not.DocumentLike,
	) (
		citation not.CitationLike,
		status Status,
	)
	SaveDraft(
		draft not.DocumentLike,
	) (
		citation not.CitationLike,
		status Status,
	)
	RetrieveDraft(
		citation not.CitationLike,
	) (
		draft not.DocumentLike,
		status Status,
	)
	DiscardDraft(
		citation not.CitationLike,
	) (
		draft not.DocumentLike,
		status Status,
	)
	NotarizeDocument(
		name doc.NameLike,
		version doc.VersionLike,
		document not.DocumentLike,
	) (
		status Status,
	)
	RetrieveDocument(
		name doc.NameLike,
		version doc.VersionLike,
	) (
		document not.DocumentLike,
		status Status,
	)
	CheckoutDocument(
		name doc.NameLike,
		version doc.VersionLike,
		level uint,
	) (
		document not.DocumentLike,
		status Status,
	)
	PostMessage(
		bag doc.NameLike,
		message not.DocumentLike,
	) (
		status Status,
	)
	RetrieveMessage(
		bag doc.NameLike,
	) (
		message not.DocumentLike,
		status Status,
	)
	AcceptMessage(
		bag doc.NameLike,
		message not.DocumentLike,
	) (
		status Status,
	)
	RejectMessage(
		bag doc.NameLike,
		message not.DocumentLike,
	) (
		status Status,
	)
	SubscribeEvents(
		bag doc.NameLike,
		type_ doc.ResourceLike,
	) (
		status Status,
	)
	UnsubscribeEvents(
		bag doc.NameLike,
		type_ doc.ResourceLike,
	) (
		status Status,
	)
	PublishEvent(
		event not.DocumentLike,
	) (
		status Status,
	)
}

// ASPECT DECLARATIONS

/*
Persistent declares the set of method signatures that must be supported by all
persistent data storage mechanisms.
*/
type Persistent interface {
	WriteCitation(
		name doc.NameLike,
		version doc.VersionLike,
		citation not.CitationLike,
	) (
		status Status,
	)
	ReadCitation(
		name doc.NameLike,
		version doc.VersionLike,
	) (
		citation not.CitationLike,
		status Status,
	)
	DeleteCitation(
		name doc.NameLike,
		version doc.VersionLike,
	) (
		citation not.CitationLike,
		status Status,
	)
	WriteMessage(
		bag doc.NameLike,
		message not.CitationLike,
	) (
		status Status,
	)
	ReadMessage(
		bag doc.NameLike,
	) (
		message not.CitationLike,
		status Status,
	)
	UnreadMessage(
		bag doc.NameLike,
		message not.CitationLike,
	) (
		status Status,
	)
	DeleteMessage(
		bag doc.NameLike,
		message not.CitationLike,
	) (
		status Status,
	)
	WriteSubscription(
		bag doc.NameLike,
		type_ doc.ResourceLike,
	) (
		status Status,
	)
	DeleteSubscription(
		bag doc.NameLike,
		type_ doc.ResourceLike,
	) (
		status Status,
	)
	WriteEvent(
		event not.DocumentLike,
	) (
		status Status,
	)
	WriteDraft(
		draft not.DocumentLike,
	) (
		citation not.CitationLike,
		status Status,
	)
	ReadDraft(
		citation not.CitationLike,
	) (
		draft not.DocumentLike,
		status Status,
	)
	DeleteDraft(
		citation not.CitationLike,
	) (
		draft not.DocumentLike,
		status Status,
	)
	WriteDocument(
		document not.DocumentLike,
	) (
		citation not.CitationLike,
		status Status,
	)
	ReadDocument(
		citation not.CitationLike,
	) (
		document not.DocumentLike,
		status Status,
	)
	DeleteDocument(
		citation not.CitationLike,
	) (
		document not.DocumentLike,
		status Status,
	)
}
