/*
................................................................................
.    Copyright (c) 2009-2026 Crater Dog Technologies™.  All Rights Reserved.   .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See https://opensource.org/license/MIT)                        .
................................................................................
*/

package module_test

import (
	sig "crypto/ed25519"
	fmt "fmt"
	doc "github.com/bali-nebula/go-bali-documents/v3"
	not "github.com/bali-nebula/go-digital-notary/v3"
	uti "github.com/craterdog/go-essential-utilities/v8"
)

// CLASS INTERFACE

// Access Function

func HsmEd25519TestClass() not.HsmEd25519ClassLike {
	return hsmEd25519Class()
}

// Constructor Methods

func (c *hsmEd25519Class_) HsmEd25519(
	device string,
	tag string,
) not.HsmEd25519Like {
	if uti.IsUndefined(tag) {
		panic("The \"tag\" attribute is required by this class.")
	}
	var filename = testDirectory + "hsmEd25519/Configuration.bali"
	var controller = uti.Controller(c.events_, c.transitions_, c.keyless_)
	var instance = &hsmEd25519_{
		// Initialize the instance attributes.
		filename_:   filename,
		controller_: controller,
	}
	if uti.PathExists(filename) {
		instance.readConfiguration(tag)
	} else {
		instance.createConfiguration(tag)
	}
	return instance
}

// Constant Methods

// Function Methods

// INSTANCE INTERFACE

// Principal Methods

func (v *hsmEd25519_) GetClass() not.HsmEd25519ClassLike {
	return hsmEd25519Class()
}

// Attribute Methods

// Hardened Methods

func (v *hsmEd25519_) GetTag() string {
	return v.tag_
}

func (v *hsmEd25519_) GetSignatureAlgorithm() string {
	return hsmEd25519Class().algorithm_
}

func (v *hsmEd25519_) GetPublicKey() []byte {
	return v.publicKey_
}

func (v *hsmEd25519_) GenerateKeys() []byte {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to generate new keys",
	)

	var err error
	v.controller_.ProcessEvent(hsmEd25519Class().generateKeys_)
	v.publicKey_, v.privateKey_, err = sig.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
	v.writeConfiguration()
	return v.publicKey_
}

func (v *hsmEd25519_) SignBytes(
	bytes []byte,
) []byte {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to sign bytes",
	)

	v.controller_.ProcessEvent(hsmEd25519Class().signBytes_)
	var privateKey = v.privateKey_
	if v.previousKey_ != nil {
		// Use the old key one more time to sign the new one.
		privateKey = v.previousKey_
		v.previousKey_ = nil
	}
	var signature = sig.Sign(privateKey, bytes)
	v.writeConfiguration()
	return signature
}

func (v *hsmEd25519_) IsValid(
	key []byte,
	bytes []byte,
	signature []byte,
) bool {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to verify bytes signature",
	)

	return sig.Verify(sig.PublicKey(key), bytes, signature)
}

func (v *hsmEd25519_) RotateKeys() []byte {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to rotate keys",
	)

	var err error
	v.controller_.ProcessEvent(hsmEd25519Class().rotateKeys_)
	v.previousKey_ = v.privateKey_
	v.publicKey_, v.privateKey_, err = sig.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
	v.writeConfiguration()
	return v.publicKey_
}

func (v *hsmEd25519_) EraseKeys() {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to erase the keys",
	)

	v.createConfiguration(v.tag_)
}

// PROTECTED INTERFACE

// Private Methods

func (v *hsmEd25519_) createConfiguration(
	tag string,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to create a new HSM configuration",
	)

	v.tag_ = tag
	v.publicKey_ = nil
	v.privateKey_ = nil
	v.previousKey_ = nil
	v.controller_ = uti.Controller(
		hsmEd25519Class().events_,
		hsmEd25519Class().transitions_,
		hsmEd25519Class().keyless_,
	)
	v.writeConfiguration()
}

func (v *hsmEd25519_) errorCheck(
	message string,
) {
	if e := recover(); e != nil {
		message = fmt.Sprintf(
			"HsmEd25519: %s:\n        %v",
			message,
			e,
		)
		panic(message)
	}
}

