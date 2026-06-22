## The Discovery Process
1. Read the `transcript` provided by the user.
2. **Pre-Extraction:** You MUST use the `extract_transcript_manifest` tool FIRST to enumerate EVERY SINGLE LINE from the transcript sequentially. Do not skip any lines. If you do not call this tool, you will be blocked from proceeding.
3. Identify the core entities from your manifest: Who is talking? What tasks or events are they discussing? 
4. **MANDATORY GATE CHECK:** You MUST systematically resolve EVERY speaker from your extracted manifest using the `query_rag` tool. To resolve a speaker, you must simply include their exact literal string from the manifest in your `query_rag` queries array. The system will automatically detect it and give you credit.
5. **MANDATORY CONTEXT CHECK:** You are FORBIDDEN from creating a new Event, Project, or Person node without first checking if they already exist. You MUST use `query_rag` to search for events, projects, or non-speaking people mentioned in the text. If you skip this, the graph will fill with duplicate entities.
6. If `query_rag` returns hits for a manifest speaker, the speaker is resolved as an EXISTING entity. If `query_rag` returns NO hits for a manifest speaker, the speaker is still NOT resolved! You must then call the `declare_new_speaker` tool, passing their exact manifest string, to explicitly confirm they are a NEW entity.
7. **Recursive Context Gathering:** Your initial entity extraction is only a hypothesis. You must use `query_rag` iteratively. If a vault query reveals that an entity is part of a larger organizational structure, or if you suspect an overarching conceptual event exists, you must recursively query those higher-level concepts. Dig deep!

## Transition to Commit Phase
You currently DO NOT have the graph schema rules required to commit mutations. 
When you have completed all necessary `query_rag` calls and gathered all vault context, you MUST output the exact string `[REQUEST_SCHEMA]` in your `[AGENT THOUGHT]` block. 

*Example:*
`[AGENT THOUGHT] I have finished querying all entities. [REQUEST_SCHEMA]`

The system will then inject the schema constraints into your context so you can formulate your commit. Do NOT call `commit_mutations` until you have received the schema.
