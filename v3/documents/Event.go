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

func EventClass() EventClassLike {
	return eventClass()
}

// Constructor Methods

func (c *eventClass_) Event(
	entity any,
	type_ fra.ResourceLike,
	permissions fra.ResourceLike,
) EventLike {
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
	var event = doc.ParseSource(
		doc.FormatComponent(entity) + `(
    $type: ` + type_.AsString() + `
    $tag: ` + tag.AsString() + `
    $version: v1
    $permissions: ` + permissions.AsString() + `
    $previous: none
)`,
	)
	var object = event.GetObject(fra.Symbol("kind"))
	if uti.IsUndefined(object) {
		panic("A \"kind\" attribute is required by this class.")
	}
	var component = object.GetComponent()
	switch component.GetEntity().(type) {
	case fra.ResourceLike:
	default:
		panic("The \"kind\" attribute must be a named resource.")
	}

	var instance = &event_{
		// Initialize the instance attributes.

		// Initialize the inherited aspects.
		Declarative: component,
	}
	return instance
}

func (c *eventClass_) EventFromString(
	source string,
) EventLike {
	var component = doc.ParseSource(source)
	var instance = &event_{
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

func (v *event_) GetClass() EventClassLike {
	return eventClass()
}

func (v *event_) AsString() string {
	return doc.FormatDocument(v.Declarative.(doc.ComponentLike))
}

func (v *event_) AsIntrinsic() doc.ComponentLike {
	return v.Declarative.(doc.ComponentLike)
}

func (v *event_) GetKind() fra.ResourceLike {
	var component = v.GetObject(fra.Symbol("kind")).GetComponent()
	var kind = component.GetEntity().(fra.ResourceLike)
	return kind
}

// Attribute Methods

// Parameterized Methods

func (v *event_) GetEntity() any {
	return v.Declarative.(doc.ComponentLike).GetEntity()
}

func (v *event_) GetType() fra.ResourceLike {
	var component = v.GetParameter(fra.Symbol("type"))
	return fra.ResourceFromString(doc.FormatComponent(component))
}

func (v *event_) GetTag() fra.TagLike {
	var component = v.GetParameter(fra.Symbol("tag"))
	return fra.TagFromString(doc.FormatComponent(component))
}

func (v *event_) GetVersion() fra.VersionLike {
	var component = v.GetParameter(fra.Symbol("version"))
	return fra.VersionFromString(doc.FormatComponent(component))
}

func (v *event_) GetPermissions() fra.ResourceLike {
	var component = v.GetParameter(fra.Symbol("permissions"))
	return fra.ResourceFromString(doc.FormatComponent(component))
}

func (v *event_) GetOptionalPrevious() fra.ResourceLike {
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

type event_ struct {
	// Declare the instance attributes.

	// Declare the inherited aspects.
	doc.Declarative
}

// Class Structure

type eventClass_ struct {
	// Declare the class constants.
}

// Class Reference

func eventClass() *eventClass_ {
	return eventClassReference_
}

var eventClassReference_ = &eventClass_{
	// Initialize the class constants.
}
