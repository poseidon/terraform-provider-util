package util

import (
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const registerInitial = `
resource "util_register" "example" {
	set = "a1b2c3"
}
`

const registerUnsetSHA = `
resource "util_register" "example" {
	set = ""
}
`

const registerUpdateSHA = `
resource "util_register" "example" {
	set = "b2c3d4"
}
`

func TestRegister(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			// set initial value
			r.TestStep{
				Config: registerInitial,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("util_register.example", "value", "a1b2c3"),
				),
			},
			// set with empty value doesn't change value
			r.TestStep{
				Config: registerUnsetSHA,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("util_register.example", "value", "a1b2c3"),
				),
			},
			// set with content updates values
			r.TestStep{
				Config: registerUpdateSHA,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("util_register.example", "value", "b2c3d4"),
				),
			},
			// suppress noisy diffs that won't affect value
			r.TestStep{
				Config:   registerUnsetSHA,
				PlanOnly: true,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("util_register.example", "value", "b2c3d4"),
				),
			},
		},
	})
}
