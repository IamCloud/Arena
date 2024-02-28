let Server = {
    initPlayer: function (playerName, teamName) {
        return new Promise((resolve, reject) => {
            fetch("/initplayer", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ PlayerName: playerName, Name: teamName, Wins: 0 })
            })
                .then(response => {
                    if (!response.ok) {
                        console.error('Network response was not ok');
                    }
                    resolve(true);
                })
                .then(data => {
                    resolve(data && data.success);
                })
                .catch(error => {
                    console.error('Error creating team:', error);
                    reject(error);
                })
        });
    },

    // Function to fetch leaderboard data from the Go server
    getLeaderboard: function () {
        fetch('/leaderboard')
            .then(response => response.json())
            .then(data => {
                // Call function to populate leaderboard with the retrieved data
                populateLeaderboard(data);
            })
            .catch(error => {
                console.error('Error fetching leaderboard:', error);
            });

            function populateLeaderboard(data) {
                const content = document.getElementById('leaderboardContent');
            
                // Clear existing leaderboard items
                content.innerHTML = '';
            
                // Iterate through the data and create list items
                data.forEach(entry => {
                    content.appendChild(createRow(entry));
                });
            
                function createRow(entry) {
                    const tr = document.createElement('tr');
            
                    for (const property in entry) {
                        if (entry.hasOwnProperty(property)) {
                            const td = document.createElement('td');
                            td.textContent = `${entry[property]}`;
                            tr.appendChild(td);
                        }
                    }
                    return tr;
                }
            }
    }
}