const sendData = (data, path, method = "post") => {
    return new Promise((resolve, reject) => {
        const API = "http://localhost:8080/api";
        const request = new XMLHttpRequest();

        request.onload = event => {
            resolve(event);
        }

        request.onerror = error => {
            reject(error);
        }

        request.open(method, API + path);

        request.send(JSON.stringify(data));
    });
}

const getData = path => {
    return new Promise((resolve, reject) => {
        const API = "http://localhost:8080/api";
        const request = new XMLHttpRequest();

        request.onload = data => {
            resolve(data);
        }

        request.onerror = error => {
            reject(error);
        }

        request.open("get", API + path);

        request.send();
    });
}

export {
    sendData,
    getData
};