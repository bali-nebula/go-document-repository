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
┌────────────────────────────────── WARNING ───────────────────────────────────┐
│         This "module_api.go" file was automatically generated using:         │
│            https://github.com/craterdog/go-development-tools/wiki            │
│                                                                              │
│      Updates to any part of this file—other than the Module Description      │
│             and the Global Functions sections may be overwritten.            │
└──────────────────────────────────────────────────────────────────────────────┘

Package "module" declares type aliases for the commonly used types declared in
the packages contained in this module.  It also provides constructors for each
commonly used class that is exported by the module.  Each constructor delegates
the actual construction process to its corresponding concrete class declared in
the corresponding package contained within this module.

For detailed documentation on this entire module refer to the wiki:
  - https://github.com/bali-nebula/go-document-repository/wiki
*/
package module

import (
	not "github.com/bali-nebula/go-digital-notary/v3"
	rep "github.com/bali-nebula/go-document-repository/v3/repository"
	sto "github.com/bali-nebula/go-document-repository/v3/storage"
	fra "github.com/craterdog/go-component-framework/v7"
)

// TYPE ALIASES

// Repository

type (
	DocumentRepositoryClassLike = rep.DocumentRepositoryClassLike
)

type (
	DocumentRepositoryLike = rep.DocumentRepositoryLike
)

type (
	Persistent = rep.Persistent
)

// Storage

type (
	CachedStorageClassLike    = sto.CachedStorageClassLike
	LocalStorageClassLike     = sto.LocalStorageClassLike
	RemoteStorageClassLike    = sto.RemoteStorageClassLike
	S3StorageClassLike        = sto.S3StorageClassLike
	ValidatedStorageClassLike = sto.ValidatedStorageClassLike
)

type (
	CachedStorageLike    = sto.CachedStorageLike
	LocalStorageLike     = sto.LocalStorageLike
	RemoteStorageLike    = sto.RemoteStorageLike
	S3StorageLike        = sto.S3StorageLike
	ValidatedStorageLike = sto.ValidatedStorageLike
)

// CLASS ACCESSORS

// Repository

func DocumentRepositoryClass() DocumentRepositoryClassLike {
	return rep.DocumentRepositoryClass()
}

func DocumentRepository(
	notary not.DigitalNotaryLike,
	storage rep.Persistent,
) DocumentRepositoryLike {
	return DocumentRepositoryClass().DocumentRepository(
		notary,
		storage,
	)
}

// Storage

func CachedStorageClass() CachedStorageClassLike {
	return sto.CachedStorageClass()
}

func CachedStorage(
	storage rep.Persistent,
) CachedStorageLike {
	return CachedStorageClass().CachedStorage(
		storage,
	)
}

func LocalStorageClass() LocalStorageClassLike {
	return sto.LocalStorageClass()
}

func LocalStorage(
	notary not.DigitalNotaryLike,
	directory string,
) LocalStorageLike {
	return LocalStorageClass().LocalStorage(
		notary,
		directory,
	)
}

func RemoteStorageClass() RemoteStorageClassLike {
	return sto.RemoteStorageClass()
}

func RemoteStorage(
	notary not.DigitalNotaryLike,
	service fra.ResourceLike,
) RemoteStorageLike {
	return RemoteStorageClass().RemoteStorage(
		notary,
		service,
	)
}

func S3StorageClass() S3StorageClassLike {
	return sto.S3StorageClass()
}

func S3Storage(
	notary not.DigitalNotaryLike,
) S3StorageLike {
	return S3StorageClass().S3Storage(
		notary,
	)
}

func ValidatedStorageClass() ValidatedStorageClassLike {
	return sto.ValidatedStorageClass()
}

func ValidatedStorage(
	notary not.DigitalNotaryLike,
	storage rep.Persistent,
) ValidatedStorageLike {
	return ValidatedStorageClass().ValidatedStorage(
		notary,
		storage,
	)
}

// GLOBAL FUNCTIONS
