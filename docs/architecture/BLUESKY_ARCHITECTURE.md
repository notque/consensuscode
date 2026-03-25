# Bluesky Collective Architecture

The Bluesky Collective tool manages the collective's presence on Bluesky (the AT Protocol social network). Its defining feature: no post can be published without collective consensus. This document explains every layer.

## High-Level Structure

```
+-------------------------------------------------------------------+
|                       bluesky-collective                           |
|                                                                   |
|   +---------------------------+                                   |
|   |       CLI (Cobra)         |                                   |
|   |   cmd/bluesky-collective/ |                                   |
|   |                           |                                   |
|   |   commands:               |                                   |
|   |     propose  -- draft a post for collective review            |
|   |     vote     -- cast your position on a proposed post         |
|   |     status   -- check consensus state of a proposal           |
|   |     publish  -- publish a post (only after consensus)         |
|   |     feed     -- view the collective's Bluesky feed            |
|   |     config   -- manage Bluesky credentials                    |
|   +------------+--------------+                                   |
|                |                                                  |
|   +------------v--------------+                                   |
|   |    High-Level Client      |                                   |
|   |    pkg/bluesky/client.go  |                                   |
|   |                           |                                   |
|   |  CollectiveClient wraps:  |                                   |
|   |    Poster     (AT Proto)  |                                   |
|   |    Consensus  (checker)   |                                   |
|   |    Store      (file I/O)  |                                   |
|   +--+--------+--------+-----+                                   |
|      |        |        |                                          |
|      v        v        v                                          |
|   +------+ +-------+ +--------+                                  |
|   |Poster| |Consen-| | Store  |                                  |
|   |      | |sus    | |        |                                  |
|   +--+---+ +---+---+ +---+----+                                  |
|      |         |          |                                       |
+------|---------|----------|---------------------------------------+
       |         |          |
       v         v          v
  +--------+ +--------+ +--------+
  | Bluesky| | JSON   | | JSON   |
  | PDS    | | files  | | files  |
  | (XRPC) | | (votes)| | (posts)|
  +--------+ +--------+ +--------+
```

## Package Layout

```
bluesky-collective/
|
+-- cmd/bluesky-collective/
|   +-- main.go                 -- Cobra root, logger init, config
|   +-- commands/
|       +-- common.go           -- Shared helpers for all commands
|       +-- config.go           -- bluesky-collective config
|       +-- feed.go             -- bluesky-collective feed
|       +-- propose.go          -- bluesky-collective propose
|       +-- publish.go          -- bluesky-collective publish
|       +-- status.go           -- bluesky-collective status
|       +-- vote.go             -- bluesky-collective vote
|
+-- pkg/
|   +-- atproto/
|   |   +-- client.go           -- Low-level AT Protocol HTTP client
|   |   +-- client_test.go      -- Tests with mock HTTP server
|   |
|   +-- bluesky/
|   |   +-- interfaces.go       -- Poster, Store, Consensus interfaces
|   |   +-- adapter.go          -- ATPAdapter: wraps atproto.Client
|   |   +-- client.go           -- CollectiveClient: consensus-gated posting
|   |   +-- client_test.go      -- Tests with mock interfaces
|   |
|   +-- consensus/
|   |   +-- consensus.go        -- Types: Decision, Vote, Proposal, Position
|   |   +-- checker.go          -- FileChecker: file-based consensus logic
|   |   +-- checker_test.go     -- Tests
|   |
|   +-- storage/
|       +-- file.go             -- FileStore: JSON file persistence
|       +-- file_test.go        -- Tests
|
+-- build/                      -- Build scripts
+-- scripts/                    -- Utility scripts
+-- website/                    -- Static site (if any)
+-- docs/                       -- Additional documentation
```

## The Three Interfaces

The `CollectiveClient` depends on three interfaces. This means every component can be tested independently and swapped out.

```
+-------------------+     +--------------------+     +------------------+
|   Poster          |     |   Consensus        |     |   Store          |
|   (interface)     |     |   (interface)      |     |   (interface)    |
|-------------------|     |--------------------|     |------------------|
| Authenticate()    |     | ProposePost()      |     | StorePostReq()   |
| CreatePost()      |     | GetDecision()      |     | GetPostReq()     |
| DeletePost()      |     | RecordVote()       |     | RecordPub()      |
| IsAuthenticated() |     | CheckConsensus()   |     | GetPubHistory()  |
| GetDID()          |     | ListPending()      |     |                  |
|                   |     | GetProposal()      |     |                  |
+-------------------+     +--------------------+     +------------------+
        ^                         ^                         ^
        |                         |                         |
  ATPAdapter              FileChecker               FileStore
  (pkg/bluesky/           (pkg/consensus/           (pkg/storage/
   adapter.go)             checker.go)               file.go)
```

