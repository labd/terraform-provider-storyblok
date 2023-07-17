# Storyblok Terraform Provider

[![Test status](https://github.com/labd/terraform-provider-storyblok/workflows/Run%20Tests/badge.svg)](https://github.com/labd/terraform-provider-storyblok/actions?query=workflow%3A%22Run+Tests%22)
[![codecov](https://codecov.io/gh/LabD/terraform-provider-storyblok/branch/master/graph/badge.svg)](https://codecov.io/gh/LabD/terraform-provider-storyblok)
[![Go Report Card](https://goreportcard.com/badge/github.com/labd/terraform-provider-storyblok)](https://goreportcard.com/report/github.com/labd/terraform-provider-storyblok)

The Terraform Storyblok provider allows you to configure your
[storyblok](https://storyblok.com/) space with infrastructure-as-code
principles.

# Commercial support

Need support implementing this terraform module in your organization? We are
able to offer support. Please contact us at opensource@labdigital.nl

# Quick start

[Read our documentation](https://registry.terraform.io/providers/labd/storyblok/latest/docs)
and check out the [examples](https://registry.terraform.io/providers/labd/storyblok/latest/docs/guides/examples).

## Usage

The provider is distributed via the Terraform registry. To use it you need to configure the [`required_provider`](https://www.terraform.io/language/providers/requirements#requiring-providers) block. For example:

```hcl
terraform {
  required_providers {
    storyblok = {
      source = "labd/storyblok"

      # It's recommended to pin the version, e.g.:
      # version = "~> 0.0.1"
    }
  }
}
```

# Binaries

Packages of the releases are available at
https://github.com/labd/terraform-provider-storyblok/releases See the
[terraform documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)
for more information about installing third-party providers.

# Contributing

## Building the provider

Clone the repository and run the following command:

```sh
$ task build-local
```

## Debugging / Troubleshooting

There are two environment settings for troubleshooting:

- `TF_LOG=INFO` enables debug output for Terraform.

Note this generates a lot of output!

## Releasing

Install "changie"

```
brew tap miniscruff/changie https://github.com/miniscruff/changie
brew install changie
```

Add unreleased change files by running for each change (add/fix/remove/etc.)

```
changie new
```

Commit this and a new PR will be created.

Once that's merged and its Github action is complete, a new release will be live.

## Testing

### Running the unit tests

```sh
$ task test
```

## Authors

This project is developed by [Lab Digital](https://www.labdigital.nl). We
welcome additional contributors. Please see our
[GitHub repository](https://github.com/labd/terraform-provider-storyblok)
for more information.
