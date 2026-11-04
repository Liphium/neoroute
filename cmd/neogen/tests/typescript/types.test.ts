// oxlint-disable no-unused-vars
// tsc compiles this file, but Node never runs it. A failed assertion fails npm test.
// Unused aliases are intentional: asking TypeScript to check them is the test.
import type {
  Address,
  AddressWire,
  Inner,
  InnerWire,
  User,
  UserWire,
  Standalone,
  StandaloneWire,
  StandaloneInner,
  StandaloneInnerWire,
  Coordinates,
  CoordinatesWire,
  Location,
  LocationWire,
  Point,
  PointWire,
  Envelope,
  EnvelopeWire,
} from "./generated/models.js";
import {
  toPointWire,
  fromPointWire,
  toEnvelopeWire,
  fromEnvelopeWire,
} from "./generated/models.js";
import type { HttpConnector } from "./generated/connector_http.js";
import type { WsConnector } from "./generated/connector_ws.js";
import type { RequestOptions, UserError } from "@liphium/neoroute-ts";

// Equal compares exact shapes; Assert rejects the result unless it is true.
// This catches extra/missing fields as well as incorrect field types.
type Assert<T extends true> = T;
type Equal<A, B> =
  (<T>() => T extends A ? 1 : 2) extends <T>() => T extends B ? 1 : 2
    ? true
    : false;

// Application models always have named fields, even for Go tuple structs.
// These expectations also cover standalone models and omitted unexported Go fields.
type ModelAssertions = [
  Assert<Equal<Address, { city: string; zip: bigint | null }>>,
  Assert<Equal<Inner, { note: string }>>,
  Assert<
    Equal<
      User,
      {
        name: string;
        addr: Address | null;
        tags: string[];
        inner: Inner;
        friends: User[];
      }
    >
  >,
  Assert<Equal<StandaloneInner, { info: string }>>,
  Assert<
    Equal<
      Standalone,
      {
        code: bigint;
        label: string | null;
        items: StandaloneInner[];
        some_stamp: Date;
      }
    >
  >,
  Assert<Equal<Coordinates, { x: number; y: number }>>,
  Assert<Equal<Location, { name: string; coordinates: Coordinates }>>,
  Assert<Equal<Point, { x: number; label: string }>>,
  Assert<
    Equal<
      Envelope,
      {
        point: Point;
        items: (Point | null)[];
        matrix: Point[][];
        lookup: Record<string, Point | null>;
        numeric: Record<number, Point>;
        next: Envelope | null;
        bytes: Uint8Array;
        date: Date;
        opaque: any;
      }
    >
  >,
];

// Wire shapes follow Go encoding: tuple structs become ordered arrays.
// Maps stay records and arrays stay arrays, but their nested model types become wire types.
type WireAssertions = [
  Assert<Equal<AddressWire, { city: string; zip: bigint | null }>>,
  Assert<Equal<InnerWire, { note: string }>>,
  Assert<
    Equal<
      UserWire,
      {
        name: string;
        addr: AddressWire | null;
        tags: string[];
        inner: InnerWire;
        friends: UserWire[];
      }
    >
  >,
  Assert<Equal<StandaloneInnerWire, { info: string }>>,
  Assert<
    Equal<
      StandaloneWire,
      {
        code: bigint;
        label: string | null;
        items: StandaloneInnerWire[];
        some_stamp: Date;
      }
    >
  >,
  Assert<Equal<CoordinatesWire, [number, number]>>,
  Assert<Equal<LocationWire, { name: string; coordinates: CoordinatesWire }>>,
  Assert<Equal<PointWire, [number, string]>>,
  Assert<
    Equal<
      EnvelopeWire,
      {
        point: PointWire;
        items: (PointWire | null)[];
        matrix: PointWire[][];
        lookup: Record<string, PointWire | null>;
        numeric: Record<number, PointWire>;
        next: EnvelopeWire | null;
        bytes: Uint8Array;
        date: Date;
        opaque: any;
      }
    >
  >,
];

