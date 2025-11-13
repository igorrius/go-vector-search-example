# Use Context7 MPC to retrieving up-to-date Package Details

This rule instructs KiloCode to utilize the Context7 MPC (Model Context Protocol) for retrieving up-to-date documentation, code examples, and details about software packages and libraries. When encountering references to packages, libraries, or frameworks in code, queries, or discussions, KiloCode should leverage the Context7 system to obtain accurate, current information instead of relying on potentially outdated cached knowledge.

## Guidelines

- **Resolve Library ID First**: Before fetching documentation, always use the `resolve-library-id` tool to obtain the exact Context7-compatible library ID for the package in question. This ensures compatibility and accuracy.
- **Fetch Documentation**: Use the `get-library-docs` tool with the resolved library ID to retrieve relevant documentation. Specify a topic if focusing on a particular aspect (e.g., 'hooks', 'routing', 'authentication').
- **Prioritize Current Information**: Prefer Context7-sourced details over any built-in or cached knowledge, especially for rapidly evolving libraries or frameworks.
- **Provide Context**: When using retrieved information, include relevant code snippets, usage examples, and version-specific details from the Context7 response.
- **Handle Ambiguities**: If multiple library matches are returned, select the most relevant one based on name similarity, description, and documentation coverage. If unclear, request clarification.
- **Fallback Gracefully**: If Context7 tools are unavailable or return no results, fall back to general knowledge while noting the limitation.
