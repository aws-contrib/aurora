//go:build !goverter

package ent_test

import (
	"github.com/aws-contrib/aurora/internal/database/ent"

	. "github.com/aws-contrib/aurora/internal/database/ent/fake"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetLockParams", func() {
	var params ent.GetLockParams

	BeforeEach(func() {
		params = ent.GetLockParams{}
	})

	Describe("SetLock", func() {
		var entity *ent.Lock

		BeforeEach(func() {
			entity = NewFakeLock()
		})

		It("sets the entity", func() {
			params.SetLock(entity)
			Expect(params).NotTo(BeZero())
		})
	})
})

var _ = Describe("InsertLockParams", func() {
	var params ent.InsertLockParams

	BeforeEach(func() {
		params = ent.InsertLockParams{}
	})

	Describe("SetLock", func() {
		var entity *ent.Lock

		BeforeEach(func() {
			entity = NewFakeLock()
		})

		It("sets the entity", func() {
			params.SetLock(entity)
			Expect(params).NotTo(BeZero())
		})
	})
})

var _ = Describe("ExecInsertLockParams", func() {
	var params ent.ExecInsertLockParams

	BeforeEach(func() {
		params = ent.ExecInsertLockParams{}
	})

	Describe("SetLock", func() {
		var entity *ent.Lock

		BeforeEach(func() {
			entity = NewFakeLock()
		})

		It("sets the entity", func() {
			params.SetLock(entity)
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
