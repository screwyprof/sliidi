# Development Log

## Setup local environment with all the necessary tools
I've added a `Makefile` to help with routine actions such as tests, formatting code and linters. 
Run `make help` to get a list of all possible targets.

## Setup CI Pipeline
One of the first things I usually do when I build a project from scratch is set up the CI pipeline as early as possible.
The tests must run on every commit. Apart from that I would love to run linters as well. For the sake of this demo,
I decided to use Github Actions to help me run the builds.
