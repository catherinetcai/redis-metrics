### redis-metrics

This is a small tool written using [go-metrics](https://github.com/rcrowley/go-metrics) and [rdbtools](https://github.com/cupcake/rdb). It has a quick regexp keymatch for redis to pull out keys and do a key count, memory count, and gets an average/std dev/variance.

#### In the werkz:
- Pull in a yaml file to specify what actual keys you want to match against
- Emit metrics to a visualizer (or visualize using d3)
