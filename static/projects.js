
const addProject = async () => {
    const name = document.getElementById("name").value;
    const response = await fetch("/api/projects", {
        method: "POST",
        body: name
    });

    if (response.ok) {
        document.getElementById("projects").insertAdjacentHTML("beforeend", await response.text());
    }
}

const logout = async () => {
    const response = await fetch("/api/logout", {
        method: "POST"
    });

    if (response.ok) {
        window.location.href = "/login.html";
    }
}
