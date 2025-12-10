# dnstom
DNS  functionality written in Go for education purposes

## Initial packages:

dnswire – encode/decode wire format

resolver – perform recursive resolution

auth – authoritative server behaviour

niosmodel – translate Infoblox objects to generic DNS views/RRs

## Try to make the entry point as dumb as possible

```go
func main() {
    cfg := loadConfig()
    if err := app.Run(cfg); err != nil {
        log.Fatal(err)
    }
}
```

## Initial layout plan

```tree
dnstom/
  go.mod

  cmd/
    dnstom-dig/
      main.go        # dig-like client (first tool)
    # later:
    # dnstom-auth/
    # dnstom-resolve/

  internal/
    dnswire/
      message.go     # DNS message structs + encode/decode + PrettyPrint
    rr/
      rr.go          # Basic RR types & helpers
    resolver/
      resolver.go    # Client that talks to upstream resolvers (stub/recursive, later)

```

### Appendix/Refernce

[https://www.isc.org/bind/]
[https://docs.infoblox.com/space/nios90/1380844355/DNS+Response+Message+Format+and+Examples]