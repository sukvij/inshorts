import spacy
from flask import Flask, request, jsonify
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.linear_model import LogisticRegression
import numpy as np # Import numpy for array conversion

app = Flask(__name__)

# --- 1. Load spaCy model for NER ---
# You might need to download this model first:
# python -m spacy download en_core_web_sm
try:
    nlp = spacy.load("en_core_web_sm")
except OSError:
    print("SpaCy model 'en_core_web_sm' not found. Please run: python -m spacy download en_core_web_sm")
    exit()

# --- 2. Train a very basic Intent Recognition model ---
# In a real application, you'd have a much larger, diverse dataset
# and potentially a more sophisticated model (e.g., BERT, fine-tuned classifier).
train_queries = [
    "Latest developments in the Elon Musk Twitter acquisition near Palo Alto",
    "News about tech companies in California",
    "Find restaurants near me",
    "What's the weather like in Mumbai",
    "Who is the CEO of Google?",
    "Recent updates on Tesla stock",
    "Events happening around Delhi",
    "How to code in Python",
    "Acquisition news for software companies",
    "Companies based in Silicon Valley",
    "Is there a park close by?",
    "Show me results for sports news",
    "What's the status of the Space X launch?"
]

train_intents = [
    "nearby_developments", # This is our target intent for the example query
    "general_news",
    "location_search",
    "weather_query",
    "person_info",
    "stock_news",
    "location_search",
    "how_to",
    "general_news",
    "location_search",
    "location_search",
    "general_news",
    "company_updates"
]

# Create TF-IDF features
vectorizer = TfidfVectorizer()
X_train_vec = vectorizer.fit_transform(train_queries)

# Train a Logistic Regression classifier
intent_classifier = LogisticRegression(max_iter=1000)
intent_classifier.fit(X_train_vec, train_intents)

# --- 3. Define a simple Concept Mapping (Rule-based for this demo) ---
# In a real system, this could be:
# - A knowledge graph lookup
# - Topic modeling (e.g., LDA, NMF)
# - Semantic similarity to predefined concepts
concept_keywords = {
    "acquisition": ["acquisition", "buyout", "merger", "takeover"],
    "technology": ["tech", "software", "innovation", "digital", "AI", "ML", "robotics", "Tesla", "SpaceX", "Google", "Twitter"],
    "finance": ["stock", "market", "economy", "investment"],
    "location": ["nearby", "near me", "around", "in", "location", "area"]
}

def extract_concepts(query_text):
    found_concepts = set()
    query_lower = query_text.lower()
    for concept, keywords in concept_keywords.items():
        for keyword in keywords:
            if keyword in query_lower:
                found_concepts.add(concept)
    return list(found_concepts)

# --- API Endpoint ---
@app.route('/predict', methods=['POST'])
def predict():
    data = request.get_json()
    query = data.get('query', '')

    if not query:
        return jsonify({"error": "No query provided"}), 400

    # 1. Entity Recognition (NER)
    doc = nlp(query)
    entities = []
    for ent in doc.ents:
        entity_type = ""
        # Map spaCy's diverse labels to your desired general types
        if ent.label_ in ["PERSON", "ORG", "GPE", "LOC", "FAC"]:
            entity_type = ent.label_
            if ent.label_ == "GPE": # Geo-Political Entity, often locations
                entity_type = "LOCATION"
            elif ent.label_ == "ORG":
                entity_type = "ORGANIZATION"
            elif ent.label_ == "PERSON":
                entity_type = "PERSON"
            elif ent.label_ == "FAC": # Building, airport, highway, etc. often relevant for locations
                entity_type = "LOCATION"
            elif ent.label_ == "LOC": # Non-GPE locations
                entity_type = "LOCATION"
            # Add more specific mappings if needed

        if entity_type: # Only add if we have a mapped type
            entities.append({
                "text": ent.text,
                "type": entity_type
            })

    # 2. Concept Extraction
    concepts = extract_concepts(query)

    # 3. Intent Recognition
    query_vec = vectorizer.transform([query])
    intent_prediction = intent_classifier.predict(query_vec)[0]

    # Map the predicted intent to your desired output ("nearby" for the example)
    # This is a rule based mapping based on your training data's intent labels
    if intent_prediction == "nearby_developments":
        final_intent = "nearby"
    elif "location_search" in intent_prediction: # if your classifier can predict this category
         final_intent = "nearby"
    else:
        final_intent = "general_query" # Default or other general intent

    # Refine intent for the specific query example
    if "near Palo Alto" in query and intent_prediction == "nearby_developments":
        final_intent = "nearby" # Explicitly set for the example

    response_data = {
        "entities": entities,
        "concepts": concepts,
        "intent": final_intent
    }

    return jsonify(response_data)

if __name__ == '__main__':
    print("ML Service: Training basic intent model and loading SpaCy...")
    print("ML Service ready. Listening on http://127.0.0.1:5000")
    app.run(host='127.0.0.1', port=5000)