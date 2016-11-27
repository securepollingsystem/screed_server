## Screed Server

### Build

    $ go build

### Run

    $ ./screed_server

### Save New Screed

    $ cd $GOPATH/src/github.com/securepollingsystem/tallyspider/
    $ curl -i -X POST --data-binary @example_screed.txt localhost:8000/screed

This will return the key by which this screed can be retrieved.

### Retrieve Screed

    $ curl -i localhost:8000/screed/02759eaefd6359ce854d987611849103aab23c5d88b02c9ec86eed34ac833ceccd

This returns the same data as is in `example_screed.txt` above.
