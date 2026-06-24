# Courtroom Simulation V34: Ontological Boundaries & Event Hubbing

**The Proposal on Trial:** The ingestion agent's tendency to create `task_mention` nodes for status updates and its failure to link participants directly to the `Event` node.

## Cycle 1: The Hallucination of the "Message Node"

**The Prosecutor (Hostile Attacker):**
The ingestion agent is hallucinating again. Look at this mutation: it successfully updated the `task_desain_logo` to "completed", but then it created a blank `task_mention` node just to represent Rendi's chat message! This is a massive ontological violation. A message is not a node.

**The Prompt Engineer:**
The agent is trying to preserve the temporal context of *how* it learned the task was completed. It thinks that if it doesn't create a node for the message, the user will lose the conversational context of Rendi saying "it's in the drive."

**The Judge:**
RULING CHANGED. The Prosecutor is correct. The prompt currently allows the agent to think a "message" is a valid entity. We must strictly enforce that the graph ontology only contains real-world operations (Tasks, Events, Projects, Insights, Circles).
**Directing the Prompt Engineer:** Add a strict rule: "You are STRICTLY FORBIDDEN from creating a node to represent a message, a chat, or a status update. If a message updates a task, use UPDATE_NODE on the task. DO NOT create a separate node for the message."

---

## Cycle 2: The Event Hub Topology

**The Prosecutor (Hostile Attacker):**
Okay, we fixed the message nodes. But look at the topology. The agent linked Rendi to the Task (`ASSIGNED_TO`). The Task is linked to the Event (`PART_OF`). But Rendi has NO direct link to the Event. The Event node is becoming a hollow ghost town with no people attached, while the Tasks act as the hubs!

**The Software Architect:**
That's technically a valid bipartite graph. You can traverse `Person -> Task -> Event` to find who was involved in the event. We don't necessarily need duplicate edges directly from `Person -> Event`.

**The Prosecutor (Hostile Attacker):**
Traversing two hops every time we want to know "who attended the event" is incredibly expensive for RAG retrieval windows! The LLM context window will be flooded with intermediate Task nodes just to answer a simple "Who" question about the Event. The Event must act as a highly connected hub!

**The Judge:**
RULING CHANGED. The Prosecutor is right regarding the RAG optimization. The Event node must be a hub of people, not just a hub of tasks.
**Directing the Prompt Engineer:** Update the rules to mandate Dual-Linking: "If a Person is linked to a Task that belongs to an Event, that Person MUST ALSO be linked directly to the parent Event via HAS_ROLE, PART_OF, or MENTIONED_IN."
