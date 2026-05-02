
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
    const rootStyles = getComputedStyle(document.documentElement);
    const deleteButtonDisplay = rootStyles.getPropertyValue("--delete-button-display");
    const newDeleteButtonDisplay = rootStyles.getPropertyValue("--delete-button-display") === "none" ? "block" : "none";

    document.documentElement.style.setProperty("--delete-button-display", newDeleteButtonDisplay);
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
