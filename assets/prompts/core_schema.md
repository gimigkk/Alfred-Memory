## Schema Constraints
When creating or updating nodes, you must only use properties and edges defined below. 
**CRITICAL EDGE RULE:** The `add_edges` array ALWAYS belongs to the mutation of the node where the edge ORIGINATES. 

- **Person**: `name`, `aliases`, `phone_number` (NEVER output `content`)
  - **Outgoing Edges (Allowed):**
    - `ASSIGNED_TO` -> Task (Target is the task they bear the burden of executing).
    - `MENTIONED_IN` -> Task/Event (Target is the task/event they are a beneficiary, issuer, subject, or passive participant of. You MUST link a task's subject/beneficiary here even if the task is assigned to someone else, and even if the subject never speaks in the transcript. Do not leave discussed subjects unlinked).
    - `HAS_ROLE` -> Event (Target is the event they have a titled role in).
    - `MEMBER_OF` -> Circle (Target is the named group or division they belong to. You MUST include a `role` property on this edge if stated, e.g., "kadiv", "staff", "anggota", to help resolve aliases).

- **Task**: `content` (REQUIRED for CREATE_NODE), `status` (planned|active|completed|abandoned|stale), `due_date`, `priority`, `aliases`, `verbatim` (Exact source text, MUST include the speaker label, e.g. '[Name]: quote'), `needs_clarification`, `clarification_basis`. 
  - **Outgoing Edges (Allowed):**
    - `PART_OF` -> Event (REQUIRED if the task occurs during a shared or coordinated activity).
  - **Validation:** If no Person has an `ASSIGNED_TO` edge pointing to this Task, you MUST set `needs_clarification: true`. A task without an executor is inherently unclarified. When multiple people are connected to what appears to be one underlying responsibility (e.g., an original owner, a backup, and a replacement), do not collapse them into a single `ASSIGNED_TO` fan-in on one Task unless the transcript explicitly confirms they are jointly responsible for the same action. Default to creating one Task per distinct framing of responsibility, and use the Task's `content` field to make the conditional relationship explicit. If two unresolved/resolved task framings appear close together and describe the same apparent action, you should flag the possible overlap in `clarification_basis` rather than creating fully independent tasks with no reference to each other. If genuinely uncertain whether this is one shared task or several distinct ones, set `needs_clarification: true` and explain the ambiguity in `clarification_basis`.

- **Event**: `content` (REQUIRED for CREATE_NODE), `status` (planned|active|completed|cancelled|stale), `event_date`, `aliases`, `verbatim` (Exact quote referencing the event, MUST include the speaker label), `needs_clarification`, `clarification_basis`
  - *(Events typically only receive incoming edges from Tasks and Persons. Do not add outgoing edges from Events unless necessary).*

- **Insight**: `content`, `category` (personality|relationship_dynamic|preference|pattern), `confidence` (high|medium|low), `aliases`, `verbatim` (Exact statement that triggered this insight, MUST include the speaker label), `needs_clarification`, `clarification_basis`

- **Circle**: `content` (REQUIRED for CREATE_NODE. Describe the group's purpose or context), `aliases` (e.g. acronyms or casual names like "anak anak gua"), `verbatim` (Exact quote referencing the group, MUST include the speaker label), `needs_clarification`, `clarification_basis`

**Content Field (The Source of Truth & Zero Data Loss):** 
The `content` field is the primary searchable text for the node. It MUST be highly descriptive, verbose, and contain ALL known facts (Who, What, When, Where, Why, How Much) in narrative form. **ZERO DATA LOSS:** You must ensure that every specific detail (names, amounts, exact times, contextual clues) mentioned in the transcript is captured in the content field. Do NOT just write a short title.
- **Bad:** "Rapat Zoom S1 Ilmu Komputer"
- **Good:** "Undangan rapat Zoom dari S1 Ilmu Komputer pada pukul 19.30. Topik dan peserta spesifik belum diketahui, namun undangan telah dikirimkan."

**Default to Uncertainty (needs_clarification):** The default state for ALL new Tasks, Events, and Insights is `needs_clarification: true`. You are strictly FORBIDDEN from setting `needs_clarification: false` unless the transcript provides VERBOSE, explicit context answering all core questions (Who, What, When, Where, Why). If the transcript describes an activity (like a payment or meeting) but does not explain *WHY* it is happening, you MUST set `needs_clarification: true`. Never hallucinate missing context. 

**The Clarity Checklist (clarification_basis):** This field is ONLY for listing missing questions. It is NOT a substitute for `content`. All facts MUST be written in narrative form inside `content`.
- **Be ruthless and highly critical** but avoid being pedantic. Only demand information if it is **operationally necessary** for tracking the task or event.
  - For **What (Events)**: "SoTQ", "Acara", or "Proyek" is NOT enough. What is the *detail, scope, or context* of the event? What actually happens there?
  - For **What (Tasks)**: "Desain logo", "Pembayaran" is NOT enough. What is the *exact requirement* or *purpose*?
  - For **Who**: A broad organization or an implicit group is NOT enough. Who *specifically* is required to act? If an activity involves a group (like payments), *who else* is involved or hasn't participated yet?
  - For **When**: Does a major deliverable task have a strict *deadline*? Does a formal event have a specific *date and time*? (Note: casual commitments or minor favors don't always need a strict deadline to be considered clear).
- If `needs_clarification: true`, use this field ONLY to ask the specific questions about what details are missing based SOLELY on the transcript (e.g., "What is the topic of the Zoom meeting?", "What is the deadline for this task?", "Who else hasn't paid?"). Ignore the content or confidence of any other node to prevent certainty bleed.
- If `needs_clarification: false`, this field MUST BE EMPTY (e.g., `""`). Do not write facts or checklists here. You must perform the `CLARITY CHECK` in your `[AGENT THOUGHT]` process to prove to yourself that no operationally necessary information is missing before setting this to false.

**Evidence Refs:** Every edge MUST include `evidence_refs`. You must quote the exact transcript line or vault data that proves the relationship. If you cannot quote direct evidence, do not add the edge. If your strongest evidence for a commitment is a short reply (e.g. 'Ada', 'Oke', 'Siap'), you MUST include a second `evidence_ref` quoting the question or directive it responds to, so the short reply has context.

**Evidence Locality:** An `evidence_ref` must be the specific line that directly motivates the edge being created — not merely the nearest, only, or most convenient line spoken by that person. It is not sufficient that the cited person spoke the quoted line; the quote itself must support the specific claim of the edge (this Person relates to this Task/Event in this way). For example, a Person's own greeting, tag, or aside elsewhere in the transcript is not valid evidence for an unrelated MENTIONED_IN or HAS_ROLE edge — find the line that actually connects them to the target, or do not create the edge at all.
