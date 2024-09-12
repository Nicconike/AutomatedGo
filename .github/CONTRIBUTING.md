# Contributing

✨ Thanks for contributing to **AutomatedGo**! ✨

[fork]: https://github.com/nicconike/automatedgo/fork
[pr]: https://github.com/nicconike/automatedgo/compare
[code-of-conduct]: CODE_OF_CONDUCT.md

As a contributor, here are the guidelines I would like you to follow:

- [Contributing](#contributing)
	- [Contributing to AutomatedGo](#contributing-to-automatedgo)
	- [Code of conduct](#code-of-conduct)
	- [How can I contribute?](#how-can-i-contribute)
		- [Improve documentation](#improve-documentation)
		- [Give feedback on issues](#give-feedback-on-issues)
		- [Fix bugs and implement features](#fix-bugs-and-implement-features)
	- [Using the issue tracker](#using-the-issue-tracker)
		- [Bug report](#bug-report)
		- [Feature request](#feature-request)
	- [Submitting a pull request](#submitting-a-pull-request)
	- [Coding rules](#coding-rules)
		- [Source code](#source-code)
		- [Documentation](#documentation)
		- [Commit message guidelines](#commit-message-guidelines)
			- [Atomic commits](#atomic-commits)
			- [Commit message format](#commit-message-format)
			- [Lint](#lint)
				- [Staticcheck](#staticcheck)
				- [Other Lint Tools](#other-lint-tools)
	- [Resources](#resources)

## Contributing to AutomatedGo

Contributions to this project are [released](https://help.github.com/articles/github-terms-of-service/#6-contributions-under-repository-license) to the public under the [project's open source license](LICENSE).

Please note that this project is released with a [Contributor Code of Conduct][code-of-conduct]. By participating in this project you agree to abide by its terms.

## Code of conduct

Help me keep **AutomatedGo** open and inclusive. Please read and follow the [Code of conduct](CODE_OF_CONDUCT.md).

## How can I contribute?

### Improve documentation

As a **AutomatedGo** user, you are the perfect candidate to help me improve the documentation: typo corrections, clarifications, more examples, etc. Take a look at the [documentation issues that need help](https://github.com/Nicconike/automatedgo/issues/new?assignees=nicconike&labels=documentation&projects=&template=documentation.yml&title=%5BDOC%5D+).

Please follow the [Documentation guidelines](#documentation).

### Give feedback on issues

Some issues are created without the information requested in the [Bug report guideline](#bug-report).
Help make them easier to resolve by adding any relevant information.

Issues with the [enhancement label](https://github.com/Nicconike/automatedgo/labels/enhancement) are meant to discuss the implementation of new features.
Participating in the discussion is a good opportunity to get involved and influence the future direction of **AutomatedGo**.

### Fix bugs and implement features

Confirmed bugs and ready-to-implement features can be marked with the [help wanted label](https://github.com/Nicconike/automatedgo/labels/help%20wanted).
Post a comment on an issue to indicate you would like to work on it and to request help from the [me](https://github.com/Nicconike) and the community.

## Using the issue tracker

The issue tracker is the channel for [bug reports](#bug-report), [features requests](#feature-request), [submitting pull requests](#submitting-a-pull-request) and many more.

Before opening an issue or a Pull Request, please use the [GitHub issue search](https://github.com/Nicconike/automatedgo/issues?q=is%3Aissue) to make sure the bug or feature request hasn't been already reported or fixed.

### Bug report

A good bug report shouldn't leave others needing to chase you for more information.
Please try to be as detailed as possible in your report and fill the information requested in the [bug report template](https://github.com/Nicconike/automatedgo/issues/new?assignees=nicconike&labels=bug&projects=&template=bug-report.yml&title=%5BBUG%5D+).

### Feature request

Feature requests are welcome, but take a moment to find out whether your idea fits with the scope and aim of the project.
It's up to you to make a strong case to convince the project's developers of the merits of this feature.
Please provide as much detail and context as possible and fill the information requested in the [feature request template](https://github.com/Nicconike/automatedgo/issues/new?assignees=nicconike&labels=enhancement&projects=&template=feature-request.yml&title=%5BFEATURE%5D+).

## Submitting a pull request

1. [Fork][fork] and clone the repository
2. Configure and install the dependencies: `pip install`
3. Create a new branch: `git checkout -b my-branch-name`
4. Make your change, add tests, and make sure the tests still pass
5. Push to your fork and [submit a pull request][pr]
6. Pat your self on the back and wait for your pull request to be reviewed and merged.

Here are a few things you can do that will increase the likelihood of your pull request being accepted:

- Follow the [guide](#how-can-i-contribute).
- Write tests.
- Keep your change as focused as possible. If there are multiple changes you would like to make that are not dependent upon each other, consider submitting them as separate pull requests.
- While doing the commits, please strictly follow [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/).

## Coding rules

### Source code

To ensure consistency and quality throughout the source code, all code modifications must have:

- No [linting](#lint) errors
- A [test](#tests) for every possible case introduced by your code change
- **100%** test coverage
- [Valid commit message(s)](#commit-message-guidelines)
- Documentation for new features
- Updated documentation for modified features

### Documentation

To ensure consistency and quality, all documentation modifications must:

- Refer to brand in [bold](https://help.github.com/articles/basic-writing-and-formatting-syntax/#styling-text) with proper capitalization, i.e. **GitHub**, **automatedgo**,etc.
- Prefer [tables](https://help.github.com/articles/organizing-information-with-tables) over [lists](https://help.github.com/articles/basic-writing-and-formatting-syntax/#lists) when listing key values, i.e. List of options with their description
- Use [links](https://help.github.com/articles/basic-writing-and-formatting-syntax/#links) when you are referring to:
  - A **automatedgo** concept described somewhere else in the documentation, i.e. How to [contribute](CONTRIBUTING.md)
  - A third-party product/brand/service, i.e. Integrate with [GitHub](https://github.com)
  - An external concept or feature, i.e. Create a [GitHub release](https://help.github.com/articles/creating-releases)
- Use the [single backtick `code` quoting](https://help.github.com/articles/basic-writing-and-formatting-syntax/#quoting-code) for:
  - programming language keywords, i.e. `function`, `async`, `String`
- Use the [triple backtick `code` formatting](https://help.github.com/articles/creating-and-highlighting-code-blocks) for:
  - code examples
  - configuration examples
  - sequence of command lines

### Commit message guidelines

#### Atomic commits

If possible, make [atomic commits](https://en.wikipedia.org/wiki/Atomic_commit), which means:

- a commit should contain exactly one self-contained functional change
- a functional change should be contained in exactly one commit
- a commit should not create an inconsistent state (such as test errors, linting errors, partial fix, feature without documentation, etc...)

A complex feature can be broken down into multiple commits as long as each one maintains a consistent state and consists of a self-contained change.

#### Commit message format

Each commit message consists of a **header**, a **body** and a **footer**.
The header has a special format that includes a **type**, a **scope** and a **subject**:

```commit
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

The **header** is mandatory and the **scope** of the header is optional.

The **footer** can contain a [closing reference to an issue](https://help.github.com/articles/closing-issues-via-commit-messages).

Please check [COMMITS.md](https://github.com/Nicconike/Steam-Stats/blob/master/COMMITS.md) for a detailed explanation about the commit message format and [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/).

#### Lint

##### Staticcheck

[AutomatedGo](https://github.com/Nicconike/automatedgo) repository uses [staticcheck](https://staticcheck.dev/) for linting. And below other tools:
- gopls
- gotests
- gomodifytags
- impl
- goplay
- dlv

##### Other Lint Tools

`golangci-lint` is often recommended as it combines multiple linters and offers good performance.

Other tools are used for specific compatibility concerns, but are less likely to result in failures in common contributions.
Please follow the guidance of these tools if failures are encountered.

## Resources

- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Using Pull Requests](https://help.github.com/articles/about-pull-requests/)
- [GitHub Help](https://help.github.com)
