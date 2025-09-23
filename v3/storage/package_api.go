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
Package "storage" provides implementations of data storage mechanisms that can
be used by a document repository.

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
package storage

import (
	doc "github.com/bali-nebula/go-bali-documents/v3"
	not "github.com/bali-nebula/go-digital-notary/v3"
	rep "github.com/bali-nebula/go-document-repository/v3/repository"
)

// TYPE DECLARATIONS

// FUNCTIONAL DECLARATIONS

// CLASS DECLARATIONS

/*
CachedStorageClassLike is a class interface that declares the complete set of
class constructors, constants and functions that must be supported by each
concrete cached-storage-like class.
*/
type CachedStorageClassLike interface {
	// Constructor Methods
	CachedStorage(
		storage rep.Persistent,
	) CachedStorageLike
}

/*
LocalStorageClassLike is a class interface that declares the complete set of
class constructors, constants and functions that must be supported by each
concrete local-storage-like class.
*/
type LocalStorageClassLike interface {
	// Constructor Methods
	LocalStorage(
		notary not.DigitalNotaryLike,
		directory string,
	) LocalStorageLike
}

/*
RemoteStorageClassLike is a class interface that declares the complete set of
class constructors, constants and functions that must be supported by each
concrete remote-storage-like class.
*/
type RemoteStorageClassLike interface {
	// Constructor Methods
	RemoteStorage(
		notary not.DigitalNotaryLike,
		service doc.ResourceLike,
	) RemoteStorageLike
}

/*
S3StorageClassLike is a class interface that declares the complete set of class
constructors, constants and functions that must be supported by each concrete
s3-storage-like class.
*/
type S3StorageClassLike interface {
	// Constructor Methods
	S3Storage(
		notary not.DigitalNotaryLike,
	) S3StorageLike
}

/*
SecureStorageClassLike is a class interface that declares the complete set of
class constructors, constants and functions that must be supported by each
concrete secure-storage-like class.
*/
type SecureStorageClassLike interface {
	// Constructor Methods
	SecureStorage(
		notary not.DigitalNotaryLike,
		storage rep.Persistent,
	) SecureStorageLike
}

// INSTANCE DECLARATIONS

/*
CachedStorageLike is an instance interface that declares the complete set of
principal, attribute and aspect methods that must be supported by each instance
of a concrete cached-storage-like class.
*/
type CachedStorageLike interface {
	// Principal Methods
	GetClass() CachedStorageClassLike

	// Aspect Interfaces
	rep.Persistent
}

/*
LocalStorageLike is an instance interface that declares the complete set of
principal, attribute and aspect methods that must be supported by each instance
of a concrete local-storage-like class.
*/
type LocalStorageLike interface {
	// Principal Methods
	GetClass() LocalStorageClassLike

	// Aspect Interfaces
	rep.Persistent
}

/*
RemoteStorageLike is an instance interface that declares the complete set of
principal, attribute and aspect methods that must be supported by each instance
of a concrete remote-storage-like class.
*/
type RemoteStorageLike interface {
	// Principal Methods
	GetClass() RemoteStorageClassLike

	// Aspect Interfaces
	rep.Persistent
}

/*
S3StorageLike is an instance interface that declares the complete set of
principal, attribute and aspect methods that must be supported by each instance
of a concrete s3-storage-like class.
*/
type S3StorageLike interface {
	// Principal Methods
	GetClass() S3StorageClassLike

	// Aspect Interfaces
	rep.Persistent
}

/*
SecureStorageLike is an instance interface that declares the complete set of
principal, attribute and aspect methods that must be supported by each instance
of a concrete secure-storage-like class.
*/
type SecureStorageLike interface {
	// Principal Methods
	GetClass() SecureStorageClassLike

	// Aspect Interfaces
	rep.Persistent
}

// ASPECT DECLARATIONS