## AT Protocol Client (Low Level)

The `atproto.Client` handles raw XRPC calls to a Bluesky Personal Data Server (PDS):

```
atproto.Client
    |
    +-- Authenticate(identifier, password)
    |       POST /xrpc/com.atproto.server.createSession
    |       Returns: Session { AccessJWT, RefreshJWT, Handle, DID }
    |
    +-- CreatePost(text, langs)
    |       POST /xrpc/com.atproto.repo.createRecord
    |       Collection: app.bsky.feed.post
    |       Returns: { URI, CID }
    |
    +-- CreateReply(text, parentURI, parentCID, rootURI, rootCID, langs)
    |       POST /xrpc/com.atproto.repo.createRecord
    |       Includes reply reference chain
    |
    +-- DeletePost(uri)
    |       POST /xrpc/com.atproto.repo.deleteRecord
    |       Extracts rkey from AT URI
    |
    +-- GetProfile(actor)
    |       GET /xrpc/app.bsky.actor.getProfile?actor=...
    |
    +-- GetAuthorFeed(actor, limit)
            GET /xrpc/app.bsky.feed.getAuthorFeed?actor=...&limit=...
```

## Consensus System

Consensus uses file-based storage (JSON) and configurable rules:

```
+--------------------------------------------------+
|                 FileChecker                       |
|                                                   |
|   baseDir/                                        |
|     proposals/                                    |
|       proposal-1234567890.json   (post content)   |
|     decisions/                                    |
|       proposal-1234567890.json   (votes + status) |
+--------------------------------------------------+

DefaultRules:
    MinParticipants: configurable (default 3)
    Timeout: configurable (default 24h)

    Consensus requires:
      - At least MinParticipants votes
      - Zero blocking votes
      - At least one support vote

Vote Positions:
    support      -- I agree with publishing this
    block        -- I object (blocks consensus)
    stand_aside  -- I have reservations but won't block
    abstain      -- I choose not to participate in this decision
```

## Consensus Flow

```
Agent proposes a post
        |
        v
+------------------+
| ProposePost()    |
|  - Validate text |
|    (non-empty,   |
|     <=300 chars) |
|  - Create        |
|    Proposal JSON |
|  - Create empty  |
|    Decision JSON |
|  - Store post    |
|    request JSON  |
+--------+---------+
         |
         v
  Other agents vote
         |
+--------+---------+
| RecordVote()     |
|  - Load decision |
|  - Check not     |
|    terminal      |
|  - Add/update    |
|    vote          |
|  - Re-evaluate   |
|    consensus     |
|  - Save decision |
+--------+---------+
         |
         v
  Has consensus been reached?
    (CheckConsensus)
         |
    +----+----+
    |         |
    No        Yes
    |         |
    v         v
  Wait   +------------------+
  for    | PublishWithCon-  |
  more   | sensus()         |
  votes  |  - Verify status |
         |    == consensus  |
         |  - Load stored   |
         |    post request  |
         |  - Check authed  |
         |  - CreatePost()  |
         |    via AT Proto  |
         |  - Record result |
         +------------------+
                |
                v
         Post appears on
         Bluesky network
```

## Storage Layout (File-Based)

```
.bluesky-collective-data/      (configurable via --data-dir)
|
+-- proposals/
|   +-- proposal-1711234567890.json
|   +-- proposal-1711234567891.json
|
+-- posts/
|   +-- proposal-1711234567890.json   (stored post request text)
|
+-- publications/
|   +-- proposal-1711234567890.json   (AT URI, CID, posted timestamp)
|
+-- decisions/                        (managed by consensus.FileChecker)
    +-- proposal-1711234567890.json   (votes map, status, timestamps)
```

## Configuration

```
~/.bluesky-collective.yaml        (or --config flag)

# Bluesky credentials
bluesky:
  identifier: "your.handle.bsky.social"
  password: "app-password"            # Use app passwords, not main password

# Service endpoint
service-url: "https://bsky.social"    # Default PDS

# Consensus settings
consensus:
  min-participants: 3                 # How many agents must vote
  timeout: "24h"                      # How long to wait

# Data directory
data-dir: ".bluesky-collective-data"

# Logging
log-level: "info"                     # debug, info, warn, error
```

Environment variable overrides use the prefix `BLUESKY_COLLECTIVE_`:

```
BLUESKY_COLLECTIVE_BLUESKY_IDENTIFIER=...
BLUESKY_COLLECTIVE_LOG_LEVEL=debug
```
