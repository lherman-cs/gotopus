# gotopus

[![codecov](https://codecov.io/gh/lherman-cs/gotopus/branch/master/graph/badge.svg)](https://codecov.io/gh/lherman-cs/gotopus)
[![Go Report Card](https://goreportcard.com/badge/github.com/lherman-cs/gotopus)](https://goreportcard.com/report/github.com/lherman-cs/gotopus)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/81a94fecb21b45bcb89ef6b8c6e3b682)](https://www.codacy.com/manual/lherman-cs/gotopus?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=lherman-cs/gotopus&amp;utm_campaign=Badge_Grade)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Gotopus is a minimalistic tool that runs arbitrary commands concurrently. You define your commands with their dependencies and Gotopus will take care of the rest, running them concurrently when possible.

# Installation

```sh
curl -sf https://gobinaries.com/lherman-cs/gotopus | sh
```

# Features

- [X] Concurrently run steps, speeding up running time
- [X] Local or remote configs
- [X] Easy to install
- [X] Circular dependency detection
- [X] Clean step definition with [YAML](https://en.wikipedia.org/wiki/YAML)