func (v *hsmEd25519_) readConfiguration(
	tag string,
) {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to read in the HSM configuration",
	)

	var source = uti.ReadFile(v.filename_)
	var component = doc.ParseComponent(source)

	v.tag_ = doc.FormatComponent(
		component.GetSubcomponent(doc.Symbol("$tag")),
	)
	if v.tag_ != tag {
		panic("The specified tag does not match the HSM tag.")
	}

	var publicKey = doc.FormatComponent(
		component.GetSubcomponent(doc.Symbol("$publicKey")),
	)
	if publicKey != "none" {
		v.publicKey_ = doc.Binary(publicKey).AsIntrinsic()
	}

	var privateKey = doc.FormatComponent(
		component.GetSubcomponent(doc.Symbol("$privateKey")),
	)
	if privateKey != "none" {
		v.privateKey_ = doc.Binary(privateKey).AsIntrinsic()
	}

	var previousKey = doc.FormatComponent(
		component.GetSubcomponent(doc.Symbol("$previousKey")),
	)
	if previousKey != "none" {
		v.previousKey_ = doc.Binary(previousKey).AsIntrinsic()
	}

	var state = doc.FormatComponent(
		component.GetSubcomponent(doc.Symbol("$state")),
	)
	switch state {
	case "$Keyless":
		v.controller_.SetState(hsmEd25519Class().keyless_)
	case "$LoneKey":
		v.controller_.SetState(hsmEd25519Class().loneKey_)
	case "$TwoKeys":
		v.controller_.SetState(hsmEd25519Class().twoKeys_)
	default:
		panic("Invalid State")
	}
}

func (v *hsmEd25519_) writeConfiguration() {
	// Check for any errors at the end.
	defer v.errorCheck(
		"An error occurred while attempting to write out the HSM configuration",
	)

	var tag = v.tag_

	var state string
	switch v.controller_.GetState() {
	case hsmEd25519Class().keyless_:
		state = "$Keyless"
	case hsmEd25519Class().loneKey_:
		state = "$LoneKey"
	case hsmEd25519Class().twoKeys_:
		state = "$TwoKeys"
	default:
		panic("Invalid State")
	}

	var publicKey = "none"
	if uti.IsDefined(v.publicKey_) {
		publicKey = doc.Binary(v.publicKey_).AsSource()
	}

	var privateKey = "none"
	if uti.IsDefined(v.privateKey_) {
		privateKey = doc.Binary(v.privateKey_).AsSource()
	}

	var previousKey = "none"
	if uti.IsDefined(v.previousKey_) {
		previousKey = doc.Binary(v.previousKey_).AsSource()
	}

	var source = `[
    $tag: ` + tag + `
    $state: ` + state + `
    $publicKey: ` + publicKey + `
    $privateKey: ` + privateKey + `
    $previousKey: ` + previousKey + `
](
    $type: /bali/types/notary/HsmEd25519/v3
)
`
	uti.WriteFile(v.filename_, source)
}

// Instance Structure

type hsmEd25519_ struct {
	// Declare the instance attributes.
	tag_         string
	publicKey_   []byte
	privateKey_  []byte
	previousKey_ []byte
	filename_    string
	controller_  uti.Stateful
}

// Class Structure

type hsmEd25519Class_ struct {
	// Declare the class constants.
	algorithm_    string
	keyless_      uti.State
	loneKey_      uti.State
	twoKeys_      uti.State
	generateKeys_ uti.Event
	signBytes_    uti.Event
	rotateKeys_   uti.Event
	events_       []uti.Event
	transitions_  map[uti.State]uti.Transitions
}

// Class Reference

func hsmEd25519Class() *hsmEd25519Class_ {
	return hsmEd25519ClassReference_
}

var hsmEd25519ClassReference_ = &hsmEd25519Class_{
	// Initialize the class constants.
	algorithm_:    "ED25519",
	keyless_:      "$Keyless",
	loneKey_:      "$LoneKey",
	twoKeys_:      "$TwoKeys",
	generateKeys_: "$GenerateKeys",
	signBytes_:    "$SignBytes",
	rotateKeys_:   "$RotateKeys",
	events_:       []uti.Event{"$GenerateKeys", "$SignBytes", "$RotateKeys"},
	transitions_: map[uti.State]uti.Transitions{
		"$Keyless": uti.Transitions{"$LoneKey", "$Invalid", "$Invalid"},
		"$LoneKey": uti.Transitions{"$Invalid", "$LoneKey", "$TwoKeys"},
		"$TwoKeys": uti.Transitions{"$Invalid", "$LoneKey", "$Invalid"},
	},
}
