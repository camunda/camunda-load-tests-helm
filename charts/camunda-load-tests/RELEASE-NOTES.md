The changelog is automatically generated using [git-chglog](https://github.com/git-chglog/git-chglog)
and it follows [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format.


<a name="camunda-load-tests-0.1.6"></a>
## [camunda-load-tests-0.1.6](https://github.com/camunda/camunda-platform-helm/compare/camunda-load-tests-0.1.5...camunda-load-tests-0.1.6) (2026-04-08)

### Docs

* add Spring Boot env var examples to values.yaml

### Feat

* allow to configure nodeSelector + tolerations
* support configurable rateDuration for starter
* add performReadBenchmarks as a first-class Helm value
* add Spring Boot env vars to worker template
* add Spring Boot env vars to starter template
* restart deployments when configuration or secrets change

### Fix

* compute config checksum based on the config content directly

### Test

* add property tests for Spring Boot env var mapping
* add golden test for worker with message configuration
* regenerate golden files with Spring Boot env vars

### Pull Requests

* Merge pull request [#37](https://github.com/camunda/camunda-platform-helm/issues/37) from camunda/cz-spring-boot-compat
* Merge pull request [#35](https://github.com/camunda/camunda-platform-helm/issues/35) from camunda/release

