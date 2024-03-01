let playerId;
let teamId;
window.addEventListener('load', function () {
    playerId = this.localStorage.getItem(STORED_PLAYERID);
    if (playerId) {
        teamId = this.localStorage.getItem(STORED_TEAMID);
        if (teamId) {
            startGame();
        } else {
            openNewTeamDialog();
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
    console.log(`Starting game for player ${playerId} with team ${teamId}`);
}

function gameOver() {
    //Remove team local storage
    localStorage.removeItem(STORED_TEAMID);

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
    dialogText.textContent = "Welcome to Arena, a game where you build your ideal medieval fantasy teams to beat other players teams across the world !";

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
                openNewTeamDialog();
            }).catch(error => {
                console.error('Error creating player:', error);
            });

    });
    document.body.appendChild(templateClone);
}

function openNewTeamDialog() {
    const template = document.querySelector("#dialog-template-new-team");

    const templateClone = template.content.cloneNode(true);
    const dialog = templateClone.querySelector("dialog");
    const dialogTitle = templateClone.querySelector("#dialog-title");
    const dialogText = templateClone.querySelector("#dialog-text");
    const dialogForm = templateClone.querySelector("#new-team-form");
    const teamname = templateClone.querySelector("#team-name-input");

    dialogTitle.textContent = "Create a new team !";
    dialogText.remove();

    dialogForm.addEventListener("submit", function (ev) {
        ev.preventDefault();

        if (!teamname.checkValidity()) return;
        localStorage.setItem(STORED_TEAMID, teamname.value);

        Server.createTeam(localStorage.getItem(STORED_PLAYERID), teamname.value)
            .then(teamId => {
                if (!teamId) {
                    console.error("createTeam failed.");
                    return;
                }
                localStorage.setItem(STORED_TEAMID, teamId);

                dialog.close();
                startGame();
            }).catch(error => {
                console.error('Error creating team:', error);
            });

    });
    document.body.appendChild(templateClone);
}