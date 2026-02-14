# Packetbeat 8.x document structure and connection handling explained

**Packetbeat creates multiple documents per connection by design**—this is expected behavior, not a bug. The agent generates two distinct document types: **flow documents** (network-layer statistics reported periodically) and **transaction documents** (application-layer events per request/response). With default settings, a single TCP connection can produce 7+ flow documents and additional transaction documents, all correlatable via the `network.community_id` field.

## Two document creation mechanisms drive this behavior

Packetbeat 8.x uses two independent document generation systems that operate simultaneously:

**Flow-based documents** capture transport-layer (Layer 4) statistics. A flow is defined as "a group of packets sent over the same time period that share common properties"—source/destination addresses, ports, and protocol. These documents contain **bidirectional statistics** in a single record: `source.bytes`, `source.packets` alongside `destination.bytes`, `destination.packets`. The critical insight is that flows are reported **periodically**, not just once per connection.

**Transaction-based documents** capture application-layer (Layer 7) events. Each HTTP request/response pair, DNS query/answer, or database command creates a separate document containing protocol-specific decoded data. A persistent HTTP connection carrying 10 requests generates **10 separate transaction documents**, plus the flow documents tracking the underlying TCP connection.

Both mechanisms can generate documents for the same connection simultaneously. A 60-second HTTP session with 5 requests might produce: 6 intermediate flow reports + 1 final flow report + 5 HTTP transaction documents = **12 total documents** for one logical connection.

## Flow configuration directly controls document volume

The `flows.period` and `flows.timeout` settings determine how many flow documents Packetbeat creates:

| Setting | Default | Purpose |
|---------|---------|---------|
| `flows.enabled` | `true` | Enable/disable flow collection entirely |
| `flows.period` | `10s` | Interval for intermediate reports; set to `-1s` to disable |
| `flows.timeout` | `30s` | Flow killed after this inactivity duration |

**Document generation math**: For a flow lasting T seconds with default settings (period=10s, timeout=30s):
- Intermediate reports = `floor(T / 10)` documents with `flow.final: false`
- Final report = 1 document with `flow.final: true` when flow times out
- Total = intermediate + final

A 65-second connection produces intermediate reports at 10s, 20s, 30s, 40s, 50s, 60s (**6 documents**) plus one final report at ~95s (65s + 30s timeout) = **7 documents** for a single TCP flow.

The `flow.final` boolean field distinguishes document types:

```json
{
  "flow": {
    "final": false,
    "id": "FQQA/wz/Dv//////Fv8BAQEBAgMEBQYGBQQDAgGrAMsAcQPGM2QC9ZdQAA"
  }
}
```

**To reduce document volume**, set `period: -1s` to disable intermediate reports—you'll get only one final document per flow when it times out.

## Bidirectional traffic appears in single documents

Packetbeat does **not** create separate documents for each direction of a flow. Instead, each flow document contains bidirectional statistics:

```json
{
  "source": {
    "ip": "203.0.113.3",
    "port": 38901,
    "bytes": 1024,
    "packets": 10,
    "mac": "01-02-03-04-05-06"
  },
  "destination": {
    "ip": "198.51.100.2",
    "port": 80,
    "bytes": 45000,
    "packets": 35,
    "mac": "06-05-04-03-02-01"
  },
  "network": {
    "bytes": 46024,
    "packets": 45,
    "transport": "tcp",
    "type": "ipv4",
    "community_id": "1:t9T66/2c66NQyftAEsr4aMZv4Hc="
  }
}
```

The `source.*` fields represent the flow originator (typically the client), while `destination.*` fields represent the responder (typically the server). The `network.bytes` and `network.packets` fields provide **aggregated totals** for both directions combined.

If you're seeing what appears to be "individual packet" source/destination IPs, you're likely looking at **transaction documents** which capture per-request details, or intermediate flow reports which still show the connection endpoints but with cumulative statistics at that point in time.

## Community ID enables reliable correlation across documents

The `network.community_id` field is the **primary correlation mechanism** for grouping all documents belonging to the same connection. It implements the Community ID v1 standard—an open specification producing identical hashes regardless of traffic direction.

The algorithm orders endpoints deterministically (smaller IP:port first) before hashing, ensuring both directions produce the same identifier:

```
Client → Server: TCP 10.0.0.1:1234 → 192.168.1.1:80
Server → Client: TCP 192.168.1.1:80 → 10.0.0.1:1234

Both produce: 1:LQU9qZlK+B5F3KDmev6m5PMibrg=
```

