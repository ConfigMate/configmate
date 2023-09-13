# ConfigMate - Your Configs' Best Friend 🤝

ConfigMate is a tool designed to validate configuration files against custom sets of rules. The goal is to reduce errors and streamline the validation process, thereby alleviating the mental burden on developers and improving operational efficiency for businesses.


## Table of Contents

- [Development](#development)
- 🌟 [Features](#features)
- 👏 [Contributors](#contributors)
- :memo: [License](#license)

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

## Features
- **Custom Rule-Based Validation**: Define your own custom rules to validate various configuration fields.
- **Pre-defined Rulesets**: Comes with a set of commonly used rules to get you started.
- **Multi-file & Multi-field Checks**: Ability to apply rules that span multiple configuration files and fields.
- **CLI Interface**: Easily run configuration checks directly from your command line.
- **VS Code Extension**: Get in-editor validation and error highlighting within Visual Studio Code.
- **Meaningful Error Messages**: Receive detailed error descriptions and exactly where the issue occurs (file, line number, and column).
- **ANTLR-Based Custom Parsers**: Utilizes ANTLR for specialized parsing to validate files beyond simple syntax checks.


## Contributors
Thanks go to these wonderful people

<table>
  <tr>
    <td align="center">
      <a href="https://github.com/javier-arango" target="_blank">
        <img src="https://avatars.githubusercontent.com/u/58098790?s=60" width="60px;"/><br />
        <sub><b>Javier Arango</b></sub>
      </a><br />
    </td>
    <td align="center">
      <a href="https://github.com/Jcabza008" target="_blank">
        <img src="https://avatars.githubusercontent.com/u/34218922?s=60" width="60px;"/><br />
        <sub><b>Julio J. Cabrera</b></sub>
      </a><br />
    </td>
    <td align="center">
      <a href="https://github.com/jeangregorfonrose" target="_blank">
        <img src="https://avatars.githubusercontent.com/u/21975726?s=60" width="60px;"/><br />
        <sub><b>Jean Gregor Fonrose</b></sub>
      </a><br />
    </td>
    <td align="center">
      <a href="https://github.com/ktminks" target="_blank">
        <img src="https://avatars.githubusercontent.com/u/19628386?s=60" width="60px;"/><br />
        <sub><b>Katie Minckler</b></sub>
      </a><br />
    </td>
  </tr>
</table>


## License
[MIT](https://github.com/ConfigMate/configmate/blob/master/LICENSE)
