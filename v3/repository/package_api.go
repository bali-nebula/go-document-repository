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
	fra "github.com/craterdog/go-component-framework/v7"
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
	SaveDraft(
		draft not.DraftLike,
	) not.CitationLike
	RetrieveDraft(
		citation not.CitationLike,
	) not.DraftLike
	DiscardDraft(
		citation not.CitationLike,
	)
	NotarizeDraft(
		name fra.NameLike,
		draft not.DraftLike,
	) not.ContractLike
	RetrieveContract(
		name fra.NameLike,
	) not.ContractLike
	CheckoutDraft(
		name fra.NameLike,
		level uint,
	) not.DraftLike
	CreateBag(
		name fra.NameLike,
		permissions fra.ResourceLike,
		capacity uint,
		leasetime uint,
	)
	MessageCount(
		bag fra.NameLike,
	) uint
	PostMessage(
		bag fra.NameLike,
		message not.DraftLike,
	)
	RetrieveMessage(
		bag fra.NameLike,
	) not.DraftLike
	AcceptMessage(
		message not.DraftLike,
	)
	RejectMessage(
		message not.DraftLike,
	)
	DeleteBag(
		bag fra.NameLike,
	)
	PublishEvent(
		event not.DraftLike,
	)
}

// ASPECT DECLARATIONS

/*
Persistent declares the set of method signatures that must be supported by all
persistent data storage mechanisms.
*/
type Persistent interface {
	CitationExists(
		name fra.NameLike,
	) bool
	ReadCitation(
		name fra.NameLike,
	) not.CitationLike
	WriteCitation(
		name fra.NameLike,
		citation not.CitationLike,
	)
	DraftExists(
		citation not.CitationLike,
	) bool
	ReadDraft(
		citation not.CitationLike,
	) not.DraftLike
	WriteDraft(
		draft not.DraftLike,
	) not.CitationLike
	DeleteDraft(
		citation not.CitationLike,
	)
	ContractExists(
		citation not.CitationLike,
	) bool
	ReadContract(
		citation not.CitationLike,
	) not.ContractLike
	WriteContract(
		contract not.ContractLike,
	) not.CitationLike
	MessageCount(
		bag not.CitationLike,
	) uint
	ReadMessage(
		bag not.CitationLike,
	) not.ContractLike
	WriteMessage(
		bag not.CitationLike,
		message not.ContractLike,
	)
	DeleteMessage(
		bag not.CitationLike,
		message not.CitationLike,
	)
	ReleaseMessage(
		bag not.CitationLike,
		message not.CitationLike,
	)
}
