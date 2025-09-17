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

func BagClass() BagClassLike {
	return bagClass()
}

// Constructor Methods

func (c *bagClass_) Bag(
	name fra.NameLike,
	capacity fra.NumberLike,
	leasetime fra.NumberLike,
	permissions fra.ResourceLike,
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
	var tag = fra.TagWithSize(20)
	var source = `[
    $name: ` + name.AsString() + `
    $capacity: ` + capacity.AsString() + `
    $leasetime: ` + leasetime.AsString() + `
](
    $type: <bali:/types/documents/Bag:v3>
    $tag: ` + tag.AsString() + `
    $version: v1
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

func (v *bag_) GetName() fra.NameLike {
	var component = v.GetObject(fra.Symbol("name")).GetComponent()
	var name = component.GetEntity().(fra.NameLike)
	return name
}

func (v *bag_) GetCapacity() fra.NumberLike {
	var component = v.GetObject(fra.Symbol("capacity")).GetComponent()
	var capacity = component.GetEntity().(fra.NumberLike)
	return capacity
}

func (v *bag_) GetLeasetime() fra.NumberLike {
	var component = v.GetObject(fra.Symbol("leasetime")).GetComponent()
	var leasetime = component.GetEntity().(fra.NumberLike)
	return leasetime
}

// Attribute Methods

// Parameterized Methods

func (v *bag_) GetEntity() any {
	return v.Declarative.(doc.ComponentLike).GetEntity()
}

func (v *bag_) GetType() fra.ResourceLike {
	var component = v.GetParameter(fra.Symbol("type"))
	return fra.ResourceFromString(doc.FormatComponent(component))
}

func (v *bag_) GetTag() fra.TagLike {
	var component = v.GetParameter(fra.Symbol("tag"))
	return fra.TagFromString(doc.FormatComponent(component))
}

func (v *bag_) GetVersion() fra.VersionLike {
	var component = v.GetParameter(fra.Symbol("version"))
	return fra.VersionFromString(doc.FormatComponent(component))
}

func (v *bag_) GetPermissions() fra.ResourceLike {
	var component = v.GetParameter(fra.Symbol("permissions"))
	return fra.ResourceFromString(doc.FormatComponent(component))
}

func (v *bag_) GetOptionalPrevious() fra.ResourceLike {
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
