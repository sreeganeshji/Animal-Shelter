import React, { useEffect, useState } from "react";

import { getData } from "../services/api";
import { isLoggedIn } from "../services/auth";
import { NavLinks } from "../components";

const VolunteerLookup = () => {
    const [state, setState] = useState({
        data: [],
        errorMessage: "",
        timeout: -1,
        cookie: {}
    });

    // check auth
    useEffect(() => {
        const cookie = isLoggedIn();
        if (cookie.hasOwnProperty("isTrustedVolunteer") && cookie.isTrustedVolunteer) {
            setState({ ...state, cookie });
        }
        /* eslint-disable react-hooks/exhaustive-deps */
    }, []);

    const [volunteerName, setVolunteerName] = useState("");

    useEffect(() => {
        if (state.cookie.hasOwnProperty("isTrustedVolunteer") && state.cookie.isTrustedVolunteer) {
            getData("/volunteer")
                .then(({ target }) => {
                    if (target) {
                        setState({ ...state, data: JSON.parse(target.response) });
                    } else {
                        setState({ ...state, errorMessage: target.response });
                    }
                })
                .catch(error => {
                    console.log(error);
                    setState({ ...state, errorMessage: JSON.stringify(error) });
                })
        }
    }, [state.cookie]);

    const getVolunteersLike = (volName) => {
        getData(`/volunteer?like=${volName}`)
            .then(({ target }) => {
                if (target.status === 200) {
                    const data = JSON.parse(target.response)
                    setState({ ...state, data });

                }
            })
            .catch(error => {
                console.log(error);
            })
            .finally(() => {
                console.log(state.data);
            })
    };

    return (
        <>
            <NavLinks isTrustedVolunteer={state.cookie.isTrustedVolunteer} />
            <main className="container">
                <h1>Volunteer Lookup</h1>
                <div className="form-unit clearLeft">
                    <label htmlFor="volunteer-name">Volunteer name:</label>
                    <input
                        type="text"
                        value={volunteerName}
                        onChange={({ target }) => {
                            setVolunteerName(target.value);
                            clearTimeout(state.timeout);
                            const timeout = setTimeout(() => {
                                getVolunteersLike(target.value)
                            }, 600)
                            setState({ ...state, timeout });
                        }}
                    />
                </div>

                {
                    state.data.length ? (

                        <table>
                            <thead>
                                <tr>
                                    <th>First Name</th>
                                    <th>Last Name</th>
                                    <th>Email</th>
                                    <th>Phone Number</th>
                                </tr>
                            </thead>
                            <tbody>
                                {
                                    state.data.map((volunteer) => {
                                        return (
                                            <tr>
                                                <td>{volunteer.firstName}</td>
                                                <td>{volunteer.lastName}</td>
                                                <td>{volunteer.email}</td>
                                                <td>{volunteer.cell}</td>
                                            </tr>
                                        )
                                    })
                                }
                            </tbody>
                        </table>
                    ) : (
                            <p> No volunteers found</p>)
                }
                {state.errorMessage && (
                    <div styles={{ color: "red" }}>{state.errorMessage}</div>
                )}
            </main>
        </>
    );
}

export default VolunteerLookup;

