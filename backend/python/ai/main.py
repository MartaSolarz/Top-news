import sys
import json

from agent_cleaner import clean_text
from agent_summarizer import summarize

if __name__ == "__main__":
    input_text, open_api_url, open_api_key = sys.stdin.read()

    cleaned_text = clean_text(input_text, open_api_url, open_api_key)

    summary = summarize(cleaned_text, open_api_url, open_api_key)

    print(json.dumps({"summary": summary}))