**Key distinction between correlation fields**:

| Field | Scope | Use case |
|-------|-------|----------|
| `network.community_id` | Cross-tool standard | Correlate with Zeek, Suricata, Wireshark |
| `flow.id` | Packetbeat internal | More specific, includes MAC/VLAN layers |

Use `network.community_id` for all correlation queries—it's consistent across flow and transaction documents for the same 5-tuple (source IP, source port, destination IP, destination port, protocol).

## Essential ECS fields for connection analysis

Packetbeat populates these Elastic Common Schema fields for network analysis:

**Flow identification**:
- `flow.id` — Internal identifier including MAC/IP/port layers
- `flow.final` — Boolean indicating final vs intermediate report
- `flow.vlan` — VLAN tag from 802.1q frames

**Network metadata**:
- `network.community_id` — Standardized flow hash
- `network.direction` — Traffic direction: `ingress`, `egress`, `inbound`, `outbound`, `internal`, `external`
- `network.transport` — Protocol: `tcp`, `udp`, `icmp`
- `network.type` — Layer 3 type: `ipv4`, `ipv6`
- `network.protocol` — Application protocol: `http`, `dns`, etc.
- `network.bytes` / `network.packets` — Bidirectional totals

**Endpoint fields** (directional—may swap between packets):
- `source.ip`, `source.port`, `source.bytes`, `source.packets`, `source.mac`
- `destination.ip`, `destination.port`, `destination.bytes`, `destination.packets`, `destination.mac`

**Role-based fields** (consistent throughout connection):
- `client.ip`, `client.port` — Connection initiator (TCP SYN sender)
- `server.ip`, `server.port` — Connection responder

**Event metadata**:
- `event.dataset` — Document type: `flow`, `http`, `dns`, etc.
- `event.start` / `event.end` — Flow timestamps
- `event.duration` — Duration in nanoseconds

## Query patterns for accurate connection analysis

**Find all events for a specific connection**:

```json
{
  "query": {
    "term": {
      "network.community_id": "1:t9T66/2c66NQyftAEsr4aMZv4Hc="
    }
  },
  "sort": [{"@timestamp": "asc"}]
}
```

**Aggregate bytes per connection (avoiding double-counting)**:

The critical technique is **filtering on `flow.final: true`** to exclude intermediate reports:

```json
{
  "query": {
    "bool": {
      "must": [
        {"term": {"event.dataset": "flow"}},
        {"term": {"flow.final": true}}
      ]
    }
  },
  "aggs": {
    "connections": {
      "terms": {"field": "network.community_id", "size": 100},
      "aggs": {
        "total_bytes": {"sum": {"field": "network.bytes"}},
        "source_bytes": {"sum": {"field": "source.bytes"}},
        "dest_bytes": {"sum": {"field": "destination.bytes"}}
      }
    }
  }
}
```

**Correlate flows with application transactions**:

```json
{
  "query": {
    "term": {"network.community_id": "1:t9T66/2c66NQyftAEsr4aMZv4Hc="}
  },
  "aggs": {
    "by_type": {
      "terms": {"field": "event.dataset"},
      "aggs": {
        "count": {"value_count": {"field": "_id"}}
      }
    }
  }
}
```

## Configuration for minimal document generation

To reduce document volume while maintaining visibility:

```yaml
packetbeat.flows:
  enabled: true
  timeout: 30s
  period: -1s  # Disable intermediate reports - only final reports

packetbeat.protocols:
  - type: http
    ports: [80, 8080]
    transaction_timeout: 10s
```

This configuration produces exactly **one flow document per connection** (when it times out), plus one transaction document per application-layer request/response pair.

For high-resolution monitoring where you need traffic trends over time:

```yaml
packetbeat.flows:
  enabled: true
  timeout: 60s
  period: 5s
  keep_counters_on_report: true  # Delta values instead of cumulative
```

The `keep_counters_on_report: true` option changes `network.bytes` from cumulative totals to **delta values** (increments since last report), simplifying rate calculations.

## Conclusion

The multiple documents per connection in Packetbeat 8.x result from periodic flow reporting (configurable via `flows.period`) combined with per-transaction application-layer events. This design supports both real-time traffic monitoring and detailed application analysis. For connection-level analysis, always filter on `flow.final: true` to avoid counting intermediate reports, and use `network.community_id` as your primary grouping field—it provides consistent bidirectional correlation across all document types and is compatible with other network monitoring tools like Zeek and Suricata.