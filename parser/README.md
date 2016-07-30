## Parser

The Parser works by to parse out incoming messages for the incoming network.
Work will be done to better build queries for remote nodes here later on. It's
definitely nothing special, as our queries are simple (and so shall remain!) so
when parsing an incoming message, we essentially can just split at commas.

An example command sent into the node looks like so:

```
GET key1,key2,key3
```

This command is essentially split into a data structure which looks like:

```
{
  "Command": "GET",
  "Args": {"key1": "", "key2": "", "key3": ""}
}
```

After the parser builds this out, it will return the command structure to
the incoming network which will then process the command and respond to the
calling node/client.

The incoming network will essentially respond like so:

```
GOT key1:value1, key2:value2, key3:value3
```

Assuming each key is available. If, however, key2 weren't available, the
incoming network would respond like so:

```
GOT key1:value1, key3:value3
```


