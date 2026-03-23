
const updateProject = async (projectID) => {
    const name = document.getElementById("name").value;
    const response = await fetch(`/api/projects/${projectID}`, {
        method: "PATCH",
        body: name
    });

    if (response.ok) {
        document.title = name;
    }
}

const deleteProject = async (projectID) => {
    const response = await fetch(`/api/projects/${projectID}`, {
        method: "DELETE",
    });

    if (response.ok) {
        window.location.href = "/projects";
    }
}
