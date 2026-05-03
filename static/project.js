
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

document.getElementById("add-path-parameter").addEventListener("click", (event) => {
    event.preventDefault();
    const pathParameterID = Math.random().toString(36).substring(2, 10);

    const pathParametersContainer = document.getElementById("path-parameters");
    const fragment = document.createDocumentFragment();

    const pathParameterContainer = document.createElement("div");
    pathParameterContainer.id = pathParameterID;
    pathParameterContainer.classList.add("path-parameter-container");

    const nameLabel = document.createElement("label");
    nameLabel.textContent = "Name";
    nameLabel.style.marginRight = "0.25rem"

    const nameInput = document.createElement("input");

    const typeLabel = document.createElement("label");
    typeLabel.textContent = "Type";
    typeLabel.classList.add("path-parameter-type-label");

    const typeInput = document.createElement("input");

    const descriptionLabel = document.createElement("label");
    descriptionLabel.textContent = "Description";
    descriptionLabel.classList.add("path-parameter-description-label");

    const descriptionInput = document.createElement("textarea");
    descriptionInput.style.height = "2.5rem";

    const deletePathParameterButton = document.createElement("button");
    deletePathParameterButton.textContent = "Delete";
    deletePathParameterButton.classList.add("delete-path-parameter");
    deletePathParameterButton.addEventListener("click", () => {
        document.getElementById(pathParameterID).remove();
    });

    pathParameterContainer.append(nameLabel, nameInput, typeLabel, typeInput, descriptionLabel, descriptionInput, deletePathParameterButton)

    fragment.appendChild(pathParameterContainer);
    pathParametersContainer.append(fragment);
});
