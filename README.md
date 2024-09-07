# ghcr-cli

`ghcr-cli` is a command-line interface tool that serves as a wrapper around the GitHub CLI (`gh`) for managing Docker images in GitHub Container Registry (ghcr.io). It provides an easy-to-use interface for listing and deleting container images.

## Note

`ghcr-cli` is a wrapper around the GitHub CLI and relies on it for authentication and API calls. Make sure your GitHub CLI is always up to date and properly authenticated with the necessary permissions.

## Features

- List Docker images in a GitHub Container Registry
- Delete specific versions of Docker images
- JSON output for easy parsing and integration with other tools

## Tested GitHub CLI Version

`ghcr-cli` has been tested with the following versions of GitHub CLI.

```
gh version 2.46.0 (2024-03-26 Debian 2.46.0-1)
https://github.com/cli/cli/releases/tag/v2.46.0
```

## Installation

Build the project:

```
$ go build -o ghcr-cli
```

## Authentication and Permissions

Before using `ghcr-cli`, you need to authenticate with GitHub and ensure you have the necessary permissions. Follow these steps:

1. Authenticate with GitHub CLI:
   ```
   $ gh auth login
   ```

2. If you haven't already, add the necessary scopes for package management:
   ```
   $ gh auth login -s 'read:packages delete:packages'
   ```

   This command adds the `read:packages` and `delete:packages` scopes, which are required for listing and deleting packages respectively.

   https://docs.github.com/en/packages/learn-github-packages/about-permissions-for-github-packages#about-scopes-and-permissions-for-package-registries

## Usage

`ghcr-cli` provides a simple interface to interact with GitHub Container Registry. All commands require specifying the owner and repository. The tool uses GitHub CLI (`gh`) under the hood, so ensure you're authenticated before running any commands.

### Listing Images

To list all Docker images in a specific repository:

```
$ ghcr-cli list OWNER/REPO
```

This will output a JSON array containing information about each image version.

### Deleting an Image Version

To delete a specific version of a Docker image:

```
$ ghcr-cli delete OWNER/REPO VERSION_ID
```

Replace `OWNER/REPO` with the appropriate GitHub username/organization and repository name, and `VERSION_ID` with the ID of the version you want to delete.

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
