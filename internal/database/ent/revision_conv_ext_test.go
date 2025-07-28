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

var _ = Describe("UpsertRevisionParams", func() {
	var params ent.UpsertRevisionParams

	BeforeEach(func() {
		params = ent.UpsertRevisionParams{}
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

var _ = Describe("ExecUpsertRevisionParams", func() {
	var params ent.ExecUpsertRevisionParams

	BeforeEach(func() {
		params = ent.ExecUpsertRevisionParams{}
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

var _ = Describe("UpdateRevisionParams", func() {
	var params ent.UpdateRevisionParams

	BeforeEach(func() {
		params = ent.UpdateRevisionParams{}
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

var _ = Describe("ExecUpdateRevisionParams", func() {
	var params ent.ExecUpdateRevisionParams

	BeforeEach(func() {
		params = ent.ExecUpdateRevisionParams{}
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
