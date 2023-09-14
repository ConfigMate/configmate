# Contributing Guide

## Development

### Workflow
There is only one main branch in this project: `master`. Code in `master` should be functional at all times.
The development workflow to follow is the following:
1. Update your `master` branch.
2. Create a new branch out of `master`. To determine your branch name see the branch naming convention section below.
3. Once you have a functional and complete contribution, you make a pull request into `mater`. To determine how to name your pull request see the versioning section below. In the most common scenario, the correct naming would be: `<versioning-prefix>/<small-description-of-the-contribution>` where `<versioning-prefix>` will be `patch`, `feature` or `major`.

When the state of the `master` branch is deemed production ready, we create a branch out of `master` to make a release. In that branch we update the Changelog with the information of the new version following the format:
```
## [<version>] - <date> 
### Added
- ...
- ...

### Fixed
- ...
- ...

### Removed
- ...
- ...
```
For the version we just write the number (e.i. 2.3.7) and for the date we write MM-DD-YYYY.

After that change we make a pull request to `master` named: `release/v<version>`.

Everytime you merge a PR onto `master`, the version will be changed based on the conventions outlined and a tag will be created with that version name. You can use that tag to make a Github release or pre-release.

### Branch Naming Convention
All branches should be named as follows: `<name-of-developer>/i<issue-number>-description`. Skip issue number if the branch doesn't address a filed issue.

### Versioning and PR Naming Convention
Versions are automatically increased when a pull request is merged into `master`. The correct new version is determined by the prefix in the name of the pull request. Note that PRs will not be allowed to be merged if a supported prefix is not used.
#### Prefixes
- `patch/`: the patch version will be increased.
- `feature/`: the minor version will be increased.
- `major/`: the major version will be increased.
- `release/`: will remove the `beta` tag to indicate production readiness.
- Anything else: nothing will occur.
#### Pull requests where a version increase is not desired:
There are situations where changes don't affect the code of ConfigMate and thus shouldn't generate version changes. This can be acomplished by not using one of the prefixes above. As a best practice, we still want to prefix some of the most common cases of these.
##### Some non-versioning prefixes
- `devops/`: used when changes are made to the Makefile, the CI/CD pipeline, or any other items that contitute part of the building, packing or testing procedure (this list is not exhaustive).
- `docs/`: used when changes only include writing, deletion, or editing of docs such as the Readme.