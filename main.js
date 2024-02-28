window.addEventListener('load', function () {    
    openInitialDialog();
});

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
        if (!Server.initPlayer(playerName.value, teamname.value)) { console.error("initPlayer failed."); return; }

        dialog.close();
        Server.getLeaderboard();
    });
    document.body.appendChild(templateClone);
}