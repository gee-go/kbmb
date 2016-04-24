# Running
- Install docker and docker-compose (on a mac I like https://github.com/nlf/dlite)

```bash
docker-compose up --build

# Open a new terminal
docker-compose scale worker=3
docker-compose run worker ./kbmb start mit.edu
```