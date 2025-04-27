# Cloud Infrastructure Engineer Tasks

## Goal

We want you to create a monitoring app that will send requests to our web server and display whether there are any problems with it. Don't panic, you do not have to reinvent Prometheus to pass the exercise. The goal is to get to know your work and style and have a technical conversation in your second interview.

## What we want to receive from you

- Do not spend more than 2-3 hours in total.
- Commit your progress, so we can follow your journey in the Git history.
- Make sure your commits are concise and well structured.
- All source files / scripts needed to start your stack should be contained in the Git repository.
- Add necessary documentation to *this* `README.md` file so we can reproduce your setup locally.
- Let us know which platform / operating system you built / tested your setup with.
- You should be able to introduce your setup in a 30 minute interview slot.
- Create a `tar.gz` or `.zip` archive of your Git repository and email it to us.

**You can take inspiration, but please do not copy & paste "solutions" from ChatGPT.**

## Tasks

- [x] Inspect the Go code and get a rough understanding of how the server works

> `post /flaky` to switch `flaky` on/off

> if `flaky` it `ture`, `get /status` will get a reponse with a random status code in randomly respond time less than 500 milisecond

- [X] Use Docker to build and start the server: `docker build -f cmd/server/Dockerfile .`

> I update the Dockerfile to generate an image with smaller image size

- [X] Create a `docker-compose.yml` file to orchestrate the services you will create
- [X] Create a client program in a language of your choice

> file: `cmd/client/main.go`

> run it: `go run cmd/client/main.go`

- [X] Your client should run periodic requests against the status endpoint and report:
  - Which endpoint it talked to
  - How long the request took
  - If the request was successful
  - Any errors it encountered

> see stdout of `go run cmd/client/main.go`

- [X] Your client should periodically tell the server to change the response via the dedicated endpoint
> implement it in the client.go

> I switch on `flaky`, then for each `get /status`, the server will change the response. Is this what this question want?

- [ ] Add Prometheus and scrape the server's metrics endpoint

## Optional tasks

If you want to further explore the tasks, here are some inspirations:

- [ ] Add a Grafana dashboard showing graphs, e.g. for status codes and response times
- [ ] Add a structured logging solution of your choice and make it queriable via Grafana.
- [ ] Add Tempo to visualise the traces from the OpenTelemetry instrumentation

**Again**, those are optional tasks, so please do not feel obligated to spend more time on them than you like.
