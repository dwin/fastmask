# fastmask

Quickly create [Fastmail](https://fastmail.com) Masked Email addresses. It includes a client package for use in other projects in addition to the CLI.

>**What is Fastmail Masked Email ?**
>
> A Masked Email address is a unique, automatically generated email address that can be used in place of your real email address.
>
> Masked Email addresses are especially useful when you need to sign up with new services online. Instead of sharing your real email address, keep it private and protect yourself from data breaches and spam by creating a new Masked Email for every service.
>
> If a Masked Email address starts receiving unwanted mail, you can simply disable that address. Masked Email addresses also make it easy to identify which service shared or leaked the address.
>
> More info: [https://www.fastmail.help/hc/en-us/articles/4406536368911-Masked-Email](https://www.fastmail.help/hc/en-us/articles/4406536368911-Masked-Email)

[![Go Reference](https://pkg.go.dev/badge/github.com/dwin/fastmask.svg)](https://pkg.go.dev/github.com/dwin/fastmask)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/dwin/fastmask)
[![Go Report Card](https://goreportcard.com/badge/github.com/dwin/fastmask)](https://goreportcard.com/report/github.com/dwin/fastmask)
[![codecov](https://codecov.io/gh/dwin/fastmask/branch/main/graph/badge.svg?token=Co4xYYdgVV)](https://codecov.io/gh/dwin/fastmask)
![GitHub branch checks state](https://img.shields.io/github/checks-status/dwin/fastmask/main)

## Installation

Get the latest releases for Linux, macOS(darwin) or Windows for armv6, armv7, arm64 and x86_64 architectures from:

[https://github.com/dwin/fastmask/releases](https://github.com/dwin/fastmask/releases) or:

```bash
go get -u github.com/dwin/fastmask
```

Or build from source.

## Usage

### CLI

```bash
fastmask login -u <email> -p <password> -m <mfa_code>
fastmask create <website> -d <description>
```

Fastmask will store the credentials in `~/.fastmask/.config.yaml`.

_Description is optional._
_MFA code is required only if enabled for your account **(it should be)**._

### Go Package

```go
import "github.com/dwin/fastmask/pkg/fastmail"

client := fastmail.NewClient("your-app-name")
```

## License

See [LICENSE](/LICENSE) for details.

## Contributing

See [CONTRIBUTING.md](/CONTRIBUTING.md) for more information.

## Future Improvements

- [ ] Improve test coverage.
- [ ] Add support for OAuth, pending Fastmail response.
- [ ] Prompt for MFA code if needed.
- [ ] Prompt for credentials if needed.
- [ ] Add support for verbose logging output.
- [ ] Add support for listing and filtering Masked Email addresses. (currently mask be managed in Fastmail settings.)
- [ ] Add support for passing credentials via environment variables or flags for scripting.
