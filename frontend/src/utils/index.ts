const providers = {
    github: {
        client_id: "a7212b06f4ba590d04ee",
        redirect_uri: "http://localhost:3000/callback",
        scope: "user:email",
    },
    gitlab: {
        client_id:
            "a88c7075c5f5c029d5803f0f6d08490140e9d1deb318b270db04194ccc3b1527",
        redirect_uri: "http://localhost:3000/gitlab",
        response_type: "code",
        scope: "read_user",
    },
};

export function getGitHubUrl(from: string, provider: keyof typeof providers) {
    const rootURl =
        provider === "github"
            ? "https://github.com/login/oauth/authorize"
            : "https://gitlab.com/oauth/authorize";

    const options = {
        ...providers[provider],
        state: from,
    };

    const qs = new URLSearchParams(options);

    return `${rootURl}?${qs.toString()}`;
}
