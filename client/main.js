let playerId;
let characterId;
window.addEventListener('load', function () {
    playerId = this.localStorage.getItem(STORED_PLAYERID);
    if (playerId) {
        characterId = this.localStorage.getItem(STORED_CHARID);
        if (characterId) {
            startGame();
        } else {
            openNewCharDialog();
        }
    } else {
        openInitialDialog();
    }

    Server.getLeaderboard();
    window.setInterval(function () {
        Server.getLeaderboard();
    }, 1000);
});

function startGame() {
    playerId = this.localStorage.getItem(STORED_PLAYERID);
    characterId = this.localStorage.getItem(STORED_CHARID);
    console.log(`Starting game for player ${playerId} with character ${characterId}`);

    Server.simulateFight(characterId)
        .then(data => {
            if (!data) {
                console.error("simulateFight failed.");
                return;
            }
            displayFightEvents(data);

        }).catch(error => {
            console.error('Error getting new upgrades:', error);
        });
}

class Character {
    constructor(name, wins, health, initiative, lives, resource, resource_max, class_id) {
        this.name = name;
        this.wins = wins;
        this.health = health;
        this.initiative = initiative;
        this.lives = lives;
        this.resource = resource;
        this.resource_max = resource_max;
        this.class_id = class_id;
    }
}
function displayFightEvents(data) {
    console.log(data);

    for (let i = 0; i < data.length; i++) {
        const event = data[i];
        displayCombatEvent(event);
    }

    function displayCombatEvent(str) {
        console.log(str);
    }
}



function gameOver() {
    //Remove character local storage
    localStorage.removeItem(STORED_CHARID);

    // TODO: Save to database for leaderboard

    // TODO: Open gameover dialog

}

function openInitialDialog() {
    const template = document.querySelector("#dialog-template");

    const templateClone = template.content.cloneNode(true);
    const dialog = templateClone.querySelector("dialog");
    const dialogTitle = templateClone.querySelector("#dialog-title");
    const dialogText = templateClone.querySelector("#dialog-text");
    const dialogForm = templateClone.querySelector("#dialog-form");

    dialogTitle.textContent = "Welcome !";
    dialogText.textContent = "Welcome to Arena, a game where you build up a medieval fantasy character to beat other players characters across the world !";

    dialogForm.addEventListener("submit", function (ev) {
        ev.preventDefault();
        const playerName = document.getElementById("player-name-input");

        if (!playerName.checkValidity()) return;

        Server.createPlayer(playerName.value)
            .then(playerGuid => {
                if (!playerGuid) {
                    console.error("createPlayer failed.");
                    return;
                }
                localStorage.setItem(STORED_PLAYERID, playerGuid);

                dialog.close();
                openNewCharDialog();
            }).catch(error => {
                console.error('Error creating player:', error);
            });

    });
    document.body.appendChild(templateClone);
}

function openNewCharDialog() {
    const template = document.querySelector("#dialog-template-new-char");

    const templateClone = template.content.cloneNode(true);
    const dialog = templateClone.querySelector("dialog");
    const dialogTitle = templateClone.querySelector("#dialog-title");
    const dialogText = templateClone.querySelector("#dialog-text");
    const dialogForm = templateClone.querySelector("#new-char-form");
    const charactername = templateClone.querySelector("#char-name-input");
    const classId = templateClone.querySelector("input[name='class']:checked").value;

    dialogTitle.textContent = "Create a new character !";
    dialogText.remove();

    dialogForm.addEventListener("submit", function (ev) {
        ev.preventDefault();

        if (!charactername.checkValidity()) return;
        localStorage.setItem(STORED_CHARID, charactername.value);

        Server.createCharacter(localStorage.getItem(STORED_PLAYERID), charactername.value, classId)
            .then(characterId => {
                if (!characterId) {
                    console.error("createCharacter failed.");
                    return;
                }
                localStorage.setItem(STORED_CHARID, characterId);

                dialog.close();
                startGame();
            }).catch(error => {
                console.error('Error creating character:', error);
            });

    });
    document.body.appendChild(templateClone);
}

function openUpgradeDialog() {
    const template = document.querySelector("#dialog-template-upgrade");

    const templateClone = template.content.cloneNode(true);
    const dialog = templateClone.querySelector("dialog");
    const dialogForm = templateClone.querySelector("#upgrade-form");
    const upgradesContainer = dialogForm.querySelector("#upgrades-container");

    Server.getNewUpgrades()
        .then(data => {
            if (!data) {
                console.error("getNewUpgrades failed.");
                return;
            }

            data.forEach(upgrade => {
                createUpgradeCard(upgrade);
            });

        }).catch(error => {
            console.error('Error getting new upgrades:', error);
        });

    function createUpgradeCard(upgrade) {
        let div = document.createElement("div");
        div.innerHTML = upgrade.name;

        upgradesContainer.appendChild(div);
    }


    dialogForm.addEventListener("submit", function (ev) {
        ev.preventDefault();

        Server.createCharacter(localStorage.getItem(STORED_PLAYERID), teamname.value)
            .then(teamId => {
                if (!teamId) {
                    console.error("createTeam failed.");
                    return;
                }
                localStorage.setItem(STORED_CHARID, teamId);

                dialog.close();
            }).catch(error => {
                console.error('Error creating team:', error);
            });

    });
    document.body.appendChild(templateClone);
}