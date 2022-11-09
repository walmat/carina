import { METADATA } from "../constants";

async function fetchGraphQL(
  text: string | null | undefined,
  variables: Record<string, unknown>
) {
  // Fetch data from GitHub's GraphQL API:
  const response = await fetch(METADATA.GRAPHQL_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      // handle authorization
    },
    body: JSON.stringify({
      query: text,
      variables,
    }),
  });

  // Get the response as JSON
  return await response.json();
}

export default fetchGraphQL;
