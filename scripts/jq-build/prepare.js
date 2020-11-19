// Tag format should be:
//   jq-[<build_id>-]<JQ_GIT_SHA>[-<serial>]
// jq- is a prefix to distinguish jq image builds from libjq-go related actions.
// build_id is an optional custom suffix. It can be a branch name, YY.MM, YY.id, etc.
// JQ_GIT_SHA is a required 8 chars of commit sha.
// serial is an optional integer to create new unique git tag and thus rebuild the image.
//
// For example:
// 1. Simple build from commit.
//    jq-b6be13d5
//  The script will checkout commit b6be13d5 from stedolan/jq and build these images:
//   flant/jq:b6be13d5-alpine
//   flant/jq:b6be13d5-ubuntu
//
// 2. to rebuild jq-b6be13d5 in the future, add a serial:
//   jq-b6be13d5-0
// The script will build jq again and push updated images:
//   flant/jq:b6be13d5-alpine
//   flant/jq:b6be13d5-ubuntu
//
// 3. Commit is not human friendly, so use build_id as a description:
//   jq-dec_literal_number-2353d034-0
// The script will checkout tag 'jq-1.6' from stedolan/jq and build these images:
//   flant/jq:dec_literal_number-2353d034-alpine
//   flant/jq:dec_literal_number-2353d034-ubuntu
//
// Note: tag for docker image may contain a maximum of 128 characters.
// The string <build_id> will be truncated to 80 symbols.
// Also do not use / in git tag.
//
// Docker documentation:
//   A tag name must be valid ASCII and may contain lowercase and uppercase letters,
//   digits, underscores, periods and dashes. A tag name may not start with a period
//   or a dash and may contain a maximum of 128 characters.

const gitTagRe = /^jq(-([\da-zA-Z][\da-zA-Z_\-.]*))?-([0-9a-fA-F]{8})(-(\d+))?$/;

function _matchRe(context) {
    let tag = context.ref.replace('refs/tags/', '')
    let match = gitTagRe.exec(tag);
    return match;
    //if (match === null) {
    //    return false;
    //}
    //return true;
}

function _checkTagFormat(context) {
    if (_matchRe(context) === null) {
        return false
    }
    return true
}


function checkTagFormat(github, context, core) {
    // Check tag format. Fail the job if tag is not in shape.
    if (!_checkTagFormat(context)) {
        core.setFailed(`Git tag ${tag} is not suitable for jq build. It should be jq-[<build_id>-]<JQ_GIT_SHA>[-<serial>]. See .github/workflows/jq-build.yaml`);
        return false;
    }
    return true;
}

function prepareEnvsForBuild(github, context, core) {
    if (!checkTagFormat(github, context, core)) {
        return;
    }

    // Ignore jq- and -serial, save buildId and SHA.
    // match[3] is always a git sha
    // match[2] is build_id and can be undefined
    let match = _matchRe(context);
    let gitSha = match[3];
    let buildId = match[2]?match[2]:"";
    let dockerTag = buildId.substring(0, (buildId.length < 80) ? buildId.length : 80);
    dockerTag = dockerTag === "" ? gitSha : `${dockerTag}-${gitSha}`;

    // Export variables for build action.
    core.exportVariable("JQ_GIT_SHA", gitSha)
    core.exportVariable("BUILD_ID", buildId)
    core.exportVariable("DOCKER_TAG", dockerTag)

    // Also set outputs
    core.setOutput('jq_git_sha', gitSha)
    core.setOutput('build_id', buildId)
    core.setOutput('docker_tag', dockerTag)

}

module.exports = {checkTagFormat, prepareEnvsForBuild }