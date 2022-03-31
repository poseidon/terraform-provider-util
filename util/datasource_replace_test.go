package util

import (
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const replaceExample1 = `
data "util_replace" "example" {
  content = "hello world"
	replacements = {
		"/(h|H)ello/": "Hallo",
		"world": "Welt",
	}
}
`

const replaceExample2 = `
data "util_replace" "example" {
  content = "Test content to change"
	replacements = {
		"/e/": "o",
		"/c/": "b",
		"st": "sts",
	}
}
`
const replaceExpected1 = "Hallo Welt"
const replaceExpected2 = "Tosts bontont to bhango"

func TestReplace(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: replaceExample1,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.util_replace.example", "replaced", replaceExpected1),
				),
			},
			{
				Config: replaceExample2,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.util_replace.example", "replaced", replaceExpected2),
				),
			},
		},
	})
}
