// js/script.js

// Event listener for DOMContentLoaded to initialize page functionalities
document.addEventListener('DOMContentLoaded', initializePageFunctions);

// Function to initialize all page functionalities
function initializePageFunctions() {
    setupShowMoreLess();
    checkLikedArticles();
    setupSearchAndReset();
    setupMenuToggle();
}

// Function to set up show more/less functionality
function setupShowMoreLess() {
    document.querySelectorAll('.show-more').forEach(item => {
        item.addEventListener('click', event => {
            let content = item.nextElementSibling;
            content.style.display = content.style.display === "none" ? "block" : "none";
            item.textContent = content.style.display === "block" ? "Show Less" : "Show More";
        });
    });
}

// Function to check and update liked articles on page load
function checkLikedArticles() {
    let favorites = JSON.parse(localStorage.getItem('favorites')) || [];
    document.querySelectorAll('.like-icon').forEach(element => {
        const articleId = element.getAttribute('data-article-id');
        if (favorites.includes(articleId)) {
            element.classList.add('liked');
        }
    });
}

// Function to set up article search and reset functionality
function setupSearchAndReset() {
    document.getElementById('searchInput').addEventListener('input', searchArticles);
    document.getElementById('dateInput').addEventListener('input', searchArticles);
    document.getElementById('resetDateFilter').addEventListener('click', resetDateFilter);
}

// Function to toggle theme between light and dark
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

// Function for toggling the like state of an article
function toggleLike(articleId, element) {
    element.classList.toggle('liked');
    let favorites = JSON.parse(localStorage.getItem('favorites')) || [];

    if (element.classList.contains('liked')) {
        if (!favorites.includes(articleId)) {
            favorites.push(articleId);
        }
    } else {
        favorites = favorites.filter(id => id !== articleId);
    }

    localStorage.setItem('favorites', JSON.stringify(favorites));
}

// Functions for search and reset
function searchArticles() {
    let input = document.getElementById('searchInput').value.toUpperCase();
    let dateInput = document.getElementById('dateInput').value;
    let dateFilter = dateInput ? new Date(dateInput).toISOString().slice(0, 10) : "";
    let articles = document.querySelectorAll('.article');

    articles.forEach(article => {
        let title = article.querySelector('.title').textContent.toUpperCase();
        let articleDate = article.getAttribute('data-date').slice(0, 10);

        article.style.display = (title.includes(input) && (articleDate === dateFilter || dateFilter === "")) ? "" : "none";
    });
}

function resetDateFilter() {
    document.getElementById('dateInput').value = "";
    searchArticles();
}

// Function to open and close the navigation menu
function setupMenuToggle() {
    document.querySelector(".menu-icon").addEventListener("click", () => {
        document.getElementById("sideMenu").style.width = "200px";
        document.getElementById("main-content").style.marginLeft = "200px";
    });

    document.querySelector(".close-btn").addEventListener("click", () => {
        document.getElementById("sideMenu").style.width = "0";
        document.getElementById("main-content").style.marginLeft = "0";
    });
}

document.querySelector(".menu-icon").addEventListener("click", openNav);
document.querySelector(".close-btn").addEventListener("click", closeNav);

// Handling SPA-like navigation changes
window.onpopstate = function(event) {
    initializePageFunctions();
};
