# terraform-provider-util

Notable changes between releases.

## Latest

## v0.2.2

* Improve `util_register` plan diff to show expected value

## v0.2.1

* Rename `util_register` field from `set` to `content` ([#9](https://github.com/poseidon/terraform-provider-util/pull/9))
* Fix `util_register` to mark `value` attribute as computed to propagate changes ([#8](https://github.com/poseidon/terraform-provider-util/pull/8))

## v0.2.0

* Add `util_register` resource for storing values ([#7](https://github.com/poseidon/terraform-provider-util/pull/7))

## v0.1.0

* Add `util_replace` data source function
