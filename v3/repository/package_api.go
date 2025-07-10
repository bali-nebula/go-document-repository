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
		notary not.NotaryLike,
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
	SaveDocument(
		document not.DocumentLike,
	) not.CitationLike
	RetrieveDocument(
		citation not.CitationLike,
	) not.DocumentLike
	DiscardDocument(
		citation not.CitationLike,
	) bool
	NotarizeDocument(
		name string,
		document not.DocumentLike,
	) not.ContractLike
	RetrieveContract(
		name string,
	) not.ContractLike
	CheckoutDocument(
		name string,
		level uint,
	) not.DocumentLike
	CreateBag(
		name string,
		permissions string,
		capacity uint,
		lease uint,
	)
	MessageCount(
		bag string,
	) uint
	PostMessage(
		bag string,
		message not.DocumentLike,
	)
	RetrieveMessage(
		bag string,
	) not.ContractLike
	AcceptMessage(
		message not.ContractLike,
	)
	RejectMessage(
		message not.ContractLike,
	)
	PublishEvent(
		event not.DocumentLike,
	)
}

// ASPECT DECLARATIONS

/*
Persistent declares the set of method signatures that must be supported by all
persistent data storage mechanisms.
*/
type Persistent interface {
	CitationExists(
		name string,
	) bool
	ReadCitation(
		name string,
	) string
	WriteCitation(
		name string,
		citation string,
	)
	DocumentExists(
		citation string,
	) bool
	ReadDocument(
		citation string,
	) string
	WriteDocument(
		document string,
	) string
	DeleteDocument(
		citation string,
	) string
	ContractExists(
		citation string,
	) bool
	ReadContract(
		citation string,
	) string
	WriteContract(
		contract string,
	) string
	MessageAvailable(
		bag string,
	) bool
	MessageCount(
		bag string,
	) uint
	AddMessage(
		bag string,
		message string,
	) string
	RetrieveMessage(
		bag string,
	) string
	ReturnMessage(
		message string,
	)
	DeleteMessage(
		message string,
	)
}
