package ent_test

import (
	"github.com/aws-contrib/aurora/internal/database/ent"

	. "github.com/aws-contrib/aurora/internal/database/ent/fake"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gateway", Ordered, func() {
	var gateway ent.Gateway

	BeforeEach(func() {
		var err error
		gateway, err = NewGateway()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		gateway.Close()
	})

	Describe("Lock", func() {
		var entity *ent.Lock

		BeforeAll(func() {
			entity = NewFakeLock()
		})

		Describe("CreateTableLocks", func() {
			It("creates the aurora_schema_revisions table", func(ctx SpecContext) {
				Expect(gateway.CreateTableLocks(ctx)).To(Succeed())
			})
		})

		Describe("InsertLock", func() {
			var params *ent.InsertLockParams

			BeforeEach(func() {
				params = &ent.InsertLockParams{}
				params.SetLock(entity)
			})

			It("inserts a revision", func(ctx SpecContext) {
				revision, err := gateway.InsertLock(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(revision).To(BeComparableTo(entity))
			})
		})

		Describe("GetLock", func() {
			var params *ent.GetLockParams

			BeforeEach(func() {
				params = &ent.GetLockParams{}
				params.SetLock(entity)
			})

			It("returns a revision", func(ctx SpecContext) {
				revision, err := gateway.GetLock(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(revision).To(BeComparableTo(entity))
			})
		})

		Describe("DeleteLock", func() {
			var params *ent.DeleteLockParams

			BeforeEach(func() {
				params = &ent.DeleteLockParams{}
				params.SetLock(entity)
			})

			It("deletes a revision", func(ctx SpecContext) {
				revision, err := gateway.DeleteLock(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(revision).To(BeComparableTo(entity))
			})
		})

		Describe("ExecInsertLock", func() {
			var params *ent.ExecInsertLockParams

			BeforeEach(func() {
				params = &ent.ExecInsertLockParams{}
				params.SetLock(entity)
			})

			It("inserts a revision", func(ctx SpecContext) {
				Expect(gateway.ExecInsertLock(ctx, params)).To(Succeed())
			})
		})

		Describe("ExecDeleteLock", func() {
			var params *ent.ExecDeleteLockParams

			BeforeEach(func() {
				params = &ent.ExecDeleteLockParams{}
				params.SetLock(entity)
			})

			It("inserts a revision", func(ctx SpecContext) {
				Expect(gateway.ExecDeleteLock(ctx, params)).To(Succeed())
			})
		})
	})
})
