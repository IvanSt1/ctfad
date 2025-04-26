export function isMale() {
    if (localStorage.getItem("gender") !== null) {
        return localStorage.getItem("gender") === "male";
    }
    return false;
}

export function setGender(gender) {
    localStorage.setItem("gender", gender);
}
