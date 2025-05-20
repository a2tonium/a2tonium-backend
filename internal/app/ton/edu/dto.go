package edu

import "github.com/xssnick/tonutils-go/address"

type CertificateInit struct {
	CollectionAddress *address.Address `json:"collection_address"`
	CertificateIndex  int64            `json:"certificate_index"`
}

type CourseInit struct {
	OwnerAddress *address.Address `json:"owner_address"`
	CourseIndex  int64            `json:"course_index"`
}
