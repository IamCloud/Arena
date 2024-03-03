let Server = {
    createTeam: function (playerId, teamName) {
        return new Promise((resolve, reject) => {
            fetch("/createteam", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ PlayerId: playerId, TeamName: teamName })
            })
                .then(response => {
                    if (!response.ok) {
                        console.error('Network response was not ok');
                        return;
                    }
                    return response.json();
                })
                .then(data => {
                    resolve(data.team_id);
                })
                .catch(error => {
                    console.error('Error creating team:', error);
                    reject(error);
                })
        });
    },
    createPlayer: function (playerName) {
        return new Promise((resolve, reject) => {
            fetch("/createplayer", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ Name: playerName })
            })
                .then(response => {
                    if (!response.ok) {
                        console.error('Network response was not ok');
                        return;
                    }
                    return response.json();
                })
                .then(data => {
                    resolve(data.player_guid);
                })
                .catch(error => {
                    console.error('Error creating player:', error);
                    reject(error);
                })
        });
    },

    // Function to fetch leaderboard data from the Go server
    getLeaderboard: function () {
        fetch('/getleaderboard')
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
    },
    getNewUpgrades: function () {
        return new Promise((resolve, reject) => {
            fetch('/getnewupgrades')
                .then(response => response.json())
                .then(data => {
                    resolve(data);
                })
                .catch(error => {
                    console.error('Error getting new upgrades:', error);
                    reject(error);
                });
        });
    }
    
}