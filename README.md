# Project for testing a potential memory leak in timescale db
The project is a small test program to show a potential memory leak with long lived connections in timescale db.
Written in Go, it uses identical set ups of a vanilla postgres and timescaleDBs docker images 
### Prerequisits:
* Docker
* Docker compose.

### Run:
```
docker-compose up test-timescale-memory
```

Running the program creates one connection to each database and inserts 2.5 million rows evenly over 5 days, 
default is that the data is inserted in a hypertable with 1 day as chunk-size for timescale, and a normal table for vanilla postgres.

After finishing the program holds the connection to each respective database to more clearly show the memory leak.

To monitor memory usage one can use htop, top or [ctop](https://github.com/bcicen/ctop). 
The easiest tool to use is probably ctop.