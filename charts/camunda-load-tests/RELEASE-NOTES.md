The changelog is automatically generated using [git-chglog](https://github.com/git-chglog/git-chglog)
and it follows [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format.


<a name="camunda-load-tests-0.1.13"></a>
## [camunda-load-tests-0.1.13](https://github.com/camunda/camunda-load-tests-helm/compare/camunda-load-tests-0.1.12...camunda-load-tests-0.1.13) (2026-09-03)

### Feat

* propagate commonLabels to all resource templates
* add global.commonLabels helper to chart

### Fix

* suppress blank lines when commonLabels is empty
* remove spurious leading newline from commonLabels helper

### Refactor

* remove outdated configurations
* make max default
* merge commonLabels into camunda-load-tests.labels template

### Test

* update golden files
* update golden
* update golden tests
* regenerate golden files after merging origin/main
* update golden files after merging label templates
* update golden files after empty-commonLabels blank-line fix
* update commonLabels golden files after leading-newline fix
* update existing golden files for commonLabels template changes
* add golden tests for commonLabels propagation

### Pull Requests

* Merge pull request [#73](https://github.com/camunda/camunda-load-tests-helm/issues/73) from camunda/release

