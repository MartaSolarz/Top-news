// js/script.js
document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.show-more').forEach(item => {
        item.addEventListener('click', event => {
            let content = item.nextElementSibling;
            if (content.style.display === "none") {
                content.style.display = "block";
                item.textContent = "Show Less";
            } else {
                content.style.display = "none";
                item.textContent = "Show More";
            }
        });
    });
});

function searchArticles() {
    let input = document.getElementById('searchInput');
    let dateInput = document.getElementById('dateInput');
    let filter = input.value.toUpperCase();
    let dateFilter = dateInput.value ? new Date(dateInput.value).toISOString().slice(0, 10) : "";
    let articles = document.getElementsByClassName('article');

    for (let i = 0; i < articles.length; i++) {
        let title = articles[i].getElementsByClassName('title')[0];
        let articleDate = articles[i].getAttribute('data-date').slice(0, 10);

        if (title.innerHTML.toUpperCase().indexOf(filter) > -1 && (articleDate === dateFilter || dateFilter === "")) {
            articles[i].style.display = "";
        } else {
            articles[i].style.display = "none";
        }
    }
}

function resetDateFilter() {
    document.getElementById('dateInput').value = "";
    searchArticles();
}


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

function toggleLike(articleId, element) {
    element.classList.toggle('liked');

    let favorites = JSON.parse(localStorage.getItem('favorites')) || [];

    if (element.classList.contains('liked')) {
        if (!favorites.includes(articleId)) {
            favorites.push(articleId);
        }
        console.log("Article " + articleId + " added to favorites.");
    } else {
        favorites = favorites.filter(id => id !== articleId);
        console.log("Article " + articleId + " removed from favorites.");
    }

    localStorage.setItem('favorites', JSON.stringify(favorites));
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