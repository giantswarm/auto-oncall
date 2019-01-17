[![CircleCI](https://circleci.com/gh/giantswarm/auto-oncall.svg?&style=shield&circle-token=03d2446a590221c8604e7f540ea4cd1a4b929f14)](https://circleci.com/gh/giantswarm/auto-oncall)

# auto-oncall
auto-oncall application is a webhook handler, responsible for creating new Opsgenie routing rules on every merge event into master branch.

# configuration
Configuration requires next data to be configured in `values.yaml` of the helm chart:

```
# opsgenie api token
opsgenieToken:

# list of applications, configured for automated oncall rules
repositories:
  - test-oncall

# user mapping between github login and Opsgenie login 
users:
  github_user: user@giantswarm.io

# organization github webhook secret
githubWebhookSecret: 
```
