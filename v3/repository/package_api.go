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
	fra "github.com/craterdog/go-component-framework/v7"
	uti "github.com/craterdog/go-missing-utilities/v7"
)

// TYPE DECLARATIONS

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
		certificate not.CertificateLike,
	) fra.ResourceLike
	SaveDraft(
		draft not.Parameterized,
	) fra.ResourceLike
	RetrieveDraft(
		citation fra.ResourceLike,
	) not.Parameterized
	DiscardDraft(
		citation fra.ResourceLike,
	)
	NotarizeDocument(
		name fra.ResourceLike,
		draft not.Parameterized,
	) not.Notarized
	RetrieveDocument(
		name fra.ResourceLike,
	) not.Notarized
	CheckoutDocument(
		name fra.ResourceLike,
		level uti.Cardinal,
	) not.Parameterized
	CreateBag(
		name fra.ResourceLike,
		permissions fra.ResourceLike,
		capacity uti.Cardinal,
		leasetime uti.Cardinal,
	)
	RemoveBag(
		name fra.ResourceLike,
	)
	MessageCount(
		bag fra.ResourceLike,
	) uti.Cardinal
	SendMessage(
		bag fra.ResourceLike,
		message doc.ItemsLike,
	)
	RetrieveMessage(
		bag fra.ResourceLike,
	) not.Notarized
	AcceptMessage(
		message not.Notarized,
	)
	RejectMessage(
		message not.Notarized,
	)
	PublishEvent(
		event doc.ItemsLike,
	)
}

// ASPECT DECLARATIONS

/*
Persistent declares the set of method signatures that must be supported by all
persistent data storage mechanisms.
*/
type Persistent interface {
	CitationExists(
		name fra.ResourceLike,
	) bool
	ReadCitation(
		name fra.ResourceLike,
	) fra.ResourceLike
	WriteCitation(
		name fra.ResourceLike,
		citation fra.ResourceLike,
	)
	DeleteCitation(
		name fra.ResourceLike,
	)
	DraftExists(
		citation fra.ResourceLike,
	) bool
	ReadDraft(
		citation fra.ResourceLike,
	) not.Parameterized
	WriteDraft(
		draft not.Parameterized,
	) fra.ResourceLike
	DeleteDraft(
		citation fra.ResourceLike,
	)
	DocumentExists(
		citation fra.ResourceLike,
	) bool
	ReadDocument(
		citation fra.ResourceLike,
	) not.Notarized
	WriteDocument(
		document not.Notarized,
	) fra.ResourceLike
	BagExists(
		citation fra.ResourceLike,
	) bool
	ReadBag(
		citation fra.ResourceLike,
	) not.Notarized
	WriteBag(
		bag not.Notarized,
	) fra.ResourceLike
	DeleteBag(
		citation fra.ResourceLike,
	)
	MessageCount(
		bag fra.ResourceLike,
	) uti.Cardinal
	ReadMessage(
		bag fra.ResourceLike,
	) not.Notarized
	WriteMessage(
		bag fra.ResourceLike,
		message not.Notarized,
	)
	DeleteMessage(
		bag fra.ResourceLike,
		citation fra.ResourceLike,
	)
	ReleaseMessage(
		bag fra.ResourceLike,
		citation fra.ResourceLike,
	)
	WriteEvent(
		event not.Notarized,
	)
}
