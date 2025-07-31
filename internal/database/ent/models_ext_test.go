package ent_test

import (
	"github.com/aws-contrib/aurora/internal/database/ent"

	. "github.com/aws-contrib/aurora/internal/database/ent/fake"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Revision", func() {
	var entity *ent.Revision

	BeforeEach(func() {
		entity = NewFakeRevision()
	})

	Describe("GetName", func() {
		It("returns the name", func() {
			Expect(entity.GetName()).NotTo(BeEmpty())
		})
	})

	Describe("SetName", func() {
		It("sets the name", func() {
			entity.SetName("id_description.sql")
			Expect(entity.ID).To(Equal("id"))
			Expect(entity.Description).To(Equal("description"))
		})
	})
})
