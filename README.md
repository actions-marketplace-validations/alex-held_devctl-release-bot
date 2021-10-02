[![Netlify Status](https://api.netlify.com/api/v1/badges/cfd72dea-e22a-463b-8e20-5748b743140a/deploy-status)](https://app.netlify.com/sites/angry-borg-f9dd47/deploys)

<a href="https://github.com/rajatjindal/devctl-release-bot"><img src="https://github.com/devctl-release-bot.png" width="100"></a><span width="10px">

`devctl-release-bot` is a bot that automates the update of plugin manifests in `devctl-index` when a new version of your `kubectl` plugin is released.
If a release is marked as a 'prerelease' in github, it will not be released to the devctl index.

To trigger `devctl-release-bot` you can use a `github-action` which sends the event to the bot.

# Basic Setup

- Make sure you have enabled github actions for your repo
- Add a `.devctl.yaml` template file at the root of your repo. Refer to [kubectl-evict-pod](https://github.com/rajatjindal/kubectl-evict-pod) repo for an example.
  - you could use https://rajatjindal.com/tools/devctl-release-bot-helper/ for generating template for your plugin
- To setup the action, add the following snippet after the step that publishes the new release and assets:
  ```yaml
  - name: Update new version in devctl-index
    uses: rajatjindal/devctl-release-bot@v0.0.40
  ```
  Check out the `goreleaser` example below for details.

##### Example when using go-releaser

`<your-git-root>/.github/workflows/release.yml`

```yaml
name: release
on:
  push:
    tags:
      - "v*.*.*"
jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@master
      - name: Setup Go
        uses: actions/setup-go@v1
        with:
          go-version: 1.16
      - name: GoReleaser
        uses: goreleaser/goreleaser-action@v1
        with:
          version: latest
          args: release --rm-dist
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      - name: Update new version in devctl-index
        uses: rajatjindal/devctl-release-bot@v0.0.40
```

\*\* You can also customize the release assets names, platforms for which build is done using .goreleaser.yml file in root of your git repo.

# Examples using devctl-release-bot in different ways

- [bash based plugins](https://github.com/ahmetb/kubectx/blob/master/.github/workflows/release.yml)
- [multiple plugins published from one repo](https://github.com/ahmetb/kubectx/blob/master/.github/workflows/release.yml)
- [circle-ci](examples/circleci.yml)
- [travis-ci](examples/travis.yml)

# Testing the template file

You can test the template file rendering before check-in to the repo by running following command

```bash
$ docker run -v /path/to/your/template-file.yaml:/tmp/template-file.yaml rajatjindal/devctl-release-bot:v0.0.40 \
  devctl-release-bot template --tag <tag-name> --template-file /tmp/template-file.yaml
```

# Inputs for the action

| Key                | Default Value          | Description                                                                          |
| ------------------ | ---------------------- | ------------------------------------------------------------------------------------ |
| workdir            | `env.GITHUB_WORKSPACE` | Overrides the GitHub workspace directory path                                        |
| devctl_template_file | `.devctl.yaml`           | The path to template file relative to $workdir. e.g. templates/misc/plugin-name.yaml |

# Limitations of devctl-release-bot

- only works for repos hosted on github right now
- The first version of plugin has to be submitted manually, by plugin author, to the devctl-index repo

# Kubernetes CLA

devctl-release-bot is just a service to open PR on your behalf to release a new version of the devctl-plugin. Your CLA agreement (that you did when submitting the new plugin to devctl-index) is still applicable on these PR's.
