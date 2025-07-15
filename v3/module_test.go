/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package module_test

import (
	fmt "fmt"
	not "github.com/bali-nebula/go-digital-notary/v3"
	doc "github.com/bali-nebula/go-document-notation/v3"
	rep "github.com/bali-nebula/go-document-repository/v3"
	fra "github.com/craterdog/go-component-framework/v7"
	uti "github.com/craterdog/go-missing-utilities/v7"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

const testDirectory = "./test/"

func TestLocalStorage(t *tes.T) {
	uti.RemakeDirectory(testDirectory)
	var ssm = not.Ssm(testDirectory)
	var hsm = ssm
	var notary = not.DigitalNotary(ssm, hsm)
	notary.ForgetKey()
	notary.GenerateKey()
	var storage = rep.LocalStorage(notary, testDirectory)
	var repository = rep.DocumentRepository(notary, storage)
	var component = doc.ParseSource("~pi").GetComponent()
	var type_ = fra.ResourceFromString("<bali:/examples/Angle:v1>")
	var tag = fra.TagWithSize(20)
	var version = fra.VersionFromString("v1.2.3")
	var permissions = fra.ResourceFromString("<bali:/permissions/public:v1>")
	var previous not.CitationLike
	var draft = not.Draft(
		component,
		type_,
		tag,
		version,
		permissions,
		previous,
	)
	var citation = repository.SaveDraft(draft)
	var same = repository.RetrieveDraft(citation)
	ass.Equal(t, draft.AsString(), same.AsString())
	repository.DiscardDraft(citation)
	var resource = fra.ResourceFromString("<bali:/contracts/Test:v1.2.3>")
	var contract = repository.NotarizeDraft(resource, draft)
	var same2 = repository.RetrieveContract(resource)
	ass.Equal(t, contract.AsString(), same2.AsString())
	draft = repository.CheckoutDraft(resource, 2)
	fmt.Println(draft.AsString())
	ass.NotEqual(t, draft.AsString(), same.AsString())
}
