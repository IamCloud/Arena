let playerId;
let teamName;
window.addEventListener('load', function () {
    playerId = this.localStorage.getItem(STORED_PLAYERID);
    if (playerId) {
        teamName = this.localStorage.getItem(STORED_TEAMNAME);
        if (teamName) {
            resumeGame(playerId, teamName);
        } else {
            openNewTeamDialog();
        }
    } else {
        openInitialDialog();
    }

    window.setInterval(function () {
        Server.getLeaderboard();
    }, 1000);
});

function resumeGame() {
    console.log("TODO: Player game is already started, resume previous session.");
}

function openInitialDialog() {
    const template = document.querySelector("#dialog-template");

    const templateClone = template.content.cloneNode(true);
    const dialog = templateClone.querySelector("dialog");
    const dialogTitle = templateClone.querySelector("#dialog-title");
    const dialogText = templateClone.querySelector("#dialog-text");
    const dialogForm = templateClone.querySelector("#dialog-form");

    dialogTitle.textContent = "Welcome !";
    dialogText.textContent = "Welcome to Arena, a game where you build your ideal medieval fantasy teams to beat other players teams across the world !";

    dialogForm.addEventListener("submit", function (ev) {
        ev.preventDefault();
        const teamname = document.getElementById("team-name-input");
        const playerName = document.getElementById("player-name-input");

        if (!teamname.checkValidity() || !playerName.checkValidity()) return;

        const newUUID = uuidv4();
        localStorage.setItem(STORED_PLAYERID, newUUID);
        if (!Server.initPlayer(newUUID, playerName.value, teamname.value)) {
            console.error("initPlayer failed.");
            return;
        }

        dialog.close();
    });
    document.body.appendChild(templateClone);
}

function openNewTeamDialog() {
    const template = document.querySelector("#dialog-template");

    const templateClone = template.content.cloneNode(true);
    const dialog = templateClone.querySelector("dialog");
    const dialogTitle = templateClone.querySelector("#dialog-title");
    const dialogText = templateClone.querySelector("#dialog-text");
    const dialogForm = templateClone.querySelector("#dialog-form");
    const teamname = templateClone.querySelector("#team-name-input");
    const playerName = templateClone.querySelector("#player-name-input");

    dialogTitle.textContent = "Create a new team !";
    dialogText.remove();
    playerName.remove();

    dialogForm.addEventListener("submit", function (ev) {
        ev.preventDefault();

        if (!teamname.checkValidity()) return;
localStorage.setItem(STORED_TEAMNAME, teamname.value);
        if (!Server.createTeam(localStorage.getItem(STORED_PLAYERID), teamname.value)) {
            console.error("createTeam failed.");
            return;
        }

        dialog.close();
    });
    document.body.appendChild(templateClone);
}


function uuidv4() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}