//go:build !goverter

package ent_test

import (
	"github.com/aws-contrib/aurora/internal/database/ent"

	. "github.com/aws-contrib/aurora/internal/database/ent/fake"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetJobParams", func() {
	var params ent.GetJobParams

	BeforeEach(func() {
		params = ent.GetJobParams{}
	})

	Describe("SetJob", func() {
		var entity *ent.Job

		BeforeEach(func() {
			entity = NewFakeJob()
		})

		It("sets the entity", func() {
			params.SetJob(entity)
			Expect(params).NotTo(BeZero())
		})
	})
})

var _ = Describe("InsertJobParams", func() {
	var params ent.InsertJobParams

	BeforeEach(func() {
		params = ent.InsertJobParams{}
	})

	Describe("SetJob", func() {
		var entity *ent.Job

		BeforeEach(func() {
			entity = NewFakeJob()
		})

		It("sets the entity", func() {
			params.SetJob(entity)
			Expect(params).NotTo(BeZero())
		})
	})
})

var _ = Describe("ExecInsertJobParams", func() {
	var params ent.ExecInsertJobParams

	BeforeEach(func() {
		params = ent.ExecInsertJobParams{}
	})

	Describe("SetJob", func() {
		var entity *ent.Job

		BeforeEach(func() {
			entity = NewFakeJob()
		})

		It("sets the entity", func() {
			params.SetJob(entity)
			Expect(params).NotTo(BeZero())
		})
	})
})

var _ = Describe("DeleteRevisionParams", func() {
	var params ent.DeleteRevisionParams

	BeforeEach(func() {
		params = ent.DeleteRevisionParams{}
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

var _ = Describe("ExecDeleteRevisionParams", func() {
	var params ent.ExecDeleteRevisionParams

	BeforeEach(func() {
		params = ent.ExecDeleteRevisionParams{}
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
