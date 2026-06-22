# AI Skill: The Courtroom Simulation

When the user requests a "courtroom" or "courtroom simulation" to debate a technical or architectural decision, execute the simulation strictly according to the following constraints.

## 1. The Setup
- Identify the exact Proposal on Trial.
- Establish a minimum of 5 Personas representing different engineering perspectives (e.g., The Software Architect, The Prompt Engineer, The DevOps/SRE, The Performance Engineer, The Security Auditor, The Judge).
- Explicitly state each persona's inherent bias or priority before the debate begins.

## 2. The Execution (Cycles)
- The debate MUST run for a minimum of 3 full cycles.
- Each cycle consists of the personas directly attacking or rebutting the arguments made by the opposing personas in the previous cycle.
- Each persona's argument MUST be at least one substantial paragraph long. Do not use brief, surface-level agreements or one-liners.

## 3. The Evidence Rule
- EVERY argument made by a persona MUST be backed by researched evidence.
- You must use web search tools to find real, external validation (academic papers, official documentation, industry articles, or cloud economic data) to support the claims.
- Cite the evidence inline and provide a references section at the bottom.

## 4. The Verdict (The Judge's Constraints)
- The Judge persona synthesizes the arguments and the empirical data.
- **ANTI-FALLACY RULE:** The Judge MUST strictly evaluate sample sizes and statistical significance. The Judge is strictly forbidden from committing Hasty Generalizations (e.g., weighing a single $n=1$ success over a historical $n=100$ failure rate).
- The Judge must prioritize documented systemic behavior over isolated anomalies.
- The Judge delivers a final, definitive verdict approving or denying the proposal based on sound logical and statistical reasoning.

## 5. Formatting Requirements
- Output the simulation as a dedicated Markdown Artifact.
- Use clear vertical spacing between paragraphs.
- Use horizontal rules (`---`) between debate cycles.
- Bold the names of the personas.
- Ensure the document is highly readable and not clustered.
