
const addEndpoint = async () => {
    const method = document.getElementById("method-select").value;
    const path = document.getElementById("path").value;
    const description = document.getElementById("description").value;
    const url = window.location.href;
    const projectID = url.substring(url.length - 22);

    const response = await fetch(`/api/projects/${projectID}/endpoints`, {
        method: "POST",
        body: JSON.stringify({
            "method": method,
            "path": path,
            "description": description
        })
    });

    if (response.ok) {
        document.getElementById("method-select").value = "POST";
        const path = document.getElementById("path").value = "";
        const description = document.getElementById("description").value = "";

        document.getElementById("endpoints").insertAdjacentHTML("beforeend", await response.text());
    }
}

const toggleEdit = () => {
    const endpointContainer = document.getElementById("endpoints");
    const endpoints = endpointContainer.querySelectorAll(".endpoint");

    endpoints.forEach((endpoint) => {
        const deleteButton = endpoint.querySelector("button");
        deleteButton.style.display = (deleteButton.style.display === "none") ? "block" : "none";
    })
}

const deleteEndpoint = async (endpointID) => {
    const url = window.location.href;
    const projectID = url.substring(url.length - 22);
    const response = await fetch(`/api/projects/${projectID}/endpoints/${endpointID}`, {
        method: "DELETE"
    });

    if (response.ok) {
        document.getElementById(endpointID).remove();
    }
}
