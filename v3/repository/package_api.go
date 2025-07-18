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
	doc "github.com/bali-nebula/go-document-notation/v3"
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
	SaveCertificate(
		certificate not.ContractLike,
	) not.CitationLike
	RetrieveCertificate(
		citation not.CitationLike,
	) not.ContractLike
	SaveDraft(
		draft not.DraftLike,
	) not.CitationLike
	RetrieveDraft(
		draft not.CitationLike,
	) not.DraftLike
	DiscardDraft(
		draft not.CitationLike,
	)
	NotarizeDraft(
		name string,
		draft not.DraftLike,
	) not.ContractLike
	RetrieveContract(
		contract string,
	) not.ContractLike
	CheckoutDraft(
		contract string,
		level int,
	) not.DraftLike
	CreateBag(
		name string,
		permissions string,
		capacity int,
		leasetime int,
	)
	RemoveBag(
		bag string,
	)
	MessageCount(
		bag string,
	) int
	SendMessage(
		bag string,
		content doc.DocumentLike,
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
		kind string,
		content doc.DocumentLike,
		permissions string,
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
	) not.CitationLike
	WriteCitation(
		name fra.ResourceLike,
		citation not.CitationLike,
	)
	DeleteCitation(
		name fra.ResourceLike,
	)
	DraftExists(
		draft not.CitationLike,
	) bool
	ReadDraft(
		draft not.CitationLike,
	) not.DraftLike
	WriteDraft(
		draft not.DraftLike,
	) not.CitationLike
	DeleteDraft(
		draft not.CitationLike,
	)
	CertificateExists(
		certificate not.CitationLike,
	) bool
	ReadCertificate(
		certificate not.CitationLike,
	) not.ContractLike
	WriteCertificate(
		certificate not.ContractLike,
	) not.CitationLike
	ContractExists(
		contract not.CitationLike,
	) bool
	ReadContract(
		contract not.CitationLike,
	) not.ContractLike
	WriteContract(
		contract not.ContractLike,
	) not.CitationLike
	BagExists(
		bag not.CitationLike,
	) bool
	ReadBag(
		bag not.CitationLike,
	) not.ContractLike
	WriteBag(
		bag not.ContractLike,
	) not.CitationLike
	DeleteBag(
		bag not.CitationLike,
	)
	MessageCount(
		bag not.CitationLike,
	) int
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
	WriteEvent(
		event not.ContractLike,
	)
}
