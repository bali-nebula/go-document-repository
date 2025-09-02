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
Package "documents" provides an implementation of wrappers for various types of
Bali Document Notation™ documents that are required by the document repository.

For detailed documentation on this package refer to the wiki:
  - https://github.com/bali-nebula/go-document-repository/wiki

This package follows the Crater Dog Technologies™ Go Coding Conventions located
here:
  - https://github.com/craterdog/go-development-tools/wiki/Coding-Conventions

Additional concrete implementations of the classes declared by this package can
be developed and used seamlessly since the interface declarations only depend on
other interfaces and intrinsic types—and the class implementations only depend
on interfaces, not on each other.
*/
package documents

import (
	doc "github.com/bali-nebula/go-bali-documents/v3"
	not "github.com/bali-nebula/go-digital-notary/v3"
	fra "github.com/craterdog/go-component-framework/v7"
)

// TYPE DECLARATIONS

// FUNCTIONAL DECLARATIONS

// CLASS DECLARATIONS

/*
EventClassLike is a class interface that declares the complete set of
class constructors, constants and functions that must be supported by each
concrete event-like class.
*/
type EventClassLike interface {
	// Constructor Methods
	Event(
		entity any,
		type_ fra.ResourceLike,
		permissions fra.ResourceLike,
	) EventLike
	EventFromString(
		source string,
	) EventLike
}

/*
MessageClassLike is a class interface that declares the complete set of
class constructors, constants and functions that must be supported by each
concrete message-like class.
*/
type MessageClassLike interface {
	// Constructor Methods
	Message(
		entity any,
		type_ fra.ResourceLike,
		permissions fra.ResourceLike,
	) MessageLike
	MessageFromString(
		source string,
	) MessageLike
}

// INSTANCE DECLARATIONS

/*
EventLike is an instance interface that declares the complete set of
principal, attribute and aspect methods that must be supported by each instance
of a concrete event-like class.
*/
type EventLike interface {
	// Principal Methods
	GetClass() EventClassLike
	AsString() string
	AsIntrinsic() doc.ComponentLike
	GetKind() fra.ResourceLike

	// Aspect Interfaces
	doc.Declarative
	not.Parameterized
}

/*
MessageLike is an instance interface that declares the complete set of
principal, attribute and aspect methods that must be supported by each instance
of a concrete message-like class.
*/
type MessageLike interface {
	// Principal Methods
	GetClass() MessageClassLike
	AsString() string
	AsIntrinsic() doc.ComponentLike
	GetBag() fra.ResourceLike

	// Aspect Interfaces
	doc.Declarative
	not.Parameterized
}
