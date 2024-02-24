# Event Scheduler

This my first (little bit larger than hello-world) go-project.
Disclaimer: May contain a lot of unconventional go code.

## Steps to run

```
go build
./goproject
```
Starts listening http://localhost:8888
## API

* POST api/v1/events
* POST api/v1/events/{event-id}/vote
* GET api/v1/events
* GET api/v1/events/{event-id}
* GET api/v1/events/{event-id}/results
