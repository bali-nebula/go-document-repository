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
	bal "github.com/bali-nebula/go-bali-documents/v3"
	not "github.com/bali-nebula/go-digital-notary/v3"
	doc "github.com/bali-nebula/go-document-repository/v3/documents"
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
		certificate not.ContractLike,
	) bal.ResourceLike
	SaveDraft(
		draft not.Parameterized,
	) bal.ResourceLike
	RetrieveDraft(
		citation bal.ResourceLike,
	) not.Parameterized
	DiscardDraft(
		citation bal.ResourceLike,
	)
	NotarizeDocument(
		name bal.NameLike,
		version bal.VersionLike,
		draft not.Parameterized,
	) not.ContractLike
	RetrieveDocument(
		name bal.NameLike,
		version bal.VersionLike,
	) not.ContractLike
	CheckoutDocument(
		name bal.NameLike,
		version bal.VersionLike,
		level uint,
	) not.Parameterized
	CreateBag(
		name bal.NameLike,
		capacity uint,
		leasetime uint,
		permissions bal.ResourceLike,
	)
	RemoveBag(
		name bal.NameLike,
	)
	PostMessage(
		bag bal.NameLike,
		message doc.MessageLike,
	)
	RetrieveMessage(
		bag bal.NameLike,
	) not.ContractLike
	AcceptMessage(
		message not.ContractLike,
	)
	RejectMessage(
		message not.ContractLike,
	)
}

// ASPECT DECLARATIONS

/*
Persistent declares the set of method signatures that must be supported by all
persistent data storage mechanisms.
*/
type Persistent interface {
	ReadCitation(
		name bal.NameLike,
		version bal.VersionLike,
	) bal.ResourceLike
	WriteCitation(
		name bal.NameLike,
		version bal.VersionLike,
		citation bal.ResourceLike,
	)
	DeleteCitation(
		name bal.NameLike,
		version bal.VersionLike,
	)
	ListCitations(
		path bal.NameLike,
	) bal.Sequential[bal.ResourceLike]
	ReadDraft(
		citation bal.ResourceLike,
	) not.Parameterized
	WriteDraft(
		draft not.Parameterized,
	) bal.ResourceLike
	DeleteDraft(
		citation bal.ResourceLike,
	)
	ReadContract(
		citation bal.ResourceLike,
	) not.ContractLike
	WriteContract(
		contract not.ContractLike,
	) bal.ResourceLike
	DeleteContract(
		citation bal.ResourceLike,
	)
}
