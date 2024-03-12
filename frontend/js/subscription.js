document.addEventListener('DOMContentLoaded', function() {
    bindMenuToggle();
    bindThemeToggle();
});

document.getElementById('subscriptionForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const email = document.getElementById('email').value;

    fetch('/api/save_email', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email: email })
    })
        .then (response => {
            if (response.ok) {
                alert("Thank you for subscribing!");
            } else if (response.status === 400) {
                alert("Email is invalid. Please enter a valid email address.");
            } else {
                console.error('Error:', response.status);
                alert("There was an error with your subscription. Please try again later.");
                throw new Error('Subscription failed with status: ' + response.status);
            }
        })
});

function bindMenuToggle() {
    document.querySelector(".menu-icon").addEventListener("click", openNav);
    document.querySelector(".close-btn").addEventListener("click", closeNav);
}

function bindThemeToggle() {
    const sunIcon = document.getElementById('sun');
    const moonIcon = document.getElementById('moon');

    if (sunIcon && moonIcon) {
        sunIcon.addEventListener('click', () => toggleTheme('light'));
        moonIcon.addEventListener('click', () => toggleTheme('dark'));
    }
}

function toggleTheme(selectedTheme) {
    const body = document.body;
    const sunIcon = document.getElementById('sun');
    const moonIcon = document.getElementById('moon');

    body.classList.toggle('dark-mode', selectedTheme === 'dark');
    sunIcon.classList.toggle('selected', selectedTheme === 'light');
    moonIcon.classList.toggle('selected', selectedTheme === 'dark');
}

function openNav() {
    document.getElementById("sideMenu").style.width = "200px";
    document.getElementById("main-content").style.marginLeft = "200px";
}

function closeNav() {
    document.getElementById("sideMenu").style.width = "0";
    document.getElementById("main-content").style.marginLeft = "0";
}
