import sys
import requests

from bs4 import BeautifulSoup


def extract_content(url):
    response = requests.get(url)

    if response.status_code == 200:
        soup = BeautifulSoup(response.text, 'html.parser')
        article = soup.find('article')
        if article:
            return article.get_text(strip=True)
        else:
            return "Article content could not be found."
    else:
        return "Failed to retrieve the URL content."


if __name__ == "__main__":
    url = sys.argv[1]
    content = extract_content(url)
    sys.stdout.write(content)
