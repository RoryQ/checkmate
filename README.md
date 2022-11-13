# Checkmate 🕵️

## Automated checklist inspector for your pull requests

## Features

- [x] Fails validation until all configured checklists are checked.
- [ ] Support for validating one-of lists i.e. radio button
- [x] Automatically creates checklists based on file changes.

## Configuration

In your PR description or issue template, place `<!--Checkmate-->` above the checklist block you want validated. A block of checklist
items is one without empty lines in-between.

```markdown
#### :heavy_check_mark: Checklist

<!--Checkmate-->

- [ ] Added or updated documentation
- [x] Tests for new functionality and regression tests for bug fixes

#### This checklist will not be validated

- [ ] This is ignored
- [ ] This is ignored
- [ ] This is ignored
```

Then configure the action in your workflow like

```yaml
on:
  # Run against PR description changes
  pull_request:
    types: [edited, opened, reopened]

name: Checklist Check
jobs:
  checkmate:
    runs-on: ubuntu-latest
    steps:
    - name: validate
      uses: roryq/checkmate@master
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```