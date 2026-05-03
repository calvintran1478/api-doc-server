
const addEndpoint = async () => {
    const method = document.getElementById("method-select").value;
    const path = document.getElementById("path").value;
    const description = document.getElementById("description").value;
    const url = window.location.href;
    const projectID = url.substring(url.length - 22);

    const pathParameters = [];
    const pathParametersContainer = document.getElementById("path-parameters");
    for (const pathParameterContainer of pathParametersContainer.children) {
        const nameTypeInputs = pathParameterContainer.querySelectorAll("input");
        const descriptionInput = pathParameterContainer.querySelector("textarea");
        pathParameters.push({
            "name": nameTypeInputs[0].value,
            "type": nameTypeInputs[1].value,
            "description": descriptionInput.value
        });
    }

    const response = await fetch(`/api/projects/${projectID}/endpoints`, {
        method: "POST",
        body: JSON.stringify({
            "method": method,
            "path": path,
            "description": description,
            "path_parameters": pathParameters
        })
    });

    if (response.ok) {
        document.getElementById("method-select").value = "POST";
        const path = document.getElementById("path").value = "";
        const description = document.getElementById("description").value = "";
        pathParametersContainer.replaceChildren();

        document.getElementById("endpoints").insertAdjacentHTML("beforeend", await response.text());
    }
}

const toggleEdit = () => {
    const editButton = document.querySelector(".edit-button");
    editButton.textContent = editButton.textContent === "Edit" ? "Done" : "Edit";

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

const deletePathParameter = (pathParameterID) => {
    document.getElementById(pathParameterID).remove();
}

document.getElementById("add-path-parameter").addEventListener("click", (event) => {
    event.preventDefault();
    const pathParameterID = Math.random().toString(36).substring(2, 10);

    const html = `
        <div id=${pathParameterID} class="path-parameter-container"><label style="margin-right: 0.25rem;">Name</label><input><label class="path-parameter-type-label">Type</label><input><label class="path-parameter-description-label">Description</label><textarea style="height: 2.5rem;"></textarea><button class="delete-path-parameter" onclick="deletePathParameter('${pathParameterID}')">Delete</button></div>
    `

    document.getElementById("path-parameters").insertAdjacentHTML("beforeend", html);
});
