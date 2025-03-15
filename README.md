## go System Monitoring API (goSMA)

### About
This RESTful API allows for fetching metrics like CPU usage or network throughput of the host system.

![A diagram depicting goSMA's functionality](assets/goSMA.png)

### How to run
There's no binary available at the moment. Use `go run .` or utilize docker (`docker compose up --build`) to run goSMA.

### Configuration
goSMA makes use of environment variables for configuration:

| Environment variable | Description                                                                                                                    |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| `PORT`               | Changes the port goSMA runs on. If you're using docker, you can configure this via the `ports:` section in the `compose.yml`.  |
| `UPDATE_INTERVALL`   | Changes the time between updates to the data that's being made available through the API. Time is in seconds (float possible). |
| `HOSTNAME`           | Overwrites the hostname that is being read from `/etc/hostname`.                                                               |
