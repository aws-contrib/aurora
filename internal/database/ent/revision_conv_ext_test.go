//go:build !goverter

package ent_test

import (
	"github.com/aws-contrib/aurora/internal/database/ent"

	. "github.com/aws-contrib/aurora/internal/database/ent/fake"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetRevisionParams", func() {
	var params ent.GetRevisionParams

	BeforeEach(func() {
		params = ent.GetRevisionParams{}
	})

	Describe("SetRevision", func() {
		var entity *ent.Revision

		BeforeEach(func() {
			entity = NewFakeRevision()
		})

		It("sets the entity", func() {
			params.SetRevision(entity)
			Expect(params).NotTo(BeZero())
		})
	})
})

var _ = Describe("InsertRevisionParams", func() {
	var params ent.InsertRevisionParams

	BeforeEach(func() {
		params = ent.InsertRevisionParams{}
	})

	Describe("SetRevision", func() {
		var entity *ent.Revision

		BeforeEach(func() {
			entity = NewFakeRevision()
		})

		It("sets the entity", func() {
			params.SetRevision(entity)
			Expect(params).NotTo(BeZero())
		})
	})
})

var _ = Describe("ExecInsertRevisionParams", func() {
	var params ent.ExecInsertRevisionParams

	BeforeEach(func() {
		params = ent.ExecInsertRevisionParams{}
	})

	Describe("SetRevision", func() {
		var entity *ent.Revision

		BeforeEach(func() {
			entity = NewFakeRevision()
		})

		It("sets the entity", func() {
			params.SetRevision(entity)
			Expect(params).NotTo(BeZero())
		})
	})
})
