# Checkmate 🕵️

## Automated checklist inspector for your pull requests

## Features

- [x] Fails validation until all configured checklists are checked.
- [ ] Support for validating one-of lists i.e. radio button
- [x] Automatically creates checklists based on file changes.

## Configuration

### PR Description Checks

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
  pull_request:
    types: [edited, opened, reopened]
    
name: Checklist Check
jobs:
  checkmate:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Validate Checklists
      uses: roryq/checkmate@master
```

### Automated checklists

To enable the automated checklists configure the paths and the github_token inputs, and add the synchronize event.

```yaml
on:
  pull_request:
    types: [edited, opened, reopened, synchronize]
  issue_comment:
    types: [edited, created, deleted]

name: Checklist Check
jobs:
  checkmate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Validate Checklists
        uses: roryq/checkmate@master
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          paths: |
            docs/**/*.md:
              - I have got the grammar police to proofread my changes.
            schema/migrations/*.sql:
              - There are no breaking changes in these migrations
              - I have notified X team of the new schema changes
```

Which will create and update a comment when the changeset files match the configured patterns.


```markdown
Thanks for your contribution!
Please complete the following tasks related to your changes and tick the checklists when complete.

### Checklist for files matching *docs/\*\*/\*.md*
<!-- Checkmate filepath=docs/**/*.md -->
- [ ] I have got the grammar police to proofread my changes.

### Checklist for files matching *schema/migrations/\*.sql*
<!-- Checkmate filepath=schema/migrations/*.sql -->
- [ ] There are no breaking changes in these migrations
- [ ] I have notified X team of the new schema changes
```
