
# ci

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/getoutreach/ci)
[![CircleCI](https://circleci.com/gh/getoutreach/ci.svg?style=shield)](https://circleci.com/gh/getoutreach/ci)
[![Generated via Stencil](https://img.shields.io/badge/Outreach-Stencil-%235951ff)](https://github.com/getoutreach/stencil)

<!--- Block(description) -->
Collection of utilities for use in CI @ Outreach
<!--- EndBlock(description) -->

----

Please read the [CONTRIBUTING.md](CONTRIBUTING.md) document for guidelines on developing and contributing changes.

<!--- Block(custom) -->
## CI

### Using in CircleCI

**Note**: To get credentials you will need to ensure your CI jobs includes the context `ghaccesstoken`

To use in CircleCI you must first download this binary as part of a base image or machine image. The following Docker images have this installed:

 * `gcr.io/outreach-docker/bootstrap/ci`: `ghaccesstoken`

For machine mode, the following snippet may be used (Bootstrap):

```bash
./scripts/shell-wrapper.sh gobin.sh github.com/getoutreach/ci/cmd/ghaccesstoken@<version-here>
```

This will vend credentials. It can be used like: `export GH_TOKEN=$(<command-string-from-above> token)`. That will set a valid github token as the env var `GH_TOKEN`.

**Note**: This Token is valid at the time of issuance for an unknown (potentially 10) API calls. If you hit a 429 you will need to call this script again to get a new token.

### Adding new credentials

In order to add new credentials to the pool you will need to create a new Github App and add it to the "pool"

#### Creating a new App

Create a new app via the UI and use the same permissions as an existing one (e.g. [`getoutreach/ci-2](https://github.com/organizations/getoutreach/settings/apps/getoutreach-ci-2)).

Permissions *MUST* be the same otherwise flakes can occur.

Once this has been done, download the private key from the bottom of the app page.

Run the following command in your terminal to generate the env var string:

```bash
echo "YOUR_APP_ID,YOUR_INSTALL_ID,$(base64 < "$HOME/Downloads/your-key.pem")" | pbcopy
```

**Note**: This is now in your clipboard on macOS.

APP_ID comes from the main app page you downloaded the private key from, while install ID comes from the URL when you go to the installed app page for the app on `getoutreach` (e.g. [`getoutreach/ci-2](https://github.com/organizations/getoutreach/settings/installations/21145320))

Take the string you generated earlier (from your clipboard) and add it to the `ghaccesstoken` CircleCI
context as `GHACCESSTOKEN_GHAPP_<NUMBER>` replacing `<NUMBER>` with the next number in the "pool".

Done!
<!--- EndBlock(custom) -->

## Dependencies and Setup

### Dependencies

<!--- Block(dependencies) -->
{[] []}
<!--- EndBlock(dependencies) -->
