# auto-oncall
auto-oncall application is a webhook handler, responsible for creating new Opsgenie routing rules on every deployment event.

# configuration
Configuration requires next data to be configured in `values.yaml` of the helm chart:

```
# github token with private repositories read access
githubToken: 
# opsgenie api token
opsgenieToken:

# user mapping between github login and Opsgenie login 
users:
  github_user: user@giantswarm.io

# organization github webhook secret
githubWebhookSecret: 
```
