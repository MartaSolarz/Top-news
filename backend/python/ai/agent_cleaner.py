import requests


def clean_text(text, url, openai_api_key, model="text-davinci-003", temperature=0.3, max_tokens=1024):
    headers = {
        "Authorization": f"Bearer {openai_api_key}"
    }
    data = {
        "model": model,
        "prompt": f"You are a Text Editor and you have to clean the following text:\n\n{text}\n\n."
                  f"Remove any unexpected characters, words, or phrases that do not belong in the text.\n"
                  f"Make sure the text does not contain any production information "
                  f"which could be added when we using request library to obtain the data.\n",
        "temperature": temperature,
        "max_tokens": max_tokens
    }

    response = requests.post(url, headers=headers, json=data)
    cleaned_text = response.json().get("choices")[0].get("text").strip()

    return cleaned_text
