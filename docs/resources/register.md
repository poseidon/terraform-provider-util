# Register

Store a value in state that persists until changed to a non-empty value.

```tf
resource "util_register" "example" {
  set = "a1b2c3"
}
```

Later, the register's value may be updated, but setting it to `null` or `""` is ignored.

```tf
resource "util_register" "example" {
  set = null
}

output "sha" {
  value = util_register.example.value  # "a1b2c3"
}
```

## Argument Reference

* `set` - set the register value (`""` or `null` values ignored)

## Attribute Reference

* `value` - computed value of the register
