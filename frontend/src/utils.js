export function fetchWithAuth(url, options = {}) {
    const token = localStorage.getItem("wms_token");

    if (!token) {
        window.location.href = "/";
        return Promise.reject("No token");
    }

    options.headers = options.headers || {};
    options.headers["Authorization"] = "Bearer " + token;

    return fetch(url, options)
        .then(async response => {
            if (response.status === 401) {
                localStorage.removeItem("wms_token");
                window.location.href = "/";
                return;
            }

            const data = await response.json().catch(() => null);
            return { ok: response.ok, data };
        });
}
