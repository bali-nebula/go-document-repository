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
	not "github.com/bali-nebula/go-digital-notary/v3"
	doc "github.com/bali-nebula/go-document-repository/v3/documents"
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
		name fra.NameLike,
		version fra.VersionLike,
		draft not.Parameterized,
	) not.Notarized
	RetrieveDocument(
		name fra.NameLike,
		version fra.VersionLike,
	) not.Notarized
	CheckoutDocument(
		name fra.NameLike,
		version fra.VersionLike,
		level uti.Cardinal,
	) not.Parameterized
	CreateBag(
		name fra.NameLike,
		capacity uti.Cardinal,
		leasetime uti.Cardinal,
		permissions fra.ResourceLike,
	)
	RemoveBag(
		name fra.NameLike,
	)
	PostMessage(
		bag fra.NameLike,
		message doc.MessageLike,
	)
	RetrieveMessage(
		bag fra.NameLike,
	) not.Notarized
	AcceptMessage(
		message not.Notarized,
	)
	RejectMessage(
		message not.Notarized,
	)
}

// ASPECT DECLARATIONS

/*
Persistent declares the set of method signatures that must be supported by all
persistent data storage mechanisms.
*/
type Persistent interface {
	ReadCitation(
		name fra.NameLike,
		version fra.VersionLike,
	) fra.ResourceLike
	WriteCitation(
		name fra.NameLike,
		version fra.VersionLike,
		citation fra.ResourceLike,
	)
	DeleteCitation(
		name fra.NameLike,
		version fra.VersionLike,
	)
	ListCitations(
		path fra.NameLike,
	) fra.Sequential[fra.ResourceLike]
	ReadDraft(
		citation fra.ResourceLike,
	) not.Parameterized
	WriteDraft(
		draft not.Parameterized,
	) fra.ResourceLike
	DeleteDraft(
		citation fra.ResourceLike,
	)
	ReadContract(
		citation fra.ResourceLike,
	) not.Notarized
	WriteContract(
		contract not.Notarized,
	) fra.ResourceLike
	DeleteContract(
		citation fra.ResourceLike,
	)
}
