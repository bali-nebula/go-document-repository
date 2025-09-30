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

package documents

import (
	doc "github.com/bali-nebula/go-bali-documents/v3"
	uti "github.com/craterdog/go-missing-utilities/v7"
)

// CLASS INTERFACE

// Access Function

func MessageClass() MessageClassLike {
	return messageClass()
}

// Constructor Methods

func (c *messageClass_) Message(
	entity any,
	type_ doc.ResourceLike,
	permissions doc.ResourceLike,
) MessageLike {
	if uti.IsUndefined(entity) {
		panic("The \"entity\" attribute is required by this class.")
	}
	if uti.IsUndefined(type_) {
		panic("The \"type\" attribute is required by this class.")
	}
	if uti.IsUndefined(permissions) {
		panic("The \"permissions\" attribute is required by this class.")
	}
	var tag = doc.Tag()
	var version = doc.Version()
	var component = doc.ParseSource(
		doc.FormatComponent(entity) + `(
    $type: ` + type_.AsString() + `
    $tag: ` + tag.AsString() + `
    $version: ` + version.AsString() + `
    $permissions: ` + permissions.AsString() + `
    $previous: none
)`,
	)
	var instance = &message_{
		// Initialize the instance attributes.

		// Initialize the inherited aspects.
		Declarative: component,
	}
	return instance
}

func (c *messageClass_) MessageFromString(
	source string,
) MessageLike {
	var component = doc.ParseSource(source)
	var instance = &message_{
		// Initialize the instance attributes.

		// Initialize the inherited aspects.
		Declarative: component,
	}
	return instance
}

// Constant Methods

// Function Methods

// INSTANCE INTERFACE

// Principal Methods

func (v *message_) GetClass() MessageClassLike {
	return messageClass()
}

func (v *message_) AsString() string {
	return doc.FormatDocument(v.Declarative.(doc.ComponentLike))
}

func (v *message_) AsIntrinsic() doc.ComponentLike {
	return v.Declarative.(doc.ComponentLike)
}

// Attribute Methods

func (v *message_) GetBag() doc.NameLike {
	var component = v.GetObject(doc.Symbol("$bag")).GetComponent()
	var bag = component.GetEntity().(doc.NameLike)
	return bag
}

func (v *message_) SetBag(
	bag doc.NameLike,
) {
	v.SetObject(bag, doc.Symbol("$bag"))
}

// Parameterized Methods

func (v *message_) GetEntity() any {
	return v.Declarative.(doc.ComponentLike).GetEntity()
}

func (v *message_) GetType() doc.ResourceLike {
	var component = v.GetParameter(doc.Symbol("$type"))
	return doc.Resource(doc.FormatComponent(component))
}

func (v *message_) GetTag() doc.TagLike {
	var component = v.GetParameter(doc.Symbol("$tag"))
	return doc.Tag(doc.FormatComponent(component))
}

func (v *message_) GetVersion() doc.VersionLike {
	var component = v.GetParameter(doc.Symbol("$version"))
	return doc.Version(doc.FormatComponent(component))
}

func (v *message_) GetOptionalPrevious() doc.ResourceLike {
	var previous doc.ResourceLike
	var component = v.GetParameter(doc.Symbol("$previous"))
	var source = doc.FormatComponent(component)
	if source != "none" {
		previous = doc.Resource(source)
	}
	return previous
}

func (v *message_) GetPermissions() doc.ResourceLike {
	var component = v.GetParameter(doc.Symbol("$permissions"))
	return doc.Resource(doc.FormatComponent(component))
}

func (v *message_) GetAccount() doc.TagLike {
	var component = v.GetParameter(doc.Symbol("$account"))
	return doc.Tag(doc.FormatComponent(component))
}

// PROTECTED INTERFACE

// Private Methods

// Instance Structure

type message_ struct {
	// Declare the instance attributes.

	// Declare the inherited aspects.
	doc.Declarative
}

// Class Structure

type messageClass_ struct {
	// Declare the class constants.
}

// Class Reference

func messageClass() *messageClass_ {
	return messageClassReference_
}

var messageClassReference_ = &messageClass_{
	// Initialize the class constants.
}
