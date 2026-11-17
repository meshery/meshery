State-Level Gold Medallist — Archery · Young Changemakers, IIT Delhi 
EXPERIENCE 
QA Intern — SignSetu (Forbes Accessibility 200 Honouree 2026) 
Remote | Jun 2026 – Sep 2026 
• Took ownership of reliability across the mobile, payments, auth and data layers of SignSetu's accessible-learning app for the Deaf community, tracing 
user-facing failures to root causes and working with developers to ship the fixes. 
• Mobile concurrency: OTA updates were silently failing on slow networks; found an uncancelled Promise leak bypassing timeout wrappers and 
resolved JavaScript-thread vs. native-scheduler race conditions, restoring reliable OTA scheduling across platforms. 
• Payments and security: Purchase stacking broke from an ordinal conflict in payment records; added webhook idempotency guards and schema 
validation, eliminating duplicate ledger writes. Uncovered a validation gap that let users buy lessons beyond the available content, a financial risk and a 
source of confusion for users with limited English. Also closed an email auto-linking vector that exposed admin accounts to takeover. 
• Data integrity: Progress writes were race-prone and nearing MongoDB's document size limit; removed non-atomic read-modify-write conflicts and 
bounded array growth, preventing data loss and overflow. 
• Scaling the team's QA: Manual review couldn't keep pace with feature velocity, so I built a 2-agent E2E test generation framework (Playwright, 
Maestro) with durable Markdown handoffs for CI/CD. The developers and CTO then adopted it and expanded it, with many further additions, into a 
company-wide testing harness. 
• Recommended in writing by the CTO (Sep 2026) for technical depth, business awareness and user empathy, and for communicating findings clearly to 
developers. 
MERN Stack Intern — Infotact Solutions 
Remote / Chennai · May 2025 – Aug 2025 
• Malformed entries were reaching the application layer unchecked; built 15+ reusable React components and designed 4+ NoSQL schemas with 
MongoDB and Mongoose validation, eliminating malformed-entry errors before they reached the database. 
PROJECTS & RESEARCH 
Algorithmic Alpha Engine (BrainQuant Generator) — WorldQuant BRAIN 
Python · Polars · FastAPI · Active Development 
• Problem: Hand-testing trading signals is slow and prone to overfitting. I designed and built an end-to-end discovery pipeline that natively compiles AST 
alpha formulas into Polars (pl.Expr) vector expressions and backtests them against Parquet snapshots of 500+ global equities, pulled by a concurrent 
8-worker pipeline (ThreadPoolExecutor) with retry and rate-limit back-off. 
• Correctness: Enforced strict dimensional parity (time-series operators per ticker, cross-sectional operators per date) for zero lookahead bias, plus daily 
demeaning and Delay-1 alignment, so local results carry over to the live BRAIN platform. 
• Discovery strategy: Built a hybrid engine splitting each batch across 6 market-archetype templates (40%), genetic evolution with crossover and 
subtree mutation (40%), and near-miss healing (20%), running continuously in the background. 
• Rigor: Validated every candidate on a 70/30 in-sample / out-of-sample split and raised the Sharpe threshold as more candidates were tested (base + 
0.3 × log10(N)) to remove data-mining bias. 
• Efficiency and visibility: Designed a 2-tier cache (formula-metrics cache plus a sub-tree column cache with 30-column / 2 GB caps) that holds memory 
flat under 1.2 GB over indefinite runs, with a live dashboard for throughput, latency, cache hit rate and RSS. 
• Differentiation: OHLCV data is among the most crowded in quantitative finance; engineered a micro-blend alpha-construction technique with near-zero 
correlation to millions of live alphas, using pure mathematical feature design rather than alternative datasets. 
• Impact: A submitted alpha passed every BRAIN check on USA/D1/TOP3000 (Sharpe 1.63, fitness 1.53, turnover 1.09%, 5.52% max drawdown, 203 
bps margin, balanced ~1,561 long / 1,473 short book). Screening 500+ candidates earned the WorldQuant Challenge Gold Level certificate and a 
top-10,000 rank among 170,000+ participants worldwide. 
Handwritten Mathematical Proof Verifier — Published Patent, Lead Co-Inventor 
App No: 202641043577 · 2026 
• Problem: Students and educators reason on paper, but formal verifiers (Z3, Lean 4) demand expert syntax, and no tool could tell a genuine mistake 
from a deliberate contradiction in a proof by contradiction. As lead co-inventor, I designed the core verification logic to close that gap. 
• Assertion-level logic: Each handwritten step is parsed into an AST and converted into a Boolean implication (step N-1 implies step N), then evaluated 
by a CAS truth engine against a Z3 / Lean 4 axiom library, verifying reasoning rather than just computing numbers. 
• Intention-aware analysis: When a step evaluates false, the system scans the proof history for markers such as "assume contradiction" and classifies it 
as a valid contradiction or a calculation error, so correct proofs are never penalised for how they argue. 
• Architecture: A neural perception engine (DenseNet-121 OCR and ViT-Base spatial layout, fused by an arbiter node with softmax confidence scoring) 
feeds a symbolic engine (FOL logic converter, schema validator, SMT solver), returning results as a real-time overlay on the student's own handwriting. 
• Impact: Fully sponsored by VIT; the application is published with 10 claims covering both the system and the method. 
TECHNICAL SKILLS 
Languages: Python (Proficient), Java, C++, JavaScript, SQL 
Systems & Backend: FastAPI, Node.js, Express.js, Docker, Kubernetes, Linux, Async/Concurrent Programming, REST APIs 
Data & ML: Polars, pandas, NumPy, scikit-learn, scipy, NetworkX, PostgreSQL, Supabase, MongoDB 
Dev Practices: Git, CI/CD, Agile / Scrum, Pytest, Jest · Frontend: React.js
