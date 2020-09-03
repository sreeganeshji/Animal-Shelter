const isLoggedIn = () => {
    const cookiePairs = document.cookie.split(";");
    const cookie = cookiePairs.reduce((accumulator, pair) => {
        const [key, value] = pair.split("=");
        switch (value) {
            case "true":
                accumulator[key] = true;
                return accumulator;
            case "false":
                accumulator[key] = false;
                return accumulator;
            default:
                accumulator[key] = null;
                return accumulator;
        }
    }, {});
    return cookie;
}

export { isLoggedIn };