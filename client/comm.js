let Server = {
    createCharacter: function (playerId, characterName, classId) {
        return new Promise((resolve, reject) => {
            fetch("/createcharacter", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ PlayerId: playerId, CharacterName: characterName, ClassId: classId })
            })
                .then(response => {
                    if (!response.ok) {
                        console.error('Network response was not ok');
                        return;
                    }
                    return response.json();
                })
                .then(data => {
                    resolve(data.character_id);
                })
                .catch(error => {
                    console.error('Error creating character:', error);
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
        return new Promise((resolve, reject) => {
        fetch('/getleaderboard')
            .then(response => response.json())
            .then(data => {
                // Call function to populate leaderboard with the retrieved data
                if (data) {
                    populateLeaderboard(data);
                }
                resolve(true);
            })
            .catch(error => {
                console.error('Error fetching leaderboard:', error);
                openErrorDialog(`Server connection error. Error fetching leaderboard. Refresh the page and if the issue persists, please contact the developer: nicolasntr11@gmail.com.`);
                resolve(false);
            });
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
    },
    simulateFight: function (characterId) {
        return new Promise((resolve, reject) => {
            fetch("/simulatefight", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ CharacterId: characterId })
            })
                .then(response => {
                    if (!response.ok) {
                        console.error('Network response was not ok');
                        return;
                    }
                    return response.json();
                })
                .then(data => {
                    resolve(data);
                })
                .catch(error => {
                    console.error('Error simulating fight:', error);
                    reject(error);
                })
        });
    },

}