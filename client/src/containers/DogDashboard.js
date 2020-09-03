import React, { useEffect, useState } from "react";
import { Redirect } from "react-router-dom";

import { getData } from "../services/api";
import { isLoggedIn } from "../services/auth";
import { NavLinks } from "../components";

const DogDashboard = () => {

    const CAPACITY = 15;
    const [state, setState] = useState({
        dogs: [],
        addDogRedirect: false,
        addApplicationRedirect: false,
        showAvailable: true,
        showUnavailable: true,
        dogDetailRedirect: "",
        loginRedirect: false,
        cookie: {}
    });

    // check auth
    useEffect(() => {
        const cookie = isLoggedIn();
        if (cookie.hasOwnProperty("isTrustedVolunteer") && cookie.isTrustedVolunteer !== null) {
            setState({ ...state, cookie });
        }
        /* eslint-disable react-hooks/exhaustive-deps */
    }, []);

    const setShowAvailable = () => {
        setState({ ...state, showAvailable: !state.showAvailable });
    }

    const setShowUnavailable = () => {
        setState({ ...state, showUnavailable: !state.showUnavailable });
    }

    // get data for dogs currently in the shelter
    useEffect(() => {
        if (state.cookie.hasOwnProperty("isTrustedVolunteer") && state.cookie.isTrustedVolunteer !== null) {
            getData("/dog?current=true")
                .then(({ target }) => {
                    setState({ ...state, dogs: JSON.parse(target.response) });
                })
                .catch(error => {
                    console.log(error);
                })
        }
    }, [state.cookie]);

    const addApplication = () => {
        setState({ ...state, addApplicationRedirect: true });
    }

    const addDog = () => {
        setState({ ...state, addDogRedirect: true });
    }

    const {
        addApplicationRedirect,
        addDogRedirect,
        dogs,
        showAvailable,
        showUnavailable,
        dogDetailRedirect,
        cookie
    } = state;

    return (
        <>
            {addApplicationRedirect && <Redirect to={{ pathname: "/add-application", state: { from: "/dog-dashboard" } }} />}
            {addDogRedirect && <Redirect to={{ pathname: "/add-dog", state: { from: "/dog-dashboard" } }} />}
            {dogDetailRedirect && <Redirect to={{ pathname: dogDetailRedirect, state: { from: "/dog-dashboard" } }} />}

            <NavLinks isTrustedVolunteer={cookie.isTrustedVolunteer} />
            <main className="container">
                <h1>Dog Dashboard</h1>

                <ControlPanel
                    addApplication={addApplication}
                    addDog={addDog}
                    availableSpaces={CAPACITY - dogs.length}
                    showAvailable={showAvailable}
                    setShowAvailable={setShowAvailable}
                    showUnavailable={showUnavailable}
                    setShowUnavailable={setShowUnavailable}
                />
                {
                    dogs.length ? (
                        <table>
                            <thead>
                                <tr>
                                    <th>Name</th>
                                    <th>Breed</th>
                                    <th>Sex</th>
                                    <th>Alteration Status</th>
                                    <th>Age</th>
                                    <th>Adoptability Status</th>
                                    <th>Date Surrendered</th>
                                </tr>
                            </thead>
                            <tbody>
                                {dogs.map((dog, index) => {
                                    const ONE_DAY = 1000 * 60 * 60 * 24;
                                    const now = Date.now();
                                    const then = Date.parse(dog.dateOfBirth);
                                    const years = Math.floor((now - then) / (ONE_DAY * 365));
                                    const yearsStr = `${years}${!years || years > 1 ? " years " : " year "}`;
                                    const months = Math.floor((now - then) / (ONE_DAY * 30)) % 12;
                                    const monthsStr = `${months}${!months || months > 1 ? " months" : " month"}`;
                                    const isAvailable = dog.microchipId.String && dog.alterationStatus;
                                    if ((isAvailable && showAvailable) || (!isAvailable && showUnavailable)) {
                                        return (
                                            <tr key={index} onClick={() => setState({ ...state, dogDetailRedirect: `/dog-details/${dog.id}` })}>
                                                <td>{dog.name}</td>
                                                <td>{dog.breed.join("/")}</td>
                                                <td>{dog.sex}</td>
                                                <td>{dog.alterationStatus ? "Altered" : "Not yet altered"}</td>
                                                <td>{yearsStr}{monthsStr}</td>
                                                <td>{isAvailable ? "Available" : "Not yet available"}</td>
                                                <td>{(new Date(dog.surrenderDate)).toDateString()}</td>
                                            </tr>
                                        );
                                    }
                                    return null;
                                })}
                            </tbody>

                        </table>
                    ) : (<h3>There are not any dogs currently in the shelter.</h3>)
                }
            </main>
        </>
    );
};

const ControlPanel = ({
    addApplication,
    addDog,
    availableSpaces,
    showAvailable,
    setShowAvailable,
    showUnavailable,
    setShowUnavailable,
    logout
}) => {
    return (
        <div className="control-panel">
            <h3>Currently Available Spaces: {availableSpaces}</h3>
            <div>
                <p>
                    <input
                        type="checkbox"
                        checked={showAvailable}
                        id="showAvailable"
                        onChange={setShowAvailable}
                    />
                    <label
                        htmlFor="showAvailable"
                        style={{ paddingLeft: "10px" }}
                    >Show dogs available for adoption</label>
                </p>
                <p>
                    <input
                        type="checkbox"
                        checked={showUnavailable}
                        id="showUnavailable"
                        onChange={setShowUnavailable}
                    />
                    <label
                        htmlFor="showUnavailable"
                        style={{ paddingLeft: "10px" }}
                    >Show dogs NOT available for adoption</label>
                </p>
            </div>
            <div>
                <p>
                    {availableSpaces && <button onClick={addDog}>Add Dog</button>}
                </p>
                <p>
                    <button onClick={addApplication}>Add Adoption Application</button>
                </p>
            </div>
        </div>
    )
}

export default DogDashboard;