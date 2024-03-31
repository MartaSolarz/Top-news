# Top News
## BBC Article Viewer with AI Summary and Subscription - portfolio project

![Image](https://upload.wikimedia.org/wikipedia/commons/thumb/a/ad/BBC_News_logo.svg/2560px-BBC_News_logo.svg.png)

**Source:** https://upload.wikimedia.org/wikipedia/commons/thumb/a/ad/BBC_News_logo.svg/2560px-BBC_News_logo.svg.png

## Overview

This application is a comprehensive solution for browsing, filtering, and managing articles from BBC. 
It features a home page displaying articles, with capabilities for filtering and searching, marking articles as favorites, and a subscription page for receiving email updates. 
The application is powered by a combination of frontend technologies (HTML, CSS, JavaScript) and backend services (Go, Python), utilizing artificial intelligence for generating article summaries.
For storing data, the application uses a MySQL database.

## Apps

In this project, we have three main apps:
- **server/main.go**: The main server application that serves the whole website.
- **rssfetcher/main.go**: A service that fetches articles from the BBC RSS feed.
- **mailing/main.py**: A service that sends emails to subscribers with the latest articles.

## Features

- **Article Browsing**: View a list of articles fetched one per day from the BBC RSS feed.
- **Search and Filter**: Easily search for articles or filter them based on categories or tags.
- **Favorites**: Mark articles as favorites for quick access in the future. View all your favorite articles in one place and manage them as needed.
- **Subscription**: Subscribe with your email to receive regular updates and newsletters containing the latest articles and summaries.
- **AI-Powered Summaries**: Leverage artificial intelligence to generate concise summaries of articles, making it easier to consume content.

### Usage

- Browse articles on the home page.
- Use the search bar or filter options to find specific articles.
- Click the heart icon to add articles from your favorites.
- Visit the Favorites page to view all articles you've marked as favorites or remove it from your favorites.
- Go to the Subscription page to subscribe with your email for updates.


## License

This project is licensed under the Apache License - see the LICENSE.md file for details.
