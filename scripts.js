const container = document.getElementById('projects-container');
const username = 'Accelertor'; // Replace 'your-username' with your actual GitHub username

fetch(`https://api.github.com/users/${username}/repos`)
    .then(response => response.json())
    .then(data => {
        data.forEach(repo => {
            const box = document.createElement('div');
            box.classList.add('project-box');

            // Create a link element for the repository
            const repoLink = document.createElement('a');
            repoLink.href = repo.html_url; // GitHub URL of the repository
            repoLink.target = '_blank'; // Open link in a new tab
            repoLink.rel = 'noopener noreferrer'; // Security best practices
            repoLink.style.textDecoration = 'none'; // Remove underline
            box.appendChild(repoLink);

            const title = document.createElement('h2');
            title.textContent = repo.name;
            repoLink.appendChild(title);

            const language = document.createElement('h3');
            language.textContent = repo.language;

            const readmeContent = document.createElement('p');
            readmeContent.textContent = "Loading...";
            fetch(`https://api.github.com/repos/${username}/${repo.name}/contents`)
                .then(response => response.json())
                .then(contents => {
                    // Searching for the image file starting with "icon"
                    const iconFile = contents.find(file => file.name.startsWith('icon'));
                    if (iconFile) {
                        // If found, create and append the image element
                        const img = document.createElement('img');
                        img.src = iconFile.download_url;
                        img.alt = repo.name;
                        box.appendChild(img);
                    } else {
                        console.log('No icon found for', repo.name);
                    }
                })
                .catch(error => console.error('Error fetching repository contents:', error));
            // Fetching README content
            fetch(`https://api.github.com/repos/${username}/${repo.name}/readme`)
                .then(response => response.json())
                .then(readme => {
                    const decodedReadme = atob(readme.content);
                    readmeContent.textContent = decodedReadme;
                })
                .catch(error => console.error('Error fetching README:', error));

            // Appending language and readme content
            const detailsContainer = document.createElement('div');
            detailsContainer.classList.add('project-details');
            detailsContainer.appendChild(language);
            detailsContainer.appendChild(readmeContent);
            box.appendChild(detailsContainer);

            // Appending the project box to the container
            container.appendChild(box);
        });
    })
    .catch(error => console.error('Error fetching GitHub projects:', error));
