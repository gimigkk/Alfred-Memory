# AI Skill: The Courtroom Simulation (V32)

## Prior Rulings
* **V29-V31:** Implemented `query_speaker_obligations` as a hard gate. Mandated `UPDATE_NODE` over `CREATE_NODE` for temporal clarifications (Rule 18).

## The Proposal on Trial
**Issue:** The agent correctly calls `query_speaker_obligations`, receives 3 unclarified nodes (2 Tasks, 1 Event), and successfully issues `UPDATE_NODE` for the 2 Tasks. However, it completely ignores the existing Event and issues a `CREATE_NODE` for a new Event (`temp_event_rapat_bph_gacoan`), resulting in duplicate Events in the graph.

**Proposal:** How do we structurally enforce that the agent updates the existing Event retrieved by `query_speaker_obligations` instead of creating a new one?

## The Personas
1. **The Software Architect (Pro-Structure):** Believes the graph schema must dictate agent behavior.
2. **The Prompt Engineer (Pro-Context):** Believes the LLM just needs clearer instructions in Rule 18.
3. **The DevOps/SRE (Pro-Observability):** Wants to look at what the agent actually saw.
4. **The Security Auditor (Pro-Validation):** Wants hard Go-side gates to reject the prompt if it duplicates nodes.
5. **The Hostile Attacker (Prosecutor):** Believes the current prompt is a semantic loophole that gives the agent an excuse to create new Events. Will exploit any ambiguity.
6. **The Judge:** Delivers the final verdict.

---

## Cycle 1: Identifying the Semantic Loophole

**The DevOps/SRE:**

Let's look at the agent's thought process in the failing trace. It wrote: "Event Identification: The transcript discusses a 'rapat bph gacoan' and the subsequent reimbursement process. This is a shared coordinated activity. I will create a new Event 'Rapat BPH Gacoan' to house these tasks." It didn't even mention the old Event in its reasoning!

**The Prompt Engineer:**

Rule 18 says: "If the existing node and the new context describe the same underlying activity, you must UPDATE the existing node." The agent probably thought "Rapat BPH Gacoan" was a *different* activity than "Koordinasi pengumpulan dana", so it justified creating a new one.

**The Hostile Attacker (Prosecutor):**

Of course it did! Look at Rule 14 (Event Inference): "If the transcript describes activity clearly occurring... you MUST create a new Event node to contain the Tasks discussed." You gave the agent two contradictory commands! Rule 14 screams "CREATE A NEW EVENT!", while Rule 18 softly suggests "update if it's the same activity." The LLM is an autocomplete engine; when it sees a shiny new noun like 'rapat bph gacoan', it triggers the Rule 14 `CREATE_NODE` pathway because the old event's generic name doesn't semantically match the new noun.

**The Software Architect:**

The Attacker is right. The agent updated the Tasks because they explicitly matched the exact bank accounts dropped earlier. But the Event had a generic name ("Koordinasi pengumpulan dana"). When the new transcript introduced a highly specific name ("Rapat BPH Gacoan"), the semantic distance was too wide for the LLM to bridge without explicit structural help.

**Outcome:**
**RULING CHANGED:** We acknowledge that Rule 14 (`CREATE_NODE`) overpowers Rule 18 (`UPDATE_NODE`) when the semantic naming of the event shifts drastically. We need a structural bridge.

---

## Cycle 2: Structural Bridging vs Prompt Tweaking

**The Prompt Engineer:**

We can fix this by adding a line to Rule 18: "If you update ANY Task that is part of an existing Event, you MUST also update that parent Event instead of creating a new one." 

**The Software Architect:**

That's still relying on the LLM's working memory. But wait... how does the LLM even know what Event the existing Task belongs to? In the output of `query_speaker_obligations`, we return: `task_apta_payment_44f8dc (Task): Apta Adi Nur Fiansah menyediakan...`. We do NOT tell the agent that `task_apta_payment` is connected to `event_coordination`! 

**The Hostile Attacker (Prosecutor):**

Exactly! You are expecting the agent to magically deduce that the old Tasks and the old Event belong together. The agent sees a list of 3 disconnected nodes. It updates the Tasks because they match. It sees the old Event, thinks "that's a generic payment event, I'm dealing with a Gacoan Reimbursement," and creates a new one. It doesn't know the Tasks it just updated were previously nested inside that old Event! 

**The DevOps/SRE:**

The Attacker nailed it. The output of `query_speaker_obligations` is just a flat list of nodes:
`task_apta_payment (Task): ...`
`event_coordination (Event): ...`
There is zero edge context showing `task_apta` -> `PART_OF` -> `event_coordination`.

**Outcome:**
**RULING CHANGED:** The root cause is data opaqueness. The agent cannot update the parent Event because `query_speaker_obligations` hides the edge topology connecting the returned obligations.

---

## Cycle 3: Designing the Graph Context

**The Software Architect:**

If `query_speaker_obligations` returns the graph topology, the agent will naturally update the whole cluster. We need to modify `tool_handlers.go` so that `query_speaker_obligations` returns the edges between the unclarified nodes it finds, or formats the output hierarchically.

**The Security Auditor:**

