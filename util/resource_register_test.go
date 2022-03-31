package util

import (
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const registerInitial = `
resource "util_register" "example" {
	content = "a1b2c3"
}

output "out" {
	value = util_register.example.value
}
`

const registerUnsetSHA = `
resource "util_register" "example" {
	content = ""
}

output "out" {
	value = util_register.example.value
}
`

const registerUpdateSHA = `
resource "util_register" "example" {
	content = "b2c3d4"
}

output "out" {
	value = util_register.example.value
}
`

// register expected values
const (
	registerInitialExpected = "a1b2c3"
	registerUpdateExpected  = "b2c3d4"
)

func TestRegister(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			// set initial value
			{
				Config: registerInitial,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("util_register.example", "value", registerInitialExpected),
					r.TestCheckOutput("out", registerInitialExpected),
				),
			},
			// Empty content does NOT change value
			{
				Config: registerUnsetSHA,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("util_register.example", "value", registerInitialExpected),
					r.TestCheckOutput("out", registerInitialExpected),
				),
			},
			// Non-empty content DOES change value
			{
				Config: registerUpdateSHA,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("util_register.example", "value", registerUpdateExpected),
					r.TestCheckOutput("out", registerUpdateExpected),
				),
			},
			// suppress noisy diffs that won't affect value
			{
				Config:   registerUnsetSHA,
				PlanOnly: true,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("util_register.example", "value", registerUpdateExpected),
					r.TestCheckOutput("out", registerUpdateExpected),
				),
			},
		},
	})
}
