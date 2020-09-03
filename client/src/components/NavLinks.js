import React, { useEffect, useState } from "react";
import { Link, Redirect } from "react-router-dom";

import { isLoggedIn } from "../services/auth";

export default function NavLinks({ isTrustedVolunteer, titleOnly = false }) {
    const [state, setState] = useState({
        loginRedirect: false,
        cookie: {}
    });

    // check auth
    useEffect(() => {
        const cookie = isLoggedIn();
        if (cookie.hasOwnProperty("isTrustedVolunteer") && cookie.isTrustedVolunteer !== null) {
            setState({ ...state, cookie });
        } else {
            setState({ ...state, loginRedirect: true });
        }
        /* eslint-disable react-hooks/exhaustive-deps */
    }, []);

    return (
        <>
            {state.loginRedirect && <Redirect to={{ pathname: "/login" }} />}
            <nav className="container">
                <h1>Mo's Mutt House</h1>
                {!titleOnly && (
                    <ul>
                        <li className="nav-link form-unit log-out">
                            <Link to="/login">Logout</Link>
                        </li>
                        {isTrustedVolunteer && (
                            <>
                                <li className="nav-link form-unit">
                                    <Link to="/reports">Reports</Link>
                                </li>
                                <li className="nav-link form-unit">
                                    <Link to="/review-applications">Review Applications</Link>
                                </li>
                            </>
                        )}

                        <li className="nav-link form-unit">
                            <Link to="/add-application">Add Application</Link>
                        </li>
                        <li className="nav-link form-unit">
                            <Link to="/add-dog">Add Dog</Link>
                        </li>
                        <li className="nav-link form-unit">
                            <Link to="/dog-dashboard">Dog Dashboard</Link>
                        </li>
                    </ul>
                )}
            </nav>
        </>
    );
}