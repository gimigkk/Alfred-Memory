# AI Skill: The Courtroom Simulation (V33)

## Prior Rulings
* **V32:** Modified Rule 18 to forbid `CREATE_NODE` for Events when a matching unclarified Event exists, forcing updates over duplication.

## The Proposal on Trial
**Issue:** The agent exhibits "tunnel vision" during the Discovery phase. It prioritizes resolving speaker names, sometimes throws in one obvious semantic query (e.g., "rapat bph gacoan"), and then immediately calls `[REQUEST_SCHEMA]`. This severely underutilizes the Hybrid RAG engine, missing broader context clues hidden in nouns, verbs, acronyms, and organizational jargon.

**Proposal:** Rewrite `skill_discovery.md` to force a mandatory `[SEMANTIC ANALYSIS]` thought block *after* speaker resolution. The agent must break down the transcript to brainstorm non-person concepts, and then aggressively execute `query_rag` on those concepts (using `""` for `target_speakers`) before being allowed to request the schema.

## The Personas
1. **The Prompt Engineer (Pro-Context):** Believes forcing the LLM to "think out loud" about abstract nouns will drastically improve vector retrieval hit rates.
2. **The Software Architect (Pro-Structure):** Wants rigid schema alignment.
3. **The DevOps/SRE (Pro-Efficiency):** Worries about token bloat, latency, and DB load.
4. **The Hostile Attacker (Prosecutor):** Believes brainstorming is an uncontrollable hallucination vector that will flood the API with garbage queries.
5. **The Judge:** Delivers the final verdict.

---

## Cycle 1: The Threat of Query Spam

**The Prompt Engineer:**

If we just tell the agent "Query for all hints!", it won't do it. We have to force it to write a list first. "Phase 3: Write down 5 keywords or concepts from the text that aren't people. Phase 4: Query all of them."

**The Hostile Attacker (Prosecutor):**

If you force it to write 5 keywords, it will query for useless conversational noise just to fill the quota. "tf", "ok", "gw", "lu", "besok". Do you want your Hybrid RAG searching the vector index for the word "ok"? You are going to blow up the Go backend with useless I/O, increase latency, and pollute the context window with completely irrelevant nodes that happened to contain the word "ok"!

**The DevOps/SRE:**

The Attacker is absolutely right. We already hit rate limits occasionally. If the LLM starts batching 10 garbage queries per ingestion loop, the LadybugDB query latency will spike, and the prompt context will fill up with useless `NO_MATCH` or, worse, false-positive node dumps that derail the commit phase.

**The Prompt Engineer:**

We can constrain the brainstorming. We tell it to extract *proper nouns, acronyms, and project-specific jargon* only. Sarcasm, slang, and common verbs are explicitly forbidden. 

**Outcome:**
**RULING CHANGED:** The proposal is modified. Brainstorming cannot be an arbitrary quota (e.g., "find 5 words"). It must be strictly filtered to "Proper Nouns, Acronyms, and Project/Event Jargon."

---

## Cycle 2: Managing the `target_speakers` Array

**The Software Architect:**

Wait, how is the agent going to query these concepts? The `query_rag` tool in `orchestrator.go` enforces a strict 1:1 match between the `queries` array and the `target_speakers` array to ensure speaker resolution is auditable. 

**The Hostile Attacker (Prosecutor):**

I know exactly what will happen. The agent will try to batch query its speakers AND its concepts at the same time:
`queries: ["apta", "reimburse gacoan"]`
`target_speakers: ["apta_ieee25"]`
The arrays mismatch, the Go Orchestrator rejects the tool call, and the loop crashes! Or, the agent will try to fake a speaker for the concept: `target_speakers: ["apta_ieee25", "NONE"]`. Then Go rejects it because "NONE" is not in the manifest!

**The Software Architect:**

Actually, `tool_handlers.go` explicitly allows an empty string `""` in the `target_speakers` array for semantic queries. If the LLM sends `""`, the orchestrator audits it as `Target: NONE` and allows it.

**The Hostile Attacker (Prosecutor):**

Then you MUST explicitly teach the agent how to do this in the `skill_discovery.md` prompt. Right now, Rule 4 says: "You MUST provide a target_speakers array of the exact same length... mapping each speaker's literal label". It barely mentions the empty string. If you don't separate Speaker Resolution from Concept Discovery, the LLM will mix them up and fail the array validation.

**Outcome:**
**RULING CHANGED:** The prompt must explicitly separate the discovery process into distinct, sequential tool calls. Call 1: Speakers (using literal labels in `target_speakers`). Call 2: Concepts (using an array of `""` in `target_speakers`).

---

## Cycle 3: The Danger of "Deep Digging"

**The Prompt Engineer:**

The proposal also includes "Recursive Context Gathering." If a query returns an organization (like "BEM"), the agent should query "BEM" to find out what it is.

**The Hostile Attacker (Prosecutor):**

No! This is an infinite loop trap. If it queries "BEM" and gets a node that says "BEM is led by Bahlil", it will query "Bahlil", which returns "Bahlil is in IEEE", so it queries "IEEE". You are building an autonomous web crawler inside an ingestion loop that is supposed to finish in 5 seconds. The Go ReAct interceptor will prune the history, but the latency will be catastrophic.

**The DevOps/SRE:**

Agreed. We are building a memory vault, not a research agent. The ingestion pipeline must be bounded. We cannot afford recursive LLM loops in Phase 1 (Discovery).

**The Prompt Engineer:**

I concede. We will drop the "recursive" requirement. The semantic discovery must be a single, broad batch query based *only* on the raw transcript text, not on the results of previous queries.

**Outcome:**
**RULING CHANGED:** "Recursive Context Gathering" is stripped from the proposal to prevent infinite loops. Semantic discovery is limited to a single pass of broad queries derived directly from the transcript manifest.

---

## The Verdict

**Findings:**
1. The agent suffers from tunnel vision because the prompt emphasizes speaker resolution and does not provide a structured thought process for extracting abstract concepts.
2. Forcing the LLM to brainstorm without constraints will result in API spam and vector search pollution (e.g., querying for conversational noise).
3. Mixing speaker queries and concept queries in a single tool call risks failing the Go Orchestrator's strict `target_speakers` array length validation.
4. Recursive querying is an infinite-loop hazard that must be avoided.

**Directives:**
- Rewrite `assets/prompts/skill_discovery.md` into strict chronological phases.
- **Phase 2 (Speaker Resolution):** The agent batches all speakers.
- **Phase 3 (Semantic Brainstorming & Query):** The agent uses an `[AGENT THOUGHT]` block to explicitly list out "Proper Nouns, Acronyms, and Project Jargon" from the transcript. It then executes a separate `query_rag` call for these concepts, using an array of `""` (empty strings) for `target_speakers`.
- Explicitly forbid recursive querying to bound the ingestion latency.
