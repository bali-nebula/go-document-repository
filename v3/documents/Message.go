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
	fra "github.com/craterdog/go-component-framework/v7"
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
	type_ fra.ResourceLike,
	permissions fra.ResourceLike,
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
	var tag = fra.TagWithSize(20)
	var component = doc.ParseSource(
		doc.FormatComponent(entity) + `(
    $type: ` + type_.AsString() + `
    $tag: ` + tag.AsString() + `
    $version: v1
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

func (v *message_) GetBag() fra.NameLike {
	var component = v.GetObject(fra.Symbol("bag")).GetComponent()
	var bag = component.GetEntity().(fra.NameLike)
	return bag
}

func (v *message_) SetBag(
	bag fra.NameLike,
) {
	v.SetObject(bag, fra.Symbol("bag"))
}

// Parameterized Methods

func (v *message_) GetEntity() any {
	return v.Declarative.(doc.ComponentLike).GetEntity()
}

func (v *message_) GetType() fra.ResourceLike {
	var component = v.GetParameter(fra.Symbol("type"))
	return fra.ResourceFromString(doc.FormatComponent(component))
}

func (v *message_) GetTag() fra.TagLike {
	var component = v.GetParameter(fra.Symbol("tag"))
	return fra.TagFromString(doc.FormatComponent(component))
}

func (v *message_) GetVersion() fra.VersionLike {
	var component = v.GetParameter(fra.Symbol("version"))
	return fra.VersionFromString(doc.FormatComponent(component))
}

func (v *message_) GetPermissions() fra.ResourceLike {
	var component = v.GetParameter(fra.Symbol("permissions"))
	return fra.ResourceFromString(doc.FormatComponent(component))
}

func (v *message_) GetOptionalPrevious() fra.ResourceLike {
	var previous fra.ResourceLike
	var component = v.GetParameter(fra.Symbol("previous"))
	var source = doc.FormatComponent(component)
	if source != "none" {
		previous = fra.ResourceFromString(source)
	}
	return previous
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
