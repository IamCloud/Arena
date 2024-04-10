let playerInfoEl;
let opponentInfoEl;
let gameLogEl;
let playerId;
let characterId;
let fightInProgress = false;
let fightBtn;

window.addEventListener('load', function () {
    playerInfoEl = document.getElementById("player-info");
    opponentInfoEl = document.getElementById("opponent-info");
    gameLogEl = document.getElementById("game-log");
    fightBtn = document.getElementById("fightBtn");

    playerId = this.localStorage.getItem(STORED_PLAYERID);
    if (playerId) {
        characterId = this.localStorage.getItem(STORED_CHARID);
        if (characterId) {
            Server.getLeaderboard();
        } else {
            openNewCharDialog();
        }
    } else {
        openInitialDialog();
    }

    fightBtn.addEventListener("click", (ev) => {
        if (!fightInProgress) {
            startFight();
        }
    });

    //Server.getLeaderboard();
    /*var ping = window.setInterval(function () {
        Server.getLeaderboard()
            .then(data => {
                if (!data) {
                    clearInterval(ping);
                    return;
                }
            }).catch(error => {
                console.error('Error creating team:', error);
                clearInterval(ping);
            });
    }, 1000);*/
});

function startFight() {
    gameLogEl.innerText = "";
    fightInProgress = true;
    fightBtn.setAttribute("disabled", "");

    playerId = this.localStorage.getItem(STORED_PLAYERID);
    characterId = this.localStorage.getItem(STORED_CHARID);
    console.log(`Starting fight for player ${playerId} with character ${characterId}`);

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

function endFight(playerWon) {
    Server.getLeaderboard();
    fightBtn.removeAttribute("disabled");
    fightInProgress = false;

    if (playerWon) {
        openUpgradeDialog();
    }
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
    const wait = (ms) =>
        new Promise(resolve =>
            setTimeout(() => resolve(true), ms)
        );
    const processFightEvents = async (data) => {

        let playerWins = false;
        for (let i = 0; i < data.length; i++) {
            const event = data[i];
            const eventData = JSON.parse(event.Data);
            switch (event.Type) {
                case "init":
                    displayCombatEvent(`${eventData.StartingCharacterName} starts !`);
                    break;
                case "upd":
                    updateInfo(playerInfoEl, eventData.Player);
                    updateInfo(opponentInfoEl, eventData.Opponent);
                    break;
                case "atk":
                    await wait(FIGHT_DELAY_MS);
                    if (eventData.Success) {
                        displayCombatEvent(`<b>${eventData.AttackerName}</b> attacks and hits ! ${eventData.DefenderName} loses <b>${eventData.Damage}</b> health !`);
                    } else {
                        displayCombatEvent(`<b>${eventData.AttackerName}</b> attacks and <i>misses</i> !`);
                    }
                    break;
                case "end":
                    displayCombatEvent(`<b><ins>${eventData.Winner.Name} wins !</ins></b>`);
                    if (eventData.Winner.CharacterId.toString() === characterId) {
                        playerWins = true;
                    }
                    break;
                case "dead":
                    displayCombatEvent(`<b>Your character dies !</b>`);                    
                    await wait(FIGHT_DELAY_MS);
                    gameOver();
                    break;
            }
        }

        endFight(playerWins);
    }

    processFightEvents(data);

    function displayCombatEvent(str) {
        let log = document.createElement("small");
        log.innerHTML = str;
        gameLogEl.appendChild(log);

        const br = document.createElement("br");
        gameLogEl.appendChild(br);

        gameLogEl.scrollTop = gameLogEl.scrollHeight;
    }

    function updateInfo(infoEl, data) {
        infoEl.querySelector(".character-name").textContent = data.Name;
        infoEl.querySelector(".class-name").textContent = `Class: ${data.ClassId}`;

        const livesElem  = infoEl.querySelector(".lives");
        if (livesElem) {
            livesElem.textContent = `Lives: ${data.Lives}`;
        }
        infoEl.querySelector(".hp").textContent = `Health: ${data.Health}/${data.HealthMax}`;
        infoEl.querySelector(".hp-bar").setAttribute("value", data.Health);
        infoEl.querySelector(".hp-bar").setAttribute("max", data.HealthMax);
    }

}



function gameOver() {
    //Remove character local storage
    localStorage.removeItem(STORED_CHARID);

    openNewCharDialog();
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

            const fieldset = document.createElement("fieldset");            
            data.forEach(upgrade => {
                fieldset.appendChild(createUpgradeCard(upgrade));
            });
            upgradesContainer.appendChild(fieldset);

        }).catch(error => {
            console.error('Error getting new upgrades:', error);
        });

    function createUpgradeCard(upgrade) {
        let lbl = document.createElement("label");
        lbl.innerText = `${upgrade.name} (${upgrade.description})`;

        const radio = document.createElement("input");
        radio.type = "radio";
        radio.name = "upgrade";
        radio.value = upgrade.id;

        lbl.appendChild(radio);
        return lbl;
    }


    dialogForm.addEventListener("submit", function (ev) {
        ev.preventDefault();
        const chosenUpgrade = dialogForm.querySelector("input[type='radio']:checked");

        console.log(chosenUpgrade.value);

        Server.chooseUpgrade(characterId, chosenUpgrade.value);
        dialog.close();
        /*Server.createCharacter(localStorage.getItem(STORED_PLAYERID), teamname.value)
            .then(teamId => {
                if (!teamId) {
                    console.error("createTeam failed.");
                    return;
                }
                localStorage.setItem(STORED_CHARID, teamId);

                dialog.close();
            }).catch(error => {
                console.error('Error creating team:', error);
            });*/

    });
    document.body.appendChild(templateClone);
}

function openErrorDialog(msg) {
    const template = document.querySelector("#dialog-error");

    const templateClone = template.content.cloneNode(true);
    const errMsg = templateClone.querySelector("#errmsg");
    errMsg.textContent = msg;

    document.body.appendChild(templateClone);
}