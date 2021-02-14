# Development Log

## Understanding the task
As almost any non-trivial tasks this one requires some understanding of the business domain. I would normal spent some
time to understand the actual requirements for the system and to discuss some pros and cons.

For example, in this case I would ask how the system is going to be used. What kind of content is expected? What is the
expected load profile? If it's a real-time highroad system to get the current quotes for an Exchange it's one thing.
If it's just a handful of users reading news - that's another thing.

Answering these questions would have impact on the actual implementation. If it's not a real-time system, then it doesn't
make a lot of sense to query the providers on the go. It would be more reasonable to create a worker or a background task
which would grab the content once in a given interval and store somewhere for a quick access.

If we're speaking about a high-load project, we should think about caching the responses. We should also think about
some things like rate-limiting the amount of queries, circuit-breakers and some smarter load-balancing techniques.

Given that I don't have a chance to ask question at the moment, and I don't really now what is expected, I decided to
follow the given requirements literally. I would be happy to discuss the solution as well as my vision on the task.

## Setup local environment with all the necessary tools
I've added a `Makefile` to help with routine actions such as tests, formatting code and linters.
Run `make help` to get a list of all possible targets.

## Setup CI Pipeline
One of the first things I usually do when I build a project from scratch is set up the CI pipeline as early as possible.
The tests must run on every commit. Apart from that I would love to run linters as well. For the sake of this demo,
I decided to use Github Actions to help me run the builds.

## Coding

When it comes to actually writing the code I practice A-TDD. Usually I start with an acceptance test to capture the
expectations. In this case however the sample goes with a test which can be used for this purpose. I just renamed it to
`api_acceptance_test.go` to emphasize its meaning.

After that I started working on the task following the Red-Green-Refactor TDD pattern. I tried to make as less changes
as possible to make the code easier to read and maintain. I almost haven't touch the original code base, but I can see
quite a big room for improvement there.

## Concurrency
The task has it that:
> Latency is crucial for this application, so fetching the items sequentially one at a time might not be good enough.
Which probably implies that the actual communication with the providers should be done concurrently.

It definitely makes sense, however this requirement contradicts to another one:
> In the case both the main provider and the fallback fail (or if the main provider fails and there is no fallback),
> the API should respond with all the items before that point. So, for example, if the configuration calls for [1,1,2,3]
> and 2 fails, the response should only contain [1,1]

We cannot guarantee the order of concurrent jobs, so this requirement complicates the implementation. If it were to me
I would change it to something more reasonable. For example, return all the positive responses. This would make the code
much less convoluted and efficient.

When it comes to the implementation, I used some very simple concurrency patterns on top of the standard library.
In a real working project it would make more sense to use some robust implementations for job queues and workers.