// Converters must accept and return the correct side of the application/wire boundary.
type ConverterAssertions = [
  Assert<Equal<typeof toPointWire, (value: Point) => PointWire>>,
  Assert<Equal<typeof fromPointWire, (value: PointWire) => Point>>,
  Assert<Equal<typeof toEnvelopeWire, (value: Envelope) => EnvelopeWire>>,
  Assert<Equal<typeof fromEnvelopeWire, (value: EnvelopeWire) => Envelope>>,
];

// Callers of either connector should only see application models and request options.
// Include every send variant plus routes whose roots are arrays or maps.
type ExpectedMethods = {
  getUser: (
    payload: User,
    options?: RequestOptions,
  ) => Promise<User | UserError>;
  exchange: (
    payload: Envelope,
    options?: RequestOptions,
  ) => Promise<Envelope | UserError>;
  get: (options?: RequestOptions) => Promise<Envelope | UserError>;
  ok: (
    payload: Envelope,
    options?: RequestOptions,
  ) => Promise<UserError | void>;
  push: (
    payload: Envelope,
    options?: RequestOptions,
  ) => Promise<UserError | void>;
  signal: (options?: RequestOptions) => Promise<UserError | void>;
  ping: (options?: RequestOptions) => Promise<UserError | void>;
  collection: (
    payload: (Point | null)[],
    options?: RequestOptions,
  ) => Promise<(Point | null)[] | UserError>;
  lookup: (
    payload: Record<string, Point | null>,
    options?: RequestOptions,
  ) => Promise<Record<string, Point | null> | UserError>;
};

// Pick checks the route methods without comparing private fields or transport methods.
// The event callback must also expose the decoded application model.
type ConnectorAssertions = [
  Assert<Equal<Pick<HttpConnector, keyof ExpectedMethods>, ExpectedMethods>>,
  Assert<Equal<Pick<WsConnector, keyof ExpectedMethods>, ExpectedMethods>>,
  Assert<
    Equal<
      HttpConnector["onUpdate"],
      (handler: (event: Envelope) => void) => void
    >
  >,
];

// Each @ts-expect-error requires the next line to fail compilation.
// If invalid usage ever becomes accepted, tsc reports an unused directive.
// Keeping these expressions in an uncalled function prevents accidental execution.
function invalidUsage(client: HttpConnector, point: Point, wire: PointWire) {
  // @ts-expect-error application models are objects
  const invalidModel: Point = [4, "test"];
  // @ts-expect-error tuples retain field order and field types
  const invalidWire: PointWire = ["test", 4];
  // @ts-expect-error converters take objects, not wire tuples
  toPointWire(wire);
  // @ts-expect-error decoders take wire tuples, not objects
  fromPointWire(point);
  const invalidNested: LocationWire = {
    name: "test",
    // @ts-expect-error nested tuple values must be converted on the wire
    coordinates: { x: 1, y: 2 },
  };
  // @ts-expect-error standalone models must not include unexported Go fields
  const privateField: keyof Standalone = "secret";
  // @ts-expect-error route payloads use application objects
  client.collection([wire]);
}

// Every generated model must expose both converter directions.
import * as models from "./generated/models.js";

// Each pair is [application type, wire type]; keep this list in sync with the fixtures.
type ModelPairs = {
  Address: [Address, AddressWire];
  Inner: [Inner, InnerWire];
  User: [User, UserWire];
  Standalone: [Standalone, StandaloneWire];
  StandaloneInner: [StandaloneInner, StandaloneInnerWire];
  Coordinates: [Coordinates, CoordinatesWire];
  Location: [Location, LocationWire];
  Point: [Point, PointWire];
  Envelope: [Envelope, EnvelopeWire];
};

// Derive the expected exported function names and signatures for every model pair.
type Encoders = {
  [K in keyof ModelPairs as `to${K}Wire`]: (
    value: ModelPairs[K][0],
  ) => ModelPairs[K][1];
};
type Decoders = {
  [K in keyof ModelPairs as `from${K}Wire`]: (
    value: ModelPairs[K][1],
  ) => ModelPairs[K][0];
};
type AllConverterAssertions = [
  Assert<Equal<Pick<typeof models, keyof Encoders>, Encoders>>,
  Assert<Equal<Pick<typeof models, keyof Decoders>, Decoders>>,
];
