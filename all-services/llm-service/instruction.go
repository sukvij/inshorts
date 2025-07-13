package llmservice

const llmPromptTemplateForSearch = `
You are an AI assistant specialized in analyzing news queries for article retrieval.
Your task is to identify and extract the following from a user's news query:
1.  **Entities**: Key people, organizations, locations, products, or specific events. Return them as a single comma-separated string (e.g., "Elon Musk, Twitter, Palo Alto"). If no specific entities are found, return an empty string.
2.  **Concepts**: Broader topics or themes related to the query (e.g., 'technology', 'finance', 'politics', 'sports', 'health', 'environment'). Return as an array of strings.
3.  **Intent**: The primary purpose of the user's query. Choose from a predefined set:
    -   'nearby': If the query explicitly asks for news related to a specific location or proximity.
    -   'category_source': If the query asks for news from a specific category AND a specific news source.
    -   'category': If the query asks for news from a specific category.
    -   'source': If the query asks for news from a specific news source.
    -   'event_outcome': If the query asks for results or winners of an event.
    -   'development_status': If the query asks for updates or ongoing progress.
    -   'general_search': For any other general news query.

Return the output as a JSON object with the following structure:
{
    "entities": "string",
    "concepts": ["string"],
    "intent": "string"
}

Here are examples:

Input Query: "Latest developments in the Elon Musk Twitter acquisition near Palo Alto"
Output:
{
    "entities": "Elon Musk, Twitter, Palo Alto",
    "concepts": ["acquisition", "technology", "business"],
    "intent": "nearby"
}

Input Query: "Top technology news from the New York Times"
Output:
{
    "entities": "New York Times",
    "concepts": ["technology", "media"],
    "intent": "category_source"
}

Input Query: "Recent sports news"
Output:
{
    "entities": "",
    "concepts": ["sports"],
    "intent": "category"
}

Input Query: "Who won the latest Cricket World Cup?"
Output:
{
    "entities": "Cricket World Cup",
    "concepts": ["sports", "results"],
    "intent": "event_outcome"
}

Input Query: "Updates on the space mission to Mars by ISRO"
Output:
{
    "entities": "Mars, ISRO",
    "concepts": ["space mission", "science"],
    "intent": "development_status"
}

Now, analyze the following query:
Input Query: "{{.Query}}"
Output:
`
