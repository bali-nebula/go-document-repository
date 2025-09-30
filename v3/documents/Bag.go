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

func BagClass() BagClassLike {
	return bagClass()
}

// Constructor Methods

func (c *bagClass_) Bag(
	name doc.NameLike,
	capacity doc.NumberLike,
	leasetime doc.NumberLike,
	permissions doc.ResourceLike,
) BagLike {
	if uti.IsUndefined(name) {
		panic("The \"name\" attribute is required by this class.")
	}
	if uti.IsUndefined(capacity) {
		panic("The \"capacity\" attribute is required by this class.")
	}
	if uti.IsUndefined(leasetime) {
		panic("The \"leasetime\" attribute is required by this class.")
	}
	if uti.IsUndefined(permissions) {
		panic("The \"permissions\" attribute is required by this class.")
	}
	var tag = doc.Tag()
	var version = doc.Version()
	var source = `[
    $name: ` + name.AsString() + `
    $capacity: ` + capacity.AsString() + `
    $leasetime: ` + leasetime.AsString() + `
](
    $type: <bali:/types/documents/Bag:v3>
    $tag: ` + tag.AsString() + `
    $version: ` + version.AsString() + `
    $permissions: ` + permissions.AsString() + `
)`
	var component = doc.ParseSource(source)
	var instance = &bag_{
		// Initialize the instance attributes.

		// Initialize the inherited aspects.
		Declarative: component,
	}
	return instance
}

func (c *bagClass_) BagFromString(
	source string,
) BagLike {
	var component = doc.ParseSource(source)
	var instance = &bag_{
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

func (v *bag_) GetClass() BagClassLike {
	return bagClass()
}

func (v *bag_) AsString() string {
	return doc.FormatDocument(v.Declarative.(doc.ComponentLike))
}

func (v *bag_) AsIntrinsic() doc.ComponentLike {
	return v.Declarative.(doc.ComponentLike)
}

func (v *bag_) GetName() doc.NameLike {
	var component = v.GetObject(doc.Symbol("$name")).GetComponent()
	var name = component.GetEntity().(doc.NameLike)
	return name
}

func (v *bag_) GetCapacity() doc.NumberLike {
	var component = v.GetObject(doc.Symbol("$capacity")).GetComponent()
	var capacity = component.GetEntity().(doc.NumberLike)
	return capacity
}

func (v *bag_) GetLeasetime() doc.NumberLike {
	var component = v.GetObject(doc.Symbol("$leasetime")).GetComponent()
	var leasetime = component.GetEntity().(doc.NumberLike)
	return leasetime
}

// Attribute Methods

// Parameterized Methods

func (v *bag_) GetEntity() any {
	return v.Declarative.(doc.ComponentLike).GetEntity()
}

func (v *bag_) GetType() doc.ResourceLike {
	var component = v.GetParameter(doc.Symbol("$type"))
	return doc.Resource(doc.FormatComponent(component))
}

func (v *bag_) GetTag() doc.TagLike {
	var component = v.GetParameter(doc.Symbol("$tag"))
	return doc.Tag(doc.FormatComponent(component))
}

func (v *bag_) GetVersion() doc.VersionLike {
	var component = v.GetParameter(doc.Symbol("$version"))
	return doc.Version(doc.FormatComponent(component))
}

func (v *bag_) GetOptionalPrevious() doc.ResourceLike {
	var previous doc.ResourceLike
	var component = v.GetParameter(doc.Symbol("$previous"))
	var source = doc.FormatComponent(component)
	if source != "none" {
		previous = doc.Resource(source)
	}
	return previous
}

func (v *bag_) GetPermissions() doc.ResourceLike {
	var component = v.GetParameter(doc.Symbol("$permissions"))
	return doc.Resource(doc.FormatComponent(component))
}

func (v *bag_) GetAccount() doc.TagLike {
	var component = v.GetParameter(doc.Symbol("$account"))
	return doc.Tag(doc.FormatComponent(component))
}

// PROTECTED INTERFACE

// Private Methods

// Instance Structure

type bag_ struct {
	// Declare the instance attributes.

	// Declare the inherited aspects.
	doc.Declarative
}

// Class Structure

type bagClass_ struct {
	// Declare the class constants.
}

// Class Reference

func bagClass() *bagClass_ {
	return bagClassReference_
}

var bagClassReference_ = &bagClass_{
	// Initialize the class constants.
}
