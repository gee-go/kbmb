# Running
- Install docker and docker-compose (on a mac I like https://github.com/nlf/dlite)

```bash
docker-compose up --build

# Open a new terminal
docker-compose scale worker=3
docker-compose run worker ./kbmb start mit.edu
```

This is a distributed crawler that is the spiritual successor to a hacky system I setup at DrinkIn to crawl liquor prices from distributors.

Not being designed to work well on a single node, it is rather complicated to setup. I've simplified the configuration to make this a bit easier to set up.

It relies on Redis for the deduplication of URLs and NSQ as a queue. Each worker subscribes to the same channel and topic to fetch and process jobs as they come in.

# Deployment topology
- Setup a nsqlookupd cluster.
- Each worker box should run a local nsqd and worker process. A worker should transmit results to its local nsqd and consume from the nsqlookupd box(es).
- A job is started by sending a message to any single worker box.

# Terms
- root host: Initially started with mit.edu, which redirects to web.mit.edu, the root host is web.mit.edu. 

# Design decisions
- Follow the first redirect to determine the host. (e.g. mit.edu -> web.mit.edu once, every page after must have the web.mit.edu host).
- Don't parse pages that are accessible via the root host, but have a canonical URL that doesn't match (e.g. given a root host of example.com, visiting example.com/news redirects to news.example.com, so it is ignored).

# Goals
- Easy to set up a network of workers.
- Simple CLI client kick starts a job and aggregates results.

# TODO
- Currently starting a job starves the nodes, such that starting a new job will not begin until the FIFO queue of previous requests is processed. It's actually quite complicated to handle this in conjunction with rate limiting, and it's something I never needed,
- Distributed rate limiting, I'm a bad person and only rate limit per worker node.
- It does not clean up after itself on redis; I ran an ephemeral redis instance for this that was easy to restart.
