# Watch a log and set to statsd!

This is a very simple tool that will watch a log and increase a StatsD metric when a line matches the specified regex.

It was initially created to watch a MySQL slow query log and plot it on a dashboard, hence the name. Now it's a more generic tool.

Just check the `--help` to see how to use it. Remember to enclose the regex with single quotes. `''`