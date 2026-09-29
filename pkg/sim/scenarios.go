package sim

// DiffBlock represents a file edit with unified diff lines
type DiffBlock struct {
	Filename string
	Lines    []string
}

// Scenario represents a complete realistic Claude Code engineering workflow
type Scenario struct {
	Prompt        string
	ScoutActions  []struct{ Tool, Args, Result string }
	Thoughts      []string
	Diffs         []DiffBlock
	TestCommand   string
	TestOutputs   []string
	ReviewNotes   []string
	CommitMessage string
	BranchName    string
	DeployTarget  string
	Conclusion    string
}

// Scenarios holds our collection of convincing developer tasks
var Scenarios = []Scenario{
	{
		Prompt: "Fix goroutine leak and lock contention in distributed connection pool",
		ScoutActions: []struct{ Tool, Args, Result string }{
			{"Grep", `"AcquireConn" path="internal/pool"`, "Found 9 matches across 3 files"},
			{"Read", `"internal/pool/manager.go" (lines 1-145)`, "Read 145 lines (4.8 KB)"},
			{"Read", `"internal/pool/worker_test.go" (lines 1-90)`, "Read 90 lines (2.9 KB)"},
		},
		Thoughts: []string{
			"Analyzing channel lifecycle in ConnectionPool.AcquireConn...",
			"Detected unbuffered channel send when context expires before worker select...",
			"Evaluating ring-buffer vs sync.Pool for idle socket recycling...",
			"Formulating lock-free atomic state transition to prevent mutex convoy...",
		},
		Diffs: []DiffBlock{
			{
				Filename: "internal/pool/manager.go",
				Lines: []string{
					"@@ -64,11 +64,16 @@ func (p *Pool) Acquire(ctx context.Context) (*Conn, error) {",
					" \tselect {",
					" \tcase conn := <-p.idleCh:",
					"-\t\tif conn.IsStale() {",
					"-\t\t\tconn.Close()",
					"-\t\t\treturn p.dialNew(ctx)",
					"-\t\t}",
					"+\t\tif conn.IsStale() || !conn.Ping(ctx) {",
					"+\t\t\t_ = conn.Close()",
					"+\t\t\tp.metrics.StaleEvictions.Add(1)",
					"+\t\t\treturn p.dialWithBackoff(ctx)",
					"+\t\t}",
					" \t\treturn conn, nil",
					" \tcase <-ctx.Done():",
					"-\t\treturn nil, ctx.Err()",
					"+\t\tp.metrics.ContextTimeouts.Add(1)",
					"+\t\treturn nil, fmt.Errorf(\"pool acquire aborted: %w\", ctx.Err())",
					" \t}",
					" }",
				},
			},
		},
		TestCommand: "go test -v -race -count=1 ./internal/pool/...",
		TestOutputs: []string{
			"=== RUN   TestPoolAcquireConcurrency",
			"=== PAUSE TestPoolAcquireConcurrency",
			"=== CONT  TestPoolAcquireConcurrency",
			"--- PASS: TestPoolAcquireConcurrency (0.14s)",
			"=== RUN   TestPoolContextCancellationNoLeak",
			"--- PASS: TestPoolContextCancellationNoLeak (0.08s)",
			"=== RUN   TestPoolStaleSocketEviction",
			"--- PASS: TestPoolStaleSocketEviction (0.05s)",
			"PASS",
			"ok  \tgithub.com/acme/core/internal/pool\t0.412s (coverage: 94.2% of statements)",
		},
		ReviewNotes: []string{
			"Security: Verified no socket file descriptor leaks on context cancellation",
			"Performance: Reduced p99 lock wait latency from 14.2ms to 0.3ms",
			"Memory: Confirmed zero orphaned goroutines via runtime.NumGoroutine() assertion",
		},
		CommitMessage: "fix(pool): prevent goroutine leak and stale socket reuse on context timeout",
		BranchName:    "claude/fix-pool-goroutine-leak",
		DeployTarget:  "staging-us-east-1 (12/12 pods healthy)",
		Conclusion:    "I've resolved the goroutine leak in `internal/pool/manager.go` by ensuring stale connections are properly closed and context timeouts increment telemetry counters without blocking the worker channel. All race detector tests pass with 94.2% coverage.",
	},
	{
		Prompt: "Patch timing side-channel vulnerability in JWT HMAC signature verification",
		ScoutActions: []struct{ Tool, Args, Result string }{
			{"Grep", `"VerifySignature" path="pkg/auth"`, "Found 5 matches across 2 files"},
			{"Read", `"pkg/auth/jwt.go" (lines 40-120)`, "Read 80 lines (3.1 KB)"},
			{"Glob", `"pkg/auth/**/*_test.go"`, "Found 4 test files"},
		},
		Thoughts: []string{
			"Inspecting byte slice comparison in pkg/auth/jwt.go...",
			"Standard bytes.Equal short-circuits on first mismatched byte, leaking timing info...",
			"Replacing with crypto/subtle.ConstantTimeCompare to guarantee constant-time evaluation...",
			"Checking key rotation fallback logic to ensure uniform execution branches...",
		},
		Diffs: []DiffBlock{
			{
				Filename: "pkg/auth/jwt.go",
				Lines: []string{
					"@@ -82,8 +82,9 @@ func (v *Verifier) VerifyToken(raw string) (*Claims, error) {",
					" \texpectedMAC := computeHMAC(payload, v.secretKey)",
					"-\tif string(signature) != string(expectedMAC) {",
					"-\t\treturn nil, ErrInvalidSignature",
					"+\tif subtle.ConstantTimeCompare(signature, expectedMAC) != 1 {",
					"+\t\tv.auditLog.RecordFailure(\"hmac_mismatch\")",
					"+\t\treturn nil, ErrInvalidSignature",
					" \t}",
					" \treturn parseClaims(payload)",
				},
			},
		},
		TestCommand: "go test -v -race ./pkg/auth/...",
		TestOutputs: []string{
			"=== RUN   TestVerifyToken_ValidHMAC",
			"--- PASS: TestVerifyToken_ValidHMAC (0.01s)",
			"=== RUN   TestVerifyToken_TimingInvariant",
			"--- PASS: TestVerifyToken_TimingInvariant (0.22s)",
			"=== RUN   TestVerifyToken_ExpiredClaims",
			"--- PASS: TestVerifyToken_ExpiredClaims (0.01s)",
			"PASS",
			"ok  \tgithub.com/acme/core/pkg/auth\t0.289s (coverage: 98.1% of statements)",
		},
		ReviewNotes: []string{
			"Security (CVE-Check): Eliminated byte-by-byte timing side-channel (CWE-208)",
			"Compliance: Constant-time comparison verified across valid and invalid key lengths",
		},
		CommitMessage: "security(auth): enforce constant-time HMAC comparison using crypto/subtle",
		BranchName:    "claude/security-hmac-constant-time",
		DeployTarget:  "auth-edge-canary (8/8 replicas ready)",
		Conclusion:    "Replaced the non-constant-time string equality check in `pkg/auth/jwt.go` with `subtle.ConstantTimeCompare` to mitigate timing side-channel attacks (CWE-208). Added timing invariance unit tests.",
	},
	{
		Prompt: "Eliminate N+1 query bottleneck in GraphQL user organization resolver",
		ScoutActions: []struct{ Tool, Args, Result string }{
			{"Grep", `"ResolveOrganizations" path="internal/graphql"`, "Found 7 matches across 4 files"},
			{"Read", `"internal/graphql/resolvers/user.go" (lines 1-160)`, "Read 160 lines (5.4 KB)"},
			{"Read", `"internal/dataloader/org_loader.go" (lines 1-85)`, "Read 85 lines (2.6 KB)"},
		},
		Thoughts: []string{
			"Tracing SQL queries emitted during nested User -> Organizations GraphQL query...",
			"Each user node triggers an individual SELECT * FROM organizations WHERE user_id = $1...",
			"Integrating batched DataLoader with 2ms coalescing window and request-scoped memoization...",
			"Verifying SQL ANY($1) array binding compatibility with pgxv5 driver...",
		},
		Diffs: []DiffBlock{
			{
				Filename: "internal/graphql/resolvers/user.go",
				Lines: []string{
					"@@ -47,9 +47,10 @@ func (r *userResolver) Organizations(ctx context.Context, obj *model.User) ([]*m",
					"-\torgs, err := r.db.QueryOrganizationsByUserID(ctx, obj.ID)",
					"-\tif err != nil {",
					"-\t\treturn nil, err",
					"-\t}",
					"-\treturn orgs, nil",
					"+\tloader := dataloader.For(ctx).OrganizationsByUser",
					"+\torgs, err := loader.Load(ctx, obj.ID)",
					"+\tif err != nil {",
					"+\t\treturn nil, fmt.Errorf(\"dataloader orgs for user %s: %w\", obj.ID, err)",
					"+\t}",
					"+\treturn orgs, nil",
				},
			},
		},
		TestCommand: "go test -v ./internal/graphql/... -bench=BenchmarkOrgResolver",
		TestOutputs: []string{
			"=== RUN   TestUserOrganizationsBatchedQuery",
			"    user_test.go:54: Executed 1 SQL batch query for 250 user nodes (previously 251 queries)",
			"--- PASS: TestUserOrganizationsBatchedQuery (0.03s)",
			"goos: darwin",
			"goarch: arm64",
			"BenchmarkOrgResolver-12    \t   14280\t     81420 ns/op\t    4120 B/op\t      18 allocs/op",
			"PASS",
			"ok  \tgithub.com/acme/core/internal/graphql\t1.320s",
		},
		ReviewNotes: []string{
			"Database: Query count for 250 items reduced from 251 to 1 (99.6% reduction)",
			"Throughput: p95 GraphQL response time improved from 480ms to 11ms",
		},
		CommitMessage: "perf(graphql): batch user organization lookups via request-scoped DataLoader",
		BranchName:    "claude/perf-graphql-dataloader",
		DeployTarget:  "api-gateway-staging (24/24 pods healthy)",
		Conclusion:    "Refactored the `Organizations` field resolver to batch database lookups through the context-bound `DataLoader`. For a 250-node payload, SQL round-trips dropped from 251 to 1, cutting p95 latency from 480ms to 11ms.",
	},
	{
		Prompt: "Implement zero-downtime vector index migration with HNSW parameters",
		ScoutActions: []struct{ Tool, Args, Result string }{
			{"Glob", `"migrations/*.sql"`, "Found 28 migration files"},
			{"Read", `"migrations/0028_embeddings_idx.sql" (lines 1-45)`, "Read 45 lines (1.4 KB)"},
			{"Grep", `"vector_cosine_ops" path="internal/search"`, "Found 6 matches"},
		},
		Thoughts: []string{
			"Standard CREATE INDEX blocks writes on embeddings table (42M rows)...",
			"Switching to CREATE INDEX CONCURRENTLY with pgvector hnsw (m=16, ef_construction=64)...",
			"Setting maintenance_work_mem = '2GB' locally within session to prevent disk spill...",
			"Adding fallback query planner hint if index build is still in phase 2...",
		},
		Diffs: []DiffBlock{
			{
				Filename: "migrations/0029_hnsw_concurrent_idx.sql",
				Lines: []string{
					"@@ -1,6 +1,11 @@",
					"-- migrate:up",
					"-CREATE INDEX idx_documents_embedding ON documents",
					"-USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);",
					"+-- migrate:up transaction:false",
					"+SET LOCAL maintenance_work_mem = '2GB';",
					"+CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_documents_embedding_hnsw",
					"+ON documents USING hnsw (embedding vector_cosine_ops)",
					"+WITH (m = 16, ef_construction = 64);",
					"+DROP INDEX CONCURRENTLY IF EXISTS idx_documents_embedding;",
				},
			},
		},
		TestCommand: "pytest tests/search/test_vector_index.py -v",
		TestOutputs: []string{
			"============================= test session starts ==============================",
			"platform darwin -- Python 3.12.4, pytest-8.2.2",
			"collected 4 items",
			"tests/search/test_vector_index.py::test_hnsw_recall_accuracy PASSED      [ 25%]",
			"tests/search/test_vector_index.py::test_concurrent_writes_unblocked PASSED [ 50%]",
			"tests/search/test_vector_index.py::test_cosine_distance_ordering PASSED  [ 75%]",
			"tests/search/test_vector_index.py::test_index_scan_explain_plan PASSED   [100%]",
			"============================== 4 passed in 0.64s ===============================",
		},
		ReviewNotes: []string{
			"Migration Safety: Disabled transaction block for CREATE INDEX CONCURRENTLY",
			"Recall@10: Improved from 0.89 (IVFFlat) to 0.985 (HNSW m=16)",
		},
		CommitMessage: "perf(search): migrate pgvector index to concurrent HNSW with m=16",
		BranchName:    "claude/vector-hnsw-migration",
		DeployTarget:  "search-cluster-canary (6/6 shards synced)",
		Conclusion:    "Updated the vector index migration to use `CREATE INDEX CONCURRENTLY` with HNSW (`m = 16, ef_construction = 64`). This avoids table write locks during deployment and boosts Recall@10 from 89.0% to 98.5%.",
	},
	{
		Prompt: "Fix memory leak in WebSocket frame decompression buffer pool",
		ScoutActions: []struct{ Tool, Args, Result string }{
			{"Grep", `"flate.NewReader" path="pkg/ws"`, "Found 4 matches across 2 files"},
			{"Read", `"pkg/ws/compression.go" (lines 1-115)`, "Read 115 lines (3.7 KB)"},
		},
		Thoughts: []string{
			"Inspecting heap profile: 68% of live heap retained by bytes.Buffer in ws.decompressFrame...",
			"Buffers growing above 64KB for large payloads are returned to sync.Pool without capacity cap...",
			"Pinning maximum pooled buffer capacity to 64KB and resetting reader state before Put()...",
		},
		Diffs: []DiffBlock{
			{
				Filename: "pkg/ws/compression.go",
				Lines: []string{
					"@@ -52,6 +52,10 @@ func releaseBuffer(buf *bytes.Buffer) {",
					"+\tconst maxPooledBufferCap = 64 * 1024 // 64 KiB",
					"+\tif buf.Cap() > maxPooledBufferCap {",
					"+\t\treturn // Let GC reclaim oversized buffers",
					"+\t}",
					" \tbuf.Reset()",
					" \tbufferPool.Put(buf)",
					" }",
				},
			},
		},
		TestCommand: "go test -v -memprofile=mem.out ./pkg/ws/...",
		TestOutputs: []string{
			"=== RUN   TestDecompressFrame_OversizedBufferDiscarded",
			"--- PASS: TestDecompressFrame_OversizedBufferDiscarded (0.02s)",
			"=== RUN   TestDecompressFrame_PoolReuse",
			"--- PASS: TestDecompressFrame_PoolReuse (0.04s)",
			"PASS",
			"ok  \tgithub.com/acme/core/pkg/ws\t0.195s (coverage: 96.4% of statements)",
		},
		ReviewNotes: []string{
			"Memory: Steady-state RSS under 10k WebSocket connections dropped from 4.2 GB to 310 MB",
			"Security: Protected against zip-bomb buffer retention attacks",
		},
		CommitMessage: "fix(ws): cap pooled decompression buffers at 64KiB to prevent heap retention",
		BranchName:    "claude/fix-ws-buffer-pool-leak",
		DeployTarget:  "realtime-gateway-prod (32/32 nodes healthy)",
		Conclusion:    "Added a 64 KiB capacity guard in `releaseBuffer` before returning buffers to `sync.Pool`. Oversized buffers from transient large frames are now collected by the GC instead of permanently inflating the pool.",
	},
}
