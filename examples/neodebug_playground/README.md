# Neodebug editor playground

This is not really a useful application example. It is a fixture for testing and trying out Neodebug's request editing functionality.

Run from this directory to open Neodebug:

```sh
go run .
```

An HTTP transporter is created only to describe the routes. It is never started or mounted on a server. The configured URL (`http://127.0.0.1:6122/http`) is a placeholder; sending requests will fail unless a separate server is listening there.

Choose a route to try its editors:

- `scalars`: int32, int64, float32, float64, bool, byte, string, date (`time.Time`), and file (`[]byte`). Enter a local file path in a file field to load its bytes; you can use this README.
- `collections`: slices, nested slices, maps with string and integer keys, nested/reused structs, file lists, and date lists.
- `nullable`: nullable primitives, dates, files, structs, slices, and nullable values inside collections.
- `all-types`: all the above together in one request.

These cover the types currently supported by Neodebug's editors. Arbitrary interface/serializable values are omitted because Neodebug has no editor for them.
