# snap collector plugin - gitstats

This plugin collects GitHub metrics.

It's used in the [snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [Installation](#installation)
  * [Configuration and Usage](configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
3. [License](#license-and-authors)

## Getting Started

### Installation

#### Download plugin binary

You can get the pre-built binaries at [GitHub Releases](https://github.com/grafana/snap-plugin-collector-gitstats/releases) page.

#### To build the plugin binary

Fork https://github.com/grafana/snap-plugin-collector-gitstats and clone repo into `$GOPATH/src/github.com/grafana/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-gitstats.git
```

Build the plugin by running make within the cloned repo:

```bash
$ ./build.sh
```

This builds the plugin binary in `/build/`

This plugin uses govendor to manage dependencies. If you want to add a dependency, then:

1. Install govendor with: `go get -u github.com/kardianos/govendor`
2. `govendor fetch <dependency path>`
3. `govendor install` to update the vendor.json file.
4. Check in the new dependency that will be in the vendor directory.

### Configuration and Usage

* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Ensure `$SNAP_PATH` is exported
`export SNAP_PATH=$GOPATH/src/github.com/intelsdi-x/snap/build`
* Three config variables must be set:
  - `access_token` is your private GitHub token.
  - `user` is your GitHub org  
  - `repo` is the GitHub repo to collect stats from. Can be one repo e.g. grafana or * for all repos for the specified user.

  Example of how to configure it in a json task manifest:
  ```json
  {
    "version": 1,
    "schedule": {
      "type": "simple",
      "interval": "1h"
    },
    "deadline": "5m",
    "workflow": {
      "collect": {
        "metrics": {
          "/raintank/apps/gitstats/repo/*":{}
        },
        "config": {
          "/raintank/apps/gitstats": {
            "access_token": "your_private_github_token",
            "repo": "grafana",
            "user": "grafana"
          }
        },
        "process": null,
        "publish": [
          {
            "plugin_name": "graphite",
            "config": {
              "prefix_tags": "",
              "prefix": "",
              "server": "127.0.0.1"
            }
          }
        ]
      }
    }
  }
  ```

## Documentation

### Collected Metrics

This plugin has the ability to gather the following metrics:

TODO (as they are probably going to change)

## License

This plugin is released under the Apache 2.0 [License](LICENSE).
