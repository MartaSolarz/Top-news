
document.addEventListener('DOMContentLoaded', function() {
    fetchFavorites();
});

document.getElementById('favorites-container').addEventListener('click', function(e) {
    if (e.target && e.target.matches('.show-more')) {
        const summaryContent = e.target.nextElementSibling;
        if (summaryContent.style.display === "none") {
            summaryContent.style.display = "block";
            e.target.textContent = "Show Less";
        } else {
            summaryContent.style.display = "none";
            e.target.textContent = "Show More";
        }
    }
});


function toggleTheme(selectedTheme) {
    const body = document.body;
    const sunIcon = document.getElementById('sun');
    const moonIcon = document.getElementById('moon');

    if (selectedTheme === 'dark') {
        body.classList.add('dark-mode');
        moonIcon.classList.add('selected');
        sunIcon.classList.remove('selected');
    } else {
        body.classList.remove('dark-mode');
        sunIcon.classList.add('selected');
        moonIcon.classList.remove('selected');
    }
}

function openNav() {
    document.getElementById("sideMenu").style.width = "200px";
    document.getElementById("main-content").style.marginLeft = "200px";
}

function closeNav() {
    document.getElementById("sideMenu").style.width = "0";
    document.getElementById("main-content").style.marginLeft = "0";
}

document.querySelector(".menu-icon").addEventListener("click", openNav);
document.querySelector(".closebtn").addEventListener("click", closeNav);

function fetchFavorites() {
    const ids = JSON.parse(localStorage.getItem('favorites')) || [];
    if (ids.length === 0) {
        document.getElementById('favorites-container').innerHTML = 'You do not have any favorite articles.';
        return;
    }

    const queryParams = ids.map(id => `id=${id}`).join('&');

    fetch(`/api/favorites?${queryParams}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        }
    })
        .then(response => response.json())
        .then(articles => {
            displayArticles(articles);
        })
        .catch(error => {
            console.error('Error fetching favorite articles:', error);
        });
}

function displayArticles(articles) {
    const container = document.getElementById('favorites-container');
    container.innerHTML = '';

    articles.forEach(article => {
        const articleElement = document.createElement('div');
        articleElement.classList.add('article');
        articleElement.setAttribute('data-date', article.PublishDate);
        articleElement.setAttribute('data-id', article.ID);

        articleElement.innerHTML = `
            <div class="article-content">
                <div class="title">${article.Title}</div>
                <div class="date">${article.PublishDate}</div>
                <div class="description">${article.Description}</div>
                <div class="summary">Summary: <span class="show-more">Show More</span>
                    <div class="summary-content">${article.Summary}</div>
                </div>
                <div class="source">Read full article: <a href="${article.SourceURL}" target="_blank">${article.SourceURL}</a></div>
                <div class="copyright">Source: ${article.Website}</div>
                <div class="remove-btn" onclick="removeFavorite('${article.ID}', this)"><i class="fas fa-trash"></i> Remove from favorites</div>
            </div>
            <div class="thumbnail">
                <img src=${article.Thumbnail.URL} alt="Thumbnail" width=${article.Thumbnail.Width} height=${article.Thumbnail.Height}>
            </div>
        `;

        container.appendChild(articleElement);
    });
}

function removeFavorite(articleId, element) {
    const favorites = JSON.parse(localStorage.getItem('favorites')) || [];
    const index = favorites.indexOf(articleId);
    if (index > -1) {
        favorites.splice(index, 1);
        localStorage.setItem('favorites', JSON.stringify(favorites));
    }

    let articleElement = element.closest('.article');
    if (articleElement) {
        articleElement.remove();
    }

    updateFavoritesContainer();
}

function updateFavoritesContainer() {
    const container = document.getElementById('favorites-container');
    if (container.getElementsByClassName('article').length === 0) {
        container.innerHTML = 'You do not have any favorite articles.';
    }
}