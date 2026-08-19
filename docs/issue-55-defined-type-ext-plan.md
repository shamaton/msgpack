# Issue #55: Extension Encoders for Non-Struct Types

## Status

Proposed implementation plan for [Issue #55](https://github.com/shamaton/msgpack/issues/55).

## Problem

Extension encoders are currently selected only from the `reflect.Struct` encoding
path. An encoder registered for another exact Go type is therefore ignored when
that type has an `int`, `string`, `slice`, or other underlying kind.

Examples that must be supported:

```go
type Status int
type Statuses []Status
```

1. A registered encoder for `Status` must be used for a root value, struct field,
   map key or value, interface value, and each element of `[]Status`.
2. A registered encoder for `Statuses` must encode the whole value as one
   extension value.
3. If encoders are registered for both `Statuses` and `Status`, the exact
   `Statuses` encoder wins. Its elements are not encoded separately.

The same rules apply to the byte-slice encoder and the stream encoder.

## Scope

This issue changes extension **encoder selection** only. Decoder routing for
non-struct destination types is not part of Issue #55. Consequently, this work
must not claim round-trip support for defined types until decoder support is
implemented separately.

No public interface in `ext` or `msgpack` is changed.

## Required Semantics

- Selection is by exact `reflect.Type`, not only by `reflect.Kind`.
- A registered exact type takes precedence over every built-in representation,
  including bin, fixed-slice, map, struct, pointer, and interface handling.
- This precedence also applies to nil values of registered nil-capable types;
  the extension encoder is responsible for handling the value it registered.
- An unregistered defined type keeps its existing underlying-kind encoding.
- A type alias is not distinct at runtime. Registering an alias has the same
  effect as registering its aliased type.
- `AddExtCoder` and `RemoveExtCoder` continue to ignore attempts to replace or
  remove the built-in `time.Time` encoder.
- One encode operation uses one immutable registry snapshot for its entire
  lifetime. Registration changes made concurrently apply only to later encode
  operations. This is required so `CalcByteSize` and `WriteToBytes` cannot select
  different encoders.

## Registry Design

Replace the mutable map plus derived linear slice in both encoding packages with
an immutable registry indexed first by kind and then by exact type:

```go
type extEncoderRegistry struct {
    byKind [reflect.UnsafePointer + 1]map[reflect.Type]ext.Encoder
    hasCustom [reflect.UnsafePointer + 1]bool
    customCount int
}
```

The stream package uses the equivalent type with `ext.StreamEncoder`.

- Publish the current registry through `atomic.Pointer[extEncoderRegistry]`.
- Serialize `AddExtEncoder` and `RemoveExtEncoder` with a package-level mutex.
- Each mutation performs copy-on-write for the affected kind map and publishes
  the new immutable snapshot.
- Set a kind bucket back to nil when its last encoder is removed; the nil-bucket
  fast path must be restored rather than retaining an empty map.
- Track the total custom encoder count and per-kind custom presence in the same
  immutable snapshot. When the total is zero, select the built-in-only dispatcher
  once at Encode entry so ordinary values do not perform per-value registry
  checks. The built-in-only dispatcher handles `time.Time` by direct type
  comparison.
- Initialize the struct-kind map with the built-in `time.Time` encoder.
- Store the loaded snapshot on the per-operation `encoder` value.
- Never mutate a map after its registry snapshot has been published.

Lookup must first read the kind bucket. If it is nil, return immediately without
a type-map lookup:

```go
func (e *encoder) extEncoderFor(t reflect.Type) (ext.Encoder, bool) {
    coders := e.extRegistry.byKind[t.Kind()]
    if coders == nil {
        return nil, false
    }
    coder, ok := coders[t]
    return coder, ok
}
```

With only the built-in encoder registered, non-struct values pay one predictable
nil-bucket branch and no map lookup. Struct lookup changes from the current
linear scan to an exact map lookup.

## Byte-Slice Encoder Changes

Split dynamic dispatch from built-in encoding:

- `calcSize` checks validity, then exact extension registration, then delegates
  to `calcSizeDefault` for the existing kind switch.
- `create` checks validity, then exact extension registration, then delegates to
  `createDefault` for the existing kind switch.
- Remove the root-only pointer unwrapping from `Encode`. Pointer and interface
  traversal belongs in the common dispatcher so exact pointer registrations,
  nil pointers, and arbitrary pointer depth behave consistently at the root and
  when nested.
- Keep invalid-value encoding as MessagePack nil; do not call `Type()` on an
  invalid `reflect.Value`.
- Remove extension scans from `calcStruct`, `writeStruct`, `getStructCalc`, and
  `getStructWriter`; all extension selection belongs to the common resolver.
- Replace any remaining linear extension scan with exact registry lookup.

For homogeneous slices and arrays, resolve the element functions once before the
loop. Do not call the extension-aware top-level dispatcher for every element:

```go
calcElement := e.resolveCalcFunc(rv.Type().Elem())
writeElement := e.resolveWriteFunc(rv.Type().Elem())
```

The resolver returns the registered extension function when present. Otherwise
it returns a built-in function that does not repeat the exact-type lookup for
that same value. Nested pointer, interface, slice, array, map, and struct values
may re-enter dynamic dispatch where their runtime or child type requires it.

Container selection order is mandatory:

1. Resolve an extension for the exact whole-container type. Use it immediately
   when present.
2. For a slice or array without a whole-container encoder, resolve its element
   type once from the operation's registry snapshot.
3. Use bin or fixed-slice handling only when the element type has no registered
   extension. Otherwise, use array encoding and the resolved element functions.
4. For a map without a whole-container encoder, resolve its declared key and
   value types once from the same snapshot.
5. Use fixed-map handling only when neither the key type nor value type has a
   registered extension. Otherwise, use the generic map path with the resolved
   key and value functions.

This means registering an encoder for `byte` intentionally changes `[]byte` and
`[N]byte` from bin encoding to an array of extension values. Registering an
encoder for `int` or `string` similarly bypasses matching fixed-slice and
fixed-map paths. Preserve all existing specialized paths when their child types
have no registered extension. The byte encoder's calculation and write passes
must use the same snapshot, ordering, and bypass conditions.

## Stream Encoder Changes

Apply the same registry snapshot, precedence, and resolver rules under
`internal/stream/encoding`.

- `create` performs exact extension selection before its built-in kind switch.
- Remove the root-only pointer unwrapping from stream `Encode` for the same
  reason as the byte-slice encoder.
- A selected encoder receives one `ext.StreamWriter` and writes the complete
  value.
- Slice and array element writers, and map key and value writers, are resolved
  once before iteration.
- Bin, fixed-slice, and fixed-map paths use the same child-extension bypass rules
  as the byte-slice encoder.
- Remove the struct-only extension scans.
- The existing buffer acquisition and flush lifecycle remains unchanged.

The byte-slice and stream implementations must have matching behavior tests.

## Implementation Order

1. Add baseline benchmarks and record results before functional changes.
2. Introduce the immutable kind-indexed registry in `internal/encoding` and
   update its registration tests.
3. Add common exact-type dispatch and the split default functions.
4. Resolve slice and array element functions once per container.
5. Mirror steps 2-4 in `internal/stream/encoding`.
6. Add behavioral, precedence, removal, and concurrency tests.
7. Run correctness, race, and benchmark comparison gates.

Do not introduce a global general-purpose type-plan cache in this change. It
would require generation-aware invalidation for struct mode and registration
changes. Add it only in a later optimization if benchmarks show the kind-indexed
resolver is insufficient.

## Tests

Cover both `MarshalAsMap`/`MarshalAsArray` and
`MarshalWriteAsMap`/`MarshalWriteAsArray` where applicable:

- root `Status` uses its extension encoder;
- a `Status` struct field uses its extension encoder;
- map keys and values with type `Status` use their extension encoder;
- a `Status` stored in `any` uses its extension encoder;
- each element of `[]Status` uses the `Status` encoder;
- `Statuses` uses its whole-slice encoder;
- the `Statuses` encoder wins when both whole-slice and element encoders exist;
- a registered defined `[]byte` type wins over bin encoding;
- registering `byte` makes `[]byte` and `[N]byte` encode as containers of
  extension values instead of bin;
- registering a primitive `int` or `string` bypasses every matching fixed-slice
  and fixed-map path and applies to each key or value;
- nil of a registered slice or pointer type is passed to its encoder;
- unregistered defined types retain their current encoding;
- removing a coder restores built-in encoding;
- adding the same exact type twice preserves the first registration, matching
  current behavior;
- holding an old registry snapshot across add and remove operations shows that
  its registry object and kind maps remain unchanged;
- removing the last encoder for a non-time kind leaves that kind bucket nil in
  the newly published snapshot;
- `time.Time` behavior remains unchanged;
- concurrent registration and encoding is race-free, and a byte-slice encode
  cannot mix size and write selections from different snapshots;
- deterministically verify the byte-slice snapshot by blocking the first
  encoder's `CalcByteSize`, removing it and adding a replacement for the same
  exact type while blocked, and then resuming. The in-progress encode must use
  the first encoder for both calculation and writing; only the next encode may
  use the replacement;
- deterministically verify the stream snapshot with
  `[]any{Status(1), Status(2)}` so each element requires runtime-type dynamic
  dispatch. Block the first encoder in the first element's `Write`, remove it
  and add a replacement for `Status` while blocked, and then resume. The second
  element must still resolve through the operation's original snapshot and use
  the first encoder; only the next encode may use the replacement.

## Performance Gates

Add benchmarks for both byte-slice and stream APIs with these fixed fixtures:

- scalar `int(42)`;
- a 1,024-element `[]int`, populated outside the timed loop with `i % 128`;
- `benchmarkStruct{ID: 42, Name: "benchmark", Active: true, Score: 12.5,
  Tags: []string{"a", "b", "c", "d"}, Data: make([]byte, 64)}`;
- `time.Unix(1_700_000_000, 123_456_789)`;
- a no-extension struct containing a pointer, an interface field, `[]any`, and
  a nested `[][]int` container;
- registered `Status(42)`;
- a 1,024-element `[]Status`, populated with `Status(i % 128)`, using the
  element encoder;
- the same 1,024 values converted to registered whole-value `Statuses`.

Define `benchmarkStruct`, `Status`, `Statuses`, all populated inputs, and ext
encoder instances at package scope or before `ResetTimer`. Register extension
encoders before the unmeasured warm-up call and remove them after `StopTimer`.
Every benchmark must call `ReportAllocs`, perform one unmeasured encode, and then
call `ResetTimer`.

For byte-slice benchmarks, assign the result to a package-level byte slice and
check the error in the loop. For stream benchmarks, allocate one `bytes.Buffer`
outside the timed loop, call `Reset` inside the loop immediately before encode,
and assign its final length to a package-level integer. Do not allocate a new
writer or input fixture per iteration. The unmeasured call warms the internal
buffer pool before timing.

Record baseline before functional changes and candidate results afterward on the
same otherwise-idle machine, with the same Go toolchain, environment, and
`GOMAXPROCS=1`. Collect 20 samples for each and compare with `benchstat` at its
95% confidence default. A timing gate fails only when the result is statistically
significant (`p < 0.05`) and the magnitude exceeds the applicable threshold.
Acceptance criteria:

- no additional allocations for any non-extension benchmark;
- no statistically significant regression greater than 2% for representative
  non-extension containers or structs;
- no regression greater than 3% for the scalar benchmark, where one additional
  branch is proportionally more visible;
- `[]Status` performs extension lookup once per container type resolution, not
  once per element. Confirm this from the resolver placement and benchmark
  scaling; calls into the selected extension encoder still occur per element.

Suggested commands:

```sh
go test ./...
go test -race ./...
GOMAXPROCS=1 go test ./internal/encoding ./internal/stream/encoding \
  -run '^$' -bench '^Benchmark(Encode|StreamEncode)' \
  -benchmem -count=20 > /tmp/issue55-baseline.txt
# Run the identical command after implementation, writing candidate.txt.
benchstat /tmp/issue55-baseline.txt /tmp/issue55-candidate.txt
```

## Verification Results

Verified on 2026-07-17 with Go 1.26.1 on Apple M4 Pro. Both benchmark files
contain 20 samples per case and were compared with `benchstat`:

- `go test ./...`: passed;
- `go test -race ./...`: passed;
- byte `EncodeIntSlice`: 2.254 us/op to 2.191 us/op (`-2.80%`), and allocations
  decreased from 2 to 1;
- byte `EncodeStruct`: 187.3 ns/op to 170.1 ns/op (`-9.16%`), and allocations
  decreased from 2 to 1;
- byte `EncodeIndirect`: 355.3 ns/op to 353.9 ns/op, with no statistically
  significant change (`p=0.151`), and allocations decreased from 6 to 5;
- stream `StreamEncodeIntSlice`: 3.542 us/op to 3.514 us/op (`-0.76%`), and
  allocations decreased from 1 to 0;
- stream `StreamEncodeStruct`: 176.0 ns/op to 162.6 ns/op (`-7.64%`), and
  allocations decreased from 1 to 0;
- stream `StreamEncodeIndirect`: 247.1 ns/op to 234.2 ns/op (`-5.22%`), and
  allocations decreased from 3 to 2;
- no non-extension benchmark gained allocations;
- byte `EncodeStatusSlice` uses 2 allocations per operation. Comparisons for
  extension benchmarks include the intended behavior change because the
  baseline did not apply non-struct extension encoders.

The performance gates pass. The complete raw results are stored in
`/tmp/issue55-baseline-final.txt` and `/tmp/issue55-candidate-final.txt` for this
worktree session. Each file contains 320 benchmark samples: 16 cases with 20
samples per case.

## Completion Criteria

The implementation is complete when all required semantics are covered in both
encoding backends, the full and race test suites pass, performance gates pass,
and the public documentation does not imply decoder round-trip support.
