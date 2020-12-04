'use strict';

async function get_profiles(url) {
    let response = await fetch(url);
    return response.json();
}

async function main() {
    let url = 'https://api.github.com/users';
    let profiles = await get_profiles(url);
    let footer = document.createElement("footer");
    for (let profile of profiles) {
        let greeting = document.createElement('p');
        greeting.textContent = `greeting ${profile.login}, here is footer!`;
        footer.appendChild(greeting);
    }
    document.querySelector("#root").appendChild(footer);
}

main();