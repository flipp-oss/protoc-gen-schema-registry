# protoc-gen-schema-registry

Register local protobuf files with the Confluent Schema Registry. If using `buf`, you will need
to run it with `--include-imports`, and it will register all dependencies as well.

## Why a Protobuf Plugin?

This kind of functionality is usually better suited to a separate script - after all, this doesn't "generate" anything locally at all! However, the convenience of using `buf` to manage dependencies means we need to hook into `protoc` to get all the information we need for all dependent protobuf files.

## Usage

### Install using Golang 1.21+

Simply run the following command:

```bash
go install github.com/flipp-oss/protoc-gen-schema-registry@latest
```

### Install Manually

Download this project from the Releases page. Put the generated binary in your
path:

```bash
mv protoc-gen-schema-registry /usr/local/bin
```

Sync with the schema registry:

```bash
protoc --schema-registry_out=. --schema-registry_opt=registry_url=http://localhost:8081 *.proto
```

Using buf, in `buf.gen.yaml`:

```yaml
version: v2
managed:
  enabled: true
plugins:
  - local: protoc-gen-schema-registry
    opt:
      - registry_url=http://127.0.0.1:8081
```

then run: `buf generate --include-imports`

Typically you'd create a separate buf template just for this purpose, e.g. `buf.gen.schema-registry.yaml`, and run it with `buf generate --template buf.gen.schema-registry.yaml --include-imports`.

## Options

- `schema_registry_url` - The URL to reach the schema registry.

---

To Do List:

- Specify which dependencies to register / not to register
- Add tests
- Homebrew?

---

This project is supported by [Flipp](https://corp.flipp.com/).
