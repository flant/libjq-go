async function getRelease(github, context, core) {
    let releases = [];
    try {
        core.startGroup('Getting list of releases...')
        const allReleases = await github.repos.listReleases({
            ...context.repo
        });
        console.log(JSON.stringify(allReleases));
        core.endGroup();
        releases = allReleases.data;
    }
    catch (error) {
        core.setFailed(`Fetch releases: ${error}`);
        return;
    }

    let tagName = context.ref.replace('refs/tags/', '');
    console.log(`Find release '${tagName}' (${context.ref})`);

    for (let release of releases) {
        if (release.tag_name === tagName) {
            core.startGroup(`Found release '${release.name}' for tag '${release.tag_name}'...`)
            console.log(JSON.stringify(release));
            core.endGroup();

            core.setOutput('id', release.id)
            core.setOutput('upload_url', release.upload_url)

            return release;
        }
        console.log(`Skip release '${release.name}' for tag '${release.tag_name}'`);
    }

    core.setFailed(`Release for tag ${ context.ref } is not found. Stop the workflow.`);
}


module.exports = getRelease