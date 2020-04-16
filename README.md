# gotopus

[![codecov](https://codecov.io/gh/lherman-cs/gotopus/branch/master/graph/badge.svg)](https://codecov.io/gh/lherman-cs/gotopus)
[![Go Report Card](https://goreportcard.com/badge/github.com/lherman-cs/gotopus)](https://goreportcard.com/report/github.com/lherman-cs/gotopus)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/81a94fecb21b45bcb89ef6b8c6e3b682)](https://www.codacy.com/manual/lherman-cs/gotopus?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=lherman-cs/gotopus&amp;utm_campaign=Badge_Grade)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Gotopus is a minimalistic tool that runs arbitrary commands concurrently. You define your commands with their dependencies and Gotopus will take care of the rest, running them concurrently when possible.

## Contents
- [Features](#features)
- [Installation](#installation)
- [Getting Started](#getting-started)
  - [Basic Usage](#basic-usage)
  - [Environment Variables](#environment-variables)
  - [Concurrency vs Parallelism](#concurrency-vs-parallelism)
- [FAQ](#faq)
  - [Why does the config format look similar to Github Actions](#why-does-the-config-format-look-similar-to-github-actions)

## Features

- [X] Concurrently run steps, speeding up running time
- [X] Local or remote configs
- [X] Easy to install
- [X] Circular dependency detection
- [X] Clean step definition with [YAML](https://en.wikipedia.org/wiki/YAML)
- [X] [Builtin and user environment variables](#environment-variables)

## Installation

```sh
curl -sf https://gobinaries.com/lherman-cs/gotopus | sh
```

## Getting Started

### Basic Usage

```
Usage: gotopus <url or filepath> ...

  -max_workers uint
    	limits the number of workers that can run concurrently (default 0 or limitless)
```

```yaml
# examples/basic.yaml
jobs:
  job1:
    steps:
      - run: sleep 1 && echo "job1"
  job2:
    needs:
      - job1
    steps:
      - run: echo "job2"
  job3:
    steps:
      - run: echo "job3"
```

To use `basic.yaml` above, you can run the following command:

```sh
gotopus basic.yaml
```

Or you can simply give a URL to this file:

```sh
gotopus https://raw.githubusercontent.com/lherman-cs/gotopus/master/examples/basic.yaml
```

By default, if you don't set `max_workers` to any number greater than 0, gotopus will create a pool of workers without a limit in lazy way. From the example above, 2 workers will be allocated instead of 3. The process looks like following:

```
job1 gets scheduled
spawn worker #0
worker #0 executes job1
job3 gets scheduled
spawn worker #1
worker #1 executes job3
job3 finishes
job1 finishes
either worker #0 or #1 executes job2
job2 finishes
```

### Environment Variables
Whenever a step runs, there are 3 kinds of environments that are going to be set and they'll have the priority order (in case of a conflict happens, the higher priority environment variable will be chosen) as listed below, where user environment variables will have the highest priority:

* User: these environment variables are defined by the user in yaml in each step.
* Builtin: environment variables that come from gotopus and they'll be prefixed with `GOTOPUS_`.

  * `GOTOPUS_JOB_ID`
  * `GOTOPUS_JOB_NAME`
  * `GOTOPUS_STEP_NAME`
  * `GOTOPUS_WORKER_ID`

* System: inherits all the environments variables from the system when you run gotopus.

Following is an example how you define and use environment variables:

```yaml
# examples/env.yaml
jobs:
  job:
    steps:
      - name: Install dependencies
        run: echo "$GOTOPUS_STEP_NAME"
      - run: echo "$name"
        env:
          name: Lukas Herman
```

### Concurrency vs Parallelism
Let's imagine that there are 2 commands that we want to execute:

1. First command: `sleep 2 && echo "job 1"`
2. Second command: `sleep 3 && echo "job 2"`

Normally, this would take **5 seconds** to finish since you need to run them sequentially. But, if you run this concurrently, this would take **~3 seconds even when you only have 1 CPU core!**. This is possible because they run **conccurrently NOT in parallel**. For more information, there's this awesome video from Rob Pike that specfically talks about "Concurency Is Not Parallelism".

[![Conccurency vs Parallelism](https://img.youtube.com/vi/cN_DpYBzKso/0.jpg)](https://www.youtube.com/watch?v=cN_DpYBzKso)

```sh
time gotopus https://raw.githubusercontent.com/lherman-cs/gotopus/master/examples/concurrency.yaml
```

## FAQ

### Why does the config format look similar to Github Actions?
This is intentional because I think Github Actions format is pretty clean and also one of my inspirations to create Gotopus. Although Gotopus and Github Actions seem to have similar functionalities, you can definitelly use them together! For example, you can run steps of a job concurrently with Gotopus (running steps concurrently is not supported currently: https://github.community/t5/GitHub-Actions/Steps-in-parallel/td-p/32635).