Wait, `query_speaker_obligations` currently returns:
`RETURN t.id, label(t), t.content, t.clarification_basis, label(e)`
It tells us the edge from the *Speaker* to the *Task* (`ASSIGNED_TO`), but it does NOT traverse the `PART_OF` edge from the *Task* to the *Event*.

**The Hostile Attacker (Prosecutor):**

Because your Cypher query is `MATCH (p)-[e]-(t) WHERE p.id IN ... AND t.needs_clarification = true`. This only finds nodes one hop away from the Speaker! The Event is likely two hops away (Speaker -> Task -> Event), or if the Speaker is directly `MENTIONED_IN` the Event, it's found, but the relationship *between* the Task and the Event is never queried or returned. You are feeding the LLM an incomplete graph and punishing it for not seeing the invisible lines!

**The DevOps/SRE:**

If we just return the full subgraph connected to those speakers via `query_rag`, it would see the topology. But `query_rag` doesn't filter for `needs_clarification=true`. 

**The Prompt Engineer:**

We don't need to change the Cypher query drastically. If the agent updates `task_apta_payment`, and it wants to link it to the new `temp_event_rapat_bph_gacoan`, it will emit an `add_edges` command: `PART_OF -> temp_event_rapat_bph_gacoan`. 
But wait... if it does that, the `task_apta_payment` will be `PART_OF` *two* events in the database! The old one and the new one. 

**The Hostile Attacker (Prosecutor):**

Which violates your graph schema! A Task should only be `PART_OF` one Event. If you let it create a new Event and link the old Task to it, you've created a corrupted graph where a Task belongs to two mutually exclusive realities.

**Outcome:**
**RULING CHANGED:** The Go Orchestrator must prevent Tasks from being linked to multiple Events, OR the prompt must strictly forbid creating a new Event if updating a Task.

---

## Cycle 4: The Final Solution

**The Software Architect:**

We cannot rely on the LLM to resolve invisible graph edges. We must explicitly tell the LLM about the `PART_OF` relationships during `query_speaker_obligations`. 

**The Security Auditor:**

Let's look at `tool_handlers.go`. 
```go
	query := fmt.Sprintf(`
		MATCH (p)-[e]-(t)
		WHERE p.id IN %s
		  AND t.needs_clarification = true
		RETURN t.id, label(t), t.content, t.clarification_basis, label(e)
	`, idsList)
```
Instead of just returning this, what if we also query for the parent Events of those Tasks?

**The Hostile Attacker (Prosecutor):**

Too much engineering. Just add a rule in `skill_commit.md`: 
"**Rule 19: Event Hijacking Prohibited:** If you use `UPDATE_NODE` to update a Task that was returned by `query_speaker_obligations`, you are FORBIDDEN from creating a new Event (`CREATE_NODE` Event) in the same commit. You MUST assume the Task already belongs to an existing Event. If an Event was also returned by `query_speaker_obligations`, you MUST update that Event instead of creating a new one."

**The Prompt Engineer:**

That is extremely brittle. What if the new transcript introduces a genuinely new Event? 

**The DevOps/SRE:**

Let's look at the actual output of the trace:
`event_coordination_44f4d1 (Event): Koordinasi pengumpulan dana antar anggota IEEE25. Apta dan Jeslyn membagikan det...`
The agent *did* see the event! It just chose to create a new one because "Rapat BPH Gacoan" sounded like a better name, and Rule 18 doesn't explicitly forbid replacing old nodes with better-named new ones.

**The Software Architect:**

So we just update Rule 18. "If `query_speaker_obligations` returns an existing Event, and the new transcript reveals the specific name or purpose of that Event (e.g., 'Rapat BPH Gacoan'), you MUST use `UPDATE_NODE` on the existing Event to update its `content` and clear its `clarification_basis`. You are STRICTLY FORBIDDEN from creating a new Event node when an unclarified Event node already exists for the same speakers."

**The Hostile Attacker (Prosecutor):**

I accept this. If you explicitly ban `CREATE_NODE` for Events when an unclarified Event is in the obligations list, the LLM's only path to resolve the context is via `UPDATE_NODE`. You close the semantic loophole.

**Outcome:**
**RULING CHANGED:** We will modify Rule 18 in `skill_commit.md` to explicitly forbid creating a new Event if an unclarified Event is returned by the obligations query, forcing the agent to update the existing Event with the newly discovered specific details.

---

## The Verdict

**Findings:**
1. The LLM hallucinates a new Event because Rule 14 (`CREATE_NODE` for specific events) overpowers Rule 18 (`UPDATE_NODE` for general semantic matches).
2. When the LLM learns a specific noun ("Rapat BPH Gacoan"), it abandons the old generic node ("Koordinasi pengumpulan") because the semantic distance is too wide without explicit guardrails.
3. The LLM does not know the old Tasks are connected to the old Event, so it feels free to link the old Tasks to a shiny new Event.

**Directives:**
- Modify `assets/prompts/skill_commit.md` Rule 18 to explicitly state that if an unclarified Event is returned by `query_speaker_obligations`, and the new transcript provides the *real name or specific purpose* of that event, the agent is **STRICTLY FORBIDDEN** from creating a new Event. It **MUST** update the existing Event.
- Explicitly state that old Tasks cannot be "hijacked" into new Events.
