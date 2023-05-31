# Checkmate 🕵️

## Automated checklist inspector for your pull requests

## Features

- [x] Fails validation until all configured checklists in the pull request description are checked.
- [x] Automatic checklists triggered on files modified in the pull request.
- [x] Support for validating select lists i.e. radio button.
- [x] [Job summary](https://github.blog/2022-05-09-supercharging-github-actions-with-job-summaries/) workflow report.

## Configuration

### Pull Request Description Checks

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


#### Oasis or Blur?
This select list validates that only one item is selected.

<!--Checkmate select=1-->
- [ ] Oasis
- [ ] Blur

```

Then configure the action in your workflow like

```yaml
on:
  pull_request:
    types: [edited, opened, reopened]

# cancel old edit events being processed
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true
  
name: Checkmate
jobs:
  validate-checklists:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Validate Checklists
      uses: marcusvnac/checkmate-evo@main
```

### Automatic checklists

Similar to the [workflow triggers on file paths](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#example-including-paths), 
you can configure Checkmate to write a checklist comment on a pull request which will be validated along with the
pull request description.

To enable the automated checklists configure the `with.paths` and the `with.github_token` inputs, and add the synchronize event for pull_request
and listen to all issue_comment events.

If you want to use a Selectlist then the first item should be the `<!--Checkmate select=1-->` comment.

```yaml
on:
  pull_request:
    types: [edited, opened, reopened, synchronize]
  issue_comment:

# cancel old edit events being processed
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true
  
name: Checkmate
jobs:
  validate-checklists:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Validate Checklists
        uses: marcusvanc/checkmate-evo@main
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          paths: |
            docs/**/*.md:
              - I have got the grammar police to proofread my changes.
            schema/migrations/*.sql:
              - There are no breaking changes in these migrations
              - I have notified X team of the new schema changes
            database/*.go:
              - <!--Checkmate select=2-->
              - Strong Consistency
              - Availability
              - Partition Tolerance
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

### Select 2 for files matching *database/\*.go*
<!-- Checkmate select=2 filepath=database/*.go -->
- [ ] Strong Consistency
- [ ] Availability
- [ ] Partition Tolerance
```
