
const deleteProject = async (projectID) => {
    const response = await fetch(`/api/projects/${projectID}`, {
        method: "DELETE",
    });

    if (response.ok) {
        window.location.href = "/projects";
    }
}
