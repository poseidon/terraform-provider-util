# replace Data Source

`replace` searches a given `content` string, performs a set of substring matches and `replacements`. Similar to the builtin [replace](https://www.terraform.io/docs/language/functions/replace.html).

## Usage

```hcl
data "util_replace" "example" {
  content      = "hello world"
  replacements = {
    "/(h|H)ello/": "Hallo",
    "world": "Welt",
  }
}

# Hallo Welt
output "example" {
  value = data.ct_replace.example.replaced
}
```

Note that maps are iterated in lexographical key order, so care may be needed if overlappping substring matchers are desired.

## Argument Reference

* `content` - content to search and replace
* `replacements` - map of substring matchers to replacements to perform. If a substring is wrapped in forward slashes, it is treated as a regular expression, using the same pattern syntax as regex

## Argument Attributes

* `replaced` - content with replacements performed

