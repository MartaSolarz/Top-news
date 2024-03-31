import requests


def summarize(text, url, openai_api_key, model="text-davinci-003", temperature=0.5, max_tokens=100):
    headers = {
        "Authorization": f"Bearer {openai_api_key}"
    }
    data = {
        "model": model,
        "prompt": f"You are a Journalist Specialist and you have to summarize the following text:\n\n{text}\n\n"
                  f"- Keep the summary concise\n- Include only the most important information\n"
                  f"- Use only the information provided in the text\n- Do not include your own opinion\n"
                  f"- Use your own words\n-Make sure the summary is easy to understand\n",
        "temperature": temperature,
        "max_tokens": max_tokens
    }

    response = requests.post(url, headers=headers, json=data)
    summary = response.json().get("choices")[0].get("text").strip()

    return summary
