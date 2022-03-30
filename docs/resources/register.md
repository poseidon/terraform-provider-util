# Register

Store a content value that persists in state until changed to a non-empty value.

```tf
resource "util_register" "example" {
  content = "a1b2c3"
}
```

Later, the register's content may be updated, but empty values (`null` or `""`) are ignored.

```tf
resource "util_register" "example" {
  content = null
}

output "sha" {
  value = util_register.example.value  # "a1b2c3"
}
```

## Argument Reference

* `content` - set the register value (`""` or `null` values ignored)

## Attribute Reference

* `value` - computed register value
