# neogen

Neogen is the code generator for the Neoroute framework. It generates code that can be used with official client SDKs.

Documentation for it is available at [liphium.dev/neoroute/utility/neogen](https://liphium.dev/neoroute/utility/neogen).

### TypeScript models and wire conversion

Generated TypeScript models expose named object fields, including Go structs marked with `neoschema.IsTuple`. Each model `Foo` also has a `FooWire` type and `toFooWire` / `fromFooWire` functions in `models.ts`. These convert recursively through nested models, arrays, record values, and nullable fields.

The behavior of these functions matches what `msgp` generates for Go structs. Tuples are basically arrays of the values in the object, and when it's not a tuple you get a map.

The generated `to` and `from` functions are used in all places to make sure the objects are correctly converted to their wire equivalents.

### What this means for you

When calling `Client` directly, use wire types (and their functions) instead of the normal objects:

```ts
const response = await client.send<UserWire, UserWire>("updateUser", toUserWire(user));
if (!(response instanceof UserError)) {
  const updatedUser = fromUserWire(response);
}
```

If you use tuples and reorder fields, you need to re-generate the code to make sure the wire conversion functions are correct.

## Testing generated code

See [the runtime test projects](tests/README.md) for the different tests we run to verify that neogen is actually generating validated code. Most of this has been created with the assistance of AI and just verified by eye since it's too much work to write by hand for a small project like ours.
