// Run the compiled generated code with a fake SDK: no server or network is involved.
// The separate types.test.ts checks types; this file checks actual values and calls.
const assert = require("node:assert/strict");
const { toEnvelopeWire, fromEnvelopeWire } = require("./dist/generated/models.js");
const { HttpConnector } = require("./dist/generated/connector_http.js");
const { UserError } = require("@liphium/neoroute-ts");
require("node:test").test("recursive wire conversions and generated connectors", async () => {
  // Point is a tuple in Go, but application code should see named object fields.
  const point = { x: 4, label: "test" };

  // Exercise tuple values inside arrays, maps, nullable fields and recursive models.
  const model = {
    point,
    items: [point, null],
    matrix: [[point]],
    lookup: { a: point, b: null },
    numeric: { 1: point },
    next: null,
    bytes: new Uint8Array([1, 2]),
    date: new Date(1000),
    opaque: { whatever: [1, 2] },
  };

  // Only tuple structs become arrays; the surrounding containers keep their shape.
  const wire = toEnvelopeWire(model);
  assert.deepEqual(wire.point, [4, "test"]);
  assert.deepEqual(wire.items, [[4, "test"], null]);
  assert.deepEqual(wire.matrix, [[[4, "test"]]]);
  assert.deepEqual(wire.lookup, { a: [4, "test"], b: null });
  assert.deepEqual(wire.numeric, { 1: [4, "test"] });

  // Bytes, dates and opaque values pass through without copying or inspecting them.
  assert.equal(wire.bytes, model.bytes);
  assert.equal(wire.date, model.date);
  assert.equal(wire.opaque, model.opaque);

  // Decoding must restore the original objects, including nullable and nested fields.
  assert.deepEqual(fromEnvelopeWire(wire), model);

  // A recursive type can contain another model, provided the runtime values are acyclic.
  assert.deepEqual(fromEnvelopeWire(toEnvelopeWire({ ...model, next: model })), {
    ...model,
    next: model,
  });

  // Conversion must not replace the caller's original nested object.
  assert.equal(model.point, point);

  // The fake Client records outgoing payloads and returns its configurable result.
  // Accessing cl here deliberately bypasses TypeScript privacy for test setup.
  const errors = [];
  const connector = new HttpConnector(
    {
      errorHandler: (e) => {
        errors.push(e);
        return e.message;
      },
    },
    "POST",
    "/",
  );

  // Pretend the server returned wire data: the connector must decode it and encode requests.
  connector.cl.result = wire;
  assert.deepEqual(await connector.exchange(model), model);
  assert.deepEqual(connector.cl.payload, wire);
  assert.deepEqual(await connector.get(), model);

  // Route roots can themselves be arrays or maps, rather than named model objects.
  connector.cl.result = [[4, "test"], null];
  assert.deepEqual(await connector.collection([point, null]), [point, null]);
  assert.deepEqual(connector.cl.payload, [[4, "test"], null]);
  connector.cl.result = { a: [4, "test"], b: null };
  assert.deepEqual(await connector.lookup({ a: point, b: null }), {
    a: point,
    b: null,
  });
  assert.deepEqual(connector.cl.payload, { a: [4, "test"], b: null });

  // No-payload routes return normally; request-only routes still encode their payloads.
  assert.equal(await connector.signal(), undefined);
  assert.equal(await connector.ping(), undefined);
  await connector.ok(model);
  assert.deepEqual(connector.cl.payload, wire);
  await connector.push(model);
  assert.deepEqual(connector.cl.payload, wire);

  // Simulate an incoming event by calling the callback saved by fake Client.receive.
  // The application handler must receive decoded objects, not nested wire tuples.
  let event;
  connector.onUpdate((value) => (event = value));
  connector.cl.handler({}, wire);
  assert.deepEqual(event, model);

  // A server UserError must pass through unchanged without being decoded as a model.
  const error = new UserError("server rejected");
  connector.cl.result = error;
  assert.equal(await connector.exchange(model), error);

  // Malformed requests/responses fail conversion and use the configured error handler.
  assert.ok((await connector.exchange(null)) instanceof UserError);
  connector.cl.result = {};
  assert.ok((await connector.get()) instanceof UserError);

  // A fire-and-forget method reports conversion errors but returns no UserError.
  assert.equal(await connector.push(null), undefined);

  // One error each came from the invalid request, invalid response and invalid push.
  assert.equal(errors.length, 3);
});

require("node:test").test("custom and standalone models preserve their values", () => {
  const {
    toUserWire,
    fromUserWire,
    toStandaloneWire,
    fromStandaloneWire,
    toLocationWire,
    fromLocationWire,
  } = require("./dist/generated/models.js");

  // Regular object models also recurse through nullable references and child arrays.
  const child = {
    name: "child",
    addr: null,
    tags: [],
    inner: { note: "nested" },
    friends: [],
  };
  const user = {
    ...child,
    name: "parent",
    addr: { city: "Tokyo", zip: 123n },
    tags: ["one"],
    friends: [child],
  };
  assert.deepEqual(fromUserWire(toUserWire(user)), user);

  // Standalone models are registered explicitly, without any route or event using them.
  const standalone = {
    code: 42n,
    label: null,
    items: [{ info: "standalone child" }],
    some_stamp: new Date(1000),
  };
  assert.deepEqual(fromStandaloneWire(toStandaloneWire(standalone)), standalone);

  // An object containing a tuple remains an object; only its coordinates become an array.
  // Different x/y values also catch accidental tuple field reordering.
  const location = { name: "home", coordinates: { x: 1, y: 2 } };
  assert.deepEqual(toLocationWire(location), {
    name: "home",
    coordinates: [1, 2],
  });
  assert.deepEqual(fromLocationWire(toLocationWire(location)), location);
});
