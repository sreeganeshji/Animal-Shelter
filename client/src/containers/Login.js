import React, { useEffect, useState } from "react";
import { Redirect } from "react-router-dom";

import { getData, sendData } from "../services/api";
import { NavLinks } from "../components";

const Login = () => {
    const [state, setState] = useState({
        email: "",
        password: "",
        errorMessage: "",
        redirect: false
    });

    useEffect(() => {
        getData("/logout")
            .then(data => {
                setState({ ...state, loginRedirect: true });
            })
            .catch(error => {
                console.log(error);
            })
        /* eslint-disable react-hooks/exhaustive-deps */
    }, [])

    const handleEmailChange = ({ target }) => {
        setState({ ...state, email: target.value });
    }

    const handlePasswordChange = ({ target }) => {
        setState({ ...state, password: target.value });
    }

    const submitForm = form => {
        const { email, password } = state;

        if (password && email.indexOf("@") !== 0 && email.length > 4) {
            sendData({ email, password }, "/login")
                .then(data => {
                    const { target } = data;
                    if (target.status === 200) {
                        // redirect to Dog Dashboard
                        const volunteer = JSON.parse(target.response);
                        localStorage.setItem("volunteerId", volunteer.id);
                        localStorage.setItem("isTrusted", volunteer.isTrusted);
                        setState({ ...state, redirect: true });
                    } else {
                        setState({ ...state, errorMessage: target.response })
                    }
                })
                .catch(error => {
                    setState({ ...state, errorMessage: error })
                });
        } else {
            if (!password) {
                setState({ ...state, errorMessage: "Please enter valid password." })
            } else {
                setState({ ...state, errorMessage: "Please enter valid email address." })
            }
        }
    }
    return (
        <>
            <NavLinks titleOnly={true} />
            <main className="container">
                <h1>Login</h1>

                <form action="/login" method="post">

                    <div className="form-unit">
                        <label htmlFor="email">Email:</label>
                        <input type="email" id="email" name="email" onChange={handleEmailChange} value={state.email} autoComplete="username" />
                    </div>

                    <div className="form-unit">
                        <label htmlFor="password">Password:</label>
                        <input type="password" id="password" name="password" onChange={handlePasswordChange} value={state.password} autoComplete="current-password" />
                    </div>

                    <button type="button" onClick={submitForm}>Log In</button>

                </form>
                {state.errorMessage && <div style={{ color: "red" }}>{state.errorMessage}</div>}
                {state.redirect && <Redirect to={{ pathname: "/dog-dashboard", state: { from: "/login" } }} />}
            </main>
        </>
    );
};

export default Login;