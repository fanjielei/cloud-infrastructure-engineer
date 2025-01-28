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

- [ ] Inspect the Go code and get a rough understanding of how the server works
- [ ] Use Docker to build and start the server
- [ ] Create a `docker-compose.yml` file to orchestrate the services you will create
- [ ] Create a client program in a language of your choice
- [ ] Your client should run periodic requests against the status endpoint and report:
  - Which endpoint it talked to
  - How long the request took
  - If the request was successful
  - Any errors it encountered
- [ ] Your client should periodically tell the server to change the response via the dedicated endpoint
- [ ] Add Prometheus and scrape the server's metrics endpoint

## Optional tasks

If you want to further explore the tasks, here are some inspirations:

- [ ] Add a Grafana dashboard showing graphs, e.g. for status codes and response times
- [ ] Add a structured logging solution of your choice and make it queriable via Grafana.
- [ ] Add Tempo to visualise the traces from the OpenTelemetry instrumentation

**Again**, those are optional tasks, so please do not feel obligated to spend more time on them than you like.
