import axios from 'axios';

export async function canAccess() {
    try {
        const res = await axios.post("http://localhost:8083/api/check", {}, {
            withCredentials: true,
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            },
        });
        if (res.status === 200) {
            return true;
        }
        return false;
    } catch(error) {
        return false;
    }
}

export function isAuthorized() {
    if (localStorage.getItem("authorized") !== null ) {
        return localStorage.getItem("authorized") == "true";
    }
    return false;
}

export function authorize() {
    if (!isAuthorized()) {
        localStorage.setItem("authorized", "true");
    }
}

export function removeAuthorize() {
    if (isAuthorized()) {
        localStorage.removeItem("authorized");
    }
}
