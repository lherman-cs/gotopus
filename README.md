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

# Getting Started

## Concurrency vs Parallelism
Let's imagine that there are 2 commands that we want to execute:

1. First command: `sleep 2 && echo "job 1"`
2. Second command: `sleep 3 && echo "job 2"`

Normally, this would take **5 seconds** to finish since you need to run them sequentially. But, if you run this concurrently, this would take **~3 seconds even if you only have 1 CPU core!**. This is possible because they run **conccurrently NOT in parallel**. For more information, there's this awesome video from Rob Pike that specfically talks about "Concurency Is Not Parallelism".

[![Conccurency vs Parallelism](https://img.youtube.com/vi/cN_DpYBzKso/0.jpg)](https://www.youtube.com/watch?v=cN_DpYBzKso)

```sh
https://raw.githubusercontent.com/lherman-cs/gotopus/master/examples/concurrency.yaml
```
