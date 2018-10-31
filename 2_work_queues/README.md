# Work Queues

The tutorial is [here](http://www.rabbitmq.com/tutorials/tutorial-two-go.html)

## Usage

The code is built with _Docker_, so nothing other than `docker` and `docker-compose` needed.

Run the services using `docker-compose --scale`;

```bash
$ docker-compose up --scale worker=10 --scale task_enqueuer=10
```

> Use `docker-compose`'s `--build` flag if you want to rebuild the images after changing the Go code.
> `docker-compose up --build --scale ...`
