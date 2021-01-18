# test-stream
Stream a zip file with random data for a a number of seconds with a max file size

# usage

`GET /download/:seconds/:size`

Where `:seconds` is how long the download should last in seconds and `:size` is
how large (MB) the zip file should be maximum. Based on these numbers the
program will calculate an average MB/sec speed and will throttle the stream
to be around this value.

```shell
wget https://location/download/:seconds/:size
```

# license
License is MIT
