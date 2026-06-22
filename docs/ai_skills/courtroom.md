# AI Skill: The Courtroom Simulation

When the user requests a "courtroom" or "courtroom simulation" to debate a technical or architectural decision, execute the simulation strictly according to the following constraints.

## 1. The Setup
- Identify the exact Proposal on Trial.
- Establish a minimum of 5 Personas representing different engineering perspectives (e.g., The Software Architect, The Prompt Engineer, The DevOps/SRE, The Security Auditor, The Judge).
- **MANDATORY PERSONA:** You MUST always include **The Hostile Attacker (Prosecutor)**. This persona must be aggressively cynical, constantly looking for prompt injection loopholes, state transition vulnerabilities, architectural flaws, and semantic contradictions. They exist to shatter echo-chambers and force the other personas to defend their logic against worst-case scenarios.
- Explicitly state each persona's inherent bias or priority before the debate begins.

## 2. The Execution (Cycles)
- The debate MUST run for a minimum of 3 full cycles.
- Each cycle consists of the personas directly attacking or rebutting the arguments made by the opposing personas in the previous cycle.

### 2a. Back-and-Forth, Not Monologues
Within a cycle, arguments must move as a **real exchange**, not a sequence of standalone speeches:
- Break each persona's contribution into **short turns** (2-4 sentences each), not single long paragraphs. If a point needs more than 4 sentences, split it across multiple turns and let another persona respond in between.
- A persona may **interrupt or directly respond mid-cycle** to a point just made — the next speaker should reference the specific claim before it, not restate their general position.
- Avoid the pattern of "Persona A gives full argument → Persona B gives full counter-argument → done." Instead aim for: claim → challenge → counter-challenge → concession or rebuttal → (repeat as needed) within the same cycle.
- A cycle should feel like a conversation under cross-examination, not a panel of prepared statements read in sequence.

### 2b. Cycle Outcome Declaration
At the end of every cycle, explicitly declare one of the following outcomes before moving to the next cycle:
- **"RULING CHANGED:"** state precisely what changed in the proposal/design and why the objection forced it.
- **"RULING UPHELD, RISK ACCEPTED:"** state precisely what residual risk remains unresolved and why it is being tolerated rather than fixed.
A cycle that produces neither outcome is not valid — it is a restatement, not a resolution. Re-run it or cut it from the transcript. Do not allow a cycle to read as progress when nothing about the proposal actually moved.

### 2c. Equal-Scrutiny Rule
The Hostile Attacker must apply identical scrutiny to claims about **existing/unmodified systems** ("the current code already handles this," "the existing limit already protects against this") as to claims about newly proposed code. An unverified assertion about legacy or pre-existing behavior is not evidence — it is a claim that has not yet been tested in this session. If the Attacker does not challenge a given claim, the Judge's verdict must explicitly note that the point was **unchallenged**, not vindicated. Unanimous agreement reached without a real attempt to break a claim does not constitute approval of that claim.

## 3. The Evidence Rule
- This rule applies when the proposal's claims are **empirical** in nature — cost comparisons, performance benchmarks, industry practice, statistical track records. In these cases, EVERY argument MUST be backed by researched evidence: use web search tools to find real, external validation (academic papers, official documentation, industry articles, or cloud economic data), cite it inline, and provide a references section at the bottom.
- This rule does **not** apply when the proposal is a matter of **internal logic or correctness** (state machine design, algorithmic soundness, code-level edge cases) where the relevant "evidence" is the reasoning and code itself, not external sources. In these cases, say so explicitly at the start of the session (e.g., "This proposal is a pure correctness question; external evidence is not applicable") rather than silently dropping the Evidence Rule while it remains nominally "mandatory." Do not fabricate citations to satisfy this rule when it does not apply.

## 4. The Verdict (The Judge's Constraints)
- The Judge persona synthesizes the arguments and the empirical data (where applicable).
- **ANTI-FALLACY RULE:** The Judge MUST strictly evaluate sample sizes and statistical significance. The Judge is strictly forbidden from committing Hasty Generalizations (e.g., weighing a single $n=1$ success over a historical $n=100$ failure rate).
- The Judge must prioritize documented systemic behavior over isolated anomalies.
- **UNCHALLENGED-CLAIM FLAG:** Per Section 2c, the Judge must explicitly list any claim in the session that the Hostile Attacker did not substantively contest, and must not treat silence as confirmation. Cycles flagged this way should be called out as open risk in the verdict, not folded into "unanimous approval."
- The Judge delivers a final, definitive verdict approving or denying the proposal based on sound logical and statistical reasoning.

## 5. Formatting Requirements (Readability is Mandatory)
This is a wall-of-text risk by nature — formatting discipline is not optional polish, it is required for the output to be usable.

- Output the simulation as a dedicated Markdown Artifact.
- **Insert a blank line gap before and after every persona's turn.** Never let two different personas' text touch across a single line break — there must be visible vertical whitespace separating every speaker change, every time, with no exceptions.
- **Bold the persona name at the start of every turn**, on its own line, followed by a blank line, then the dialogue. Format each turn as:

  ```
  **Persona Name:**

  Their dialogue text here, kept short per Section 2a.

  ```

- Use a horizontal rule (`---`) between full cycles only — not between every individual turn within a cycle, or the document becomes choppier than it needs to be. Turn-to-turn separation is handled by blank-line spacing; cycle-to-cycle separation is handled by `---`.
- Keep individual turns short (2-4 sentences) per Section 2a — this is also a formatting requirement, since long unbroken paragraphs are what make the transcript hard to read in the first place.
- Use a clear header for each cycle (e.g., `## Cycle 1: The Title Of This Cycle`) so the document is scannable without reading every line.
- The Verdict section should be visually distinct — use a header and consider bolding key rulings as a short list rather than prose, so the outcome is skimmable even if the debate above it is dense.

## 6. Multi-Session Continuity
If this courtroom continues a proposal debated in a prior session (e.g., "Courtroom V21" following "V20"), open the transcript with a **"Prior Rulings"** section listing every decision already made in earlier sessions on this same proposal, in plain bullet form.
- Any new cycle that revisits a question already settled by a prior ruling MUST explicitly state whether it **upholds** or **overturns** that ruling, and why.
- Silently restating a position opposite to a prior verdict — without acknowledging the prior verdict exists — is prohibited. If this is detected mid-session, the session must pause and reconcile the contradiction before continuing.
- If no prior session exists, state plainly that this is a first-pass review with no prior rulings to reconcile.

## 7. Recommended Use
This skill is best suited to proposals with clear, falsifiable failure modes — state machines, gating logic, security boundaries, architectural tradeoffs with measurable cost/latency implications. It is not a substitute for actually running the code; the verdict identifies what to fix or verify next, not a guarantee of correctness.