import React, { useEffect, useReducer, useState } from "react"
import { Redirect } from "react-router-dom"

import { getData, sendData } from "../services/api";
import { isLoggedIn } from "../services/auth";
import { NavLinks } from "../components";

const AddDog = () => {
    const [pageState, setPageState] = useState({
        loginRedirect: false,
        cookie: {}
    });

    // check auth
    useEffect(() => {
        const cookie = isLoggedIn();
        if (cookie.hasOwnProperty("isTrustedVolunteer") && cookie.isTrustedVolunteer !== null) {
            setPageState({ ...pageState, cookie });
        }
        /* eslint-disable react-hooks/exhaustive-deps */
    }, []);

    const reducer = (state, action) => {
        switch (action.type) {
            case "dogName":
                return { ...state, dogName: action.payload };
            case "breed":
                let newBreed = state.breed;
                if (action.payload.checked) {
                    newBreed[action.payload.breed] = true;
                    if (action.payload.breed === "Mixed" || action.payload.breed === "Unknown") {
                        // delete all other breeds
                        for (const key in newBreed) {
                            if (key !== action.payload.breed) {
                                delete newBreed[key];
                                document.querySelector(`#${key.split(" ").join("-")}`).checked = false;
                            }
                        }
                    }
                } else {
                    delete newBreed[action.payload.breed];
                }
                return { ...state, breed: newBreed };
            case "sex":
                return { ...state, sex: action.payload }
            case "alterationStatus":
                return { ...state, alterationStatus: action.payload }
            case "dateOfBirth":
                return { ...state, dateOfBirth: action.payload }
            case "description":
                return { ...state, description: action.payload }
            case "microchipId":
                return { ...state, microchipId: action.payload }
            case "expenses":
                return { ...state, expenses: action.payload }
            case "surrenderDate":
                return { ...state, surrenderDate: action.payload }
            case "surrenderReason":
                return { ...state, surrenderReason: action.payload }
            case "surrenderedByAnimalControl":
                return { ...state, surrenderedByAnimalControl: action.payload }
            case "availableBreeds":
                return { ...state, availableBreeds: action.payload }
            case "errorMessage":
                return { ...state, errorMessage: action.payload }
            case "showModal":
                return { ...state, id: action.payload, showModal: true }
            case "dogDetail":
                return { ...state, detailRedirect: true }
            case "dogDashboard":
                return { ...state, dashboardRedirect: true }
            default:
                throw new Error();
        }
    };

    const [state, dispatch] = useReducer(
        reducer,
        {
            dogName: "",
            breed: {},
            sex: null,
            alterationStatus: null,
            dateOfBirth: "",
            description: "",
            microchipId: "",
            surrenderDate: "",
            surrenderReason: "",
            surrenderedByAnimalControl: null,
            availableBreeds: [],
            errorMessage: "",
            showModal: false,
            detailRedirect: false,
            dashboardRedirect: false,
            id: -1
        });

    const {
        dogName,
        breed,
        sex,
        alterationStatus,
        dateOfBirth,
        description,
        microchipId,
        surrenderDate,
        surrenderReason,
        surrenderedByAnimalControl,
        availableBreeds,
        errorMessage,
        showModal,
        detailRedirect,
        dashboardRedirect,
        id
    } = state;

    useEffect(() => {
        if (pageState.cookie.hasOwnProperty("isTrustedVolunteer") && pageState.cookie.isTrustedVolunteer !== null) {
            getData("/breed")
                .then(({ target }) => {
                    dispatch({ type: 'availableBreeds', payload: JSON.parse(target.response) });
                })
                .catch(error => {
                    dispatch({ type: "errorMessage", payload: error });
                    console.log("error:", error);
                })
        }
    }, [pageState.cookie]);

    const addDog = () => {
        const volunteerId = JSON.parse(localStorage.getItem("volunteerId"));
        const dogData = {
            name: dogName,
            breed: Object.keys(breed),
            sex,
            alterationStatus: alterationStatus === "true" ? true : false,
            dateOfBirth,
            description,
            surrenderDate,
            surrenderReason,
            surrenderedByAnimalControl,
            volunteerId
        };
        if (microchipId) {
            dogData.microchipId = { "String": microchipId, "Valid": true };
        }
        if (validateDog()) {
            sendData(dogData, "/dog")
                .then(data => {
                    const { target } = data;
                    if (target.status === 201) {
                        // show success modal
                        dispatch({ type: "showModal", payload: JSON.parse(target.response).id });
                    } else {
                        console.log("error adding dog (non-201):", target.response);
                        dispatch({ type: "errorMessage", payload: target.response });
                    }
                })
                .catch(error => {
                    // show error modal
                    console.log("error adding dog (catch):", error);
                    dispatch({ type: "errorMessage", payload: error });
                });
        }
    }

    const validateDog = () => {
        if (!dogName) {
            dispatch({ type: "errorMessage", payload: "Please enter dog name." });
            return false;
        }
        if (dogName.toLowerCase() === "uga") {
            let isError = false;
            Object.keys(breed).forEach(b => {
                // prohibit any type of bulldog named Uga
                if (b.toLowerCase().indexOf("bulldog") !== -1) {
                    dispatch({ type: "errorMessage", payload: "BullDog can not be named 'Uga'." });
                    isError = true;
                }
            })
            if (isError) {
                return false;
            }
        }
        if (!Object.keys(breed).length) {
            dispatch({ type: "errorMessage", payload: "Please select at least one dog breed." });
            return false;
        }
        if (!sex) {
            dispatch({ type: "errorMessage", payload: "Please indicate dog's sex (or select unknown)." });
            return false;
        }
        if (!alterationStatus) {
            dispatch({ type: "errorMessage", payload: "Please indicate dog's alteration status." });
            return false;
        }
        if (!dateOfBirth) {
            dispatch({ type: "errorMessage", payload: "Please enter dog's date of birth." });
            return false;
        }
        if (!description) {
            dispatch({ type: "errorMessage", payload: "Please enter dog description." });
            return false;
        }
        if (!surrenderDate) {
            dispatch({ type: "errorMessage", payload: "Please enter dog's surrender date." });
            return false;
        }
        if (!surrenderReason) {
            dispatch({ type: "errorMessage", payload: "Please enter reason for dog surrender." });
            return false;
        }
        if (!surrenderedByAnimalControl) {
            dispatch({ type: "errorMessage", payload: "Please indicate whether dog was surrendered by Animal Control." });
            return false;
        }
        return true;
    }

    return (
        <>
            {detailRedirect && <Redirect to={{ pathname: `/dog-details/${id}`, state: { from: "/add-dog" } }} />}
            {dashboardRedirect && <Redirect to={{ pathname: "/dog-dashboard", state: { from: "/add-dog" } }} />}
            <NavLinks isTrustedVolunteer={pageState.cookie.isTrustedVolunteer} />
            <main className="container">
                <h1>Add Dog</h1>
                <div style={{ overflow: showModal ? "hidden" : "scroll", clear: "left" }}>
                    <h3>All fields are required unless otherwise stated.</h3>
                    <form>
                        <div className="form-unit">
                            <label htmlFor="dogName">Dog Name:</label>
                            <input type="text"
                                value={dogName}
                                id="dogName"
                                onChange={({ target }) => {
                                    dispatch({ type: 'dogName', payload: target.value })
                                }}
                            />
                        </div>
                        <div className="form-unit">
                            <label htmlFor="breed">Dog Breed: (Please select at least one.)</label>
                            <div className="scroll-box">
                                {availableBreeds.map(({ id, breed: b = "" }) => (
                                    <p key={id}>
                                        <input type="checkbox"
                                            value={b}
                                            id={b.split(" ").join("-")}
                                            onChange={({ target }) => {
                                                dispatch({ type: 'breed', payload: { breed: target.value, checked: target.checked } })
                                            }}
                                            disabled={(breed["Mixed"] || breed["Unknown"]) && !(b === "Mixed" || b === "Unknown")}
                                        /><label htmlFor={b} style={{ display: "inline", paddingLeft: "10px" }}>{b}</label>
                                    </p>
                                ))}
                            </div>
                        </div>
                        <div className="form-unit">
                            <label htmlFor="sex">Dog Sex:</label>
                            <p>
                                <input type="radio"
                                    id="male"
                                    name="sex"
                                    value="male"
                                    onChange={({ target }) => {
                                        dispatch({ type: 'sex', payload: target.value })
                                    }}
                                />
                                <label htmlFor="male" style={{ display: "inline" }}>Male</label>
                            </p>
                            <p>
                                <input type="radio"
                                    id="female"
                                    name="sex"
                                    value="female"
                                    onChange={({ target }) => {
                                        dispatch({ type: 'sex', payload: target.value })
                                    }}
                                />
                                <label htmlFor="female" style={{ display: "inline" }}>Female</label>
                            </p>
                            <p>
                                <input type="radio"
                                    id="unknown"
                                    name="sex"
                                    value="unknown"
                                    onChange={({ target }) => {
                                        dispatch({ type: 'sex', payload: target.value })
                                    }}
                                />
                                <label htmlFor="unknown" style={{ display: "inline" }}>Unknown</label>
                            </p>

                        </div>
                        <div className="form-unit">
                            <label htmlFor="alterationStatus">Dog Alteration Status:</label>
                            <p>
                                <input type="radio"
                                    id="altered"
                                    name="alterationStatus"
                                    value={true}
                                    onChange={({ target }) => {
                                        dispatch({ type: 'alterationStatus', payload: target.value })
                                    }}
                                />
                                <label htmlFor="altered" style={{ display: "inline" }}>Spayed/Neutered</label>
                            </p>
                            <p>
                                <input type="radio"
                                    id="unaltered"
                                    name="alterationStatus"
                                    value={false}
                                    onChange={({ target }) => {
                                        dispatch({ type: 'alterationStatus', payload: target.value })
                                    }}
                                />
                                <label htmlFor="unaltered" style={{ display: "inline" }}>Unaltered</label>
                            </p>
                        </div>
                        <div className="form-unit">
                            <label htmlFor="dateOfBirth">Dog Date of Birth:</label>
                            <input type="date"
                                value={dateOfBirth}
                                id="dateOfBirth"
                                onChange={({ target }) => {
                                    dispatch({ type: 'dateOfBirth', payload: target.value })
                                }}
                            />
                        </div>
                        <div className="form-unit">
                            <label htmlFor="description">Dog Description:</label>
                            <textarea value={description}
                                id="description"
                                onChange={({ target }) => {
                                    dispatch({ type: 'description', payload: target.value })
                                }}
                            />
                        </div>
                        <div className="form-unit">
                            <label htmlFor="microchipId">Dog Microchip ID: (leave blank if none)</label>
                            <input type="text"
                                value={microchipId}
                                id="microchipId"
                                onChange={({ target }) => {
                                    dispatch({ type: 'microchipId', payload: target.value })
                                }}
                            />
                        </div>
                        <div className="form-unit">
                            <label htmlFor="surrenderDate">Dog Surrender Date:</label>
                            <input type="date"
                                value={surrenderDate}
                                id="surrenderDate"
                                onChange={({ target }) => {
                                    dispatch({ type: 'surrenderDate', payload: target.value })
                                }}
                            />
                        </div>
                        <div className="form-unit">
                            <label htmlFor="surrenderReason">Dog Surrender Reason:</label>
                            <textarea value={surrenderReason}
                                id="surrenderReason"
                                onChange={({ target }) => {
                                    dispatch({ type: 'surrenderReason', payload: target.value })
                                }}
                            />
                        </div>
                        <div className="form-unit">
                            <label htmlFor="surrenderedByAnimalControl">Dog Was Surrendered by Animal Control:</label>
                            <p>
                                <input type="radio"
                                    id="was"
                                    name="surrenderedByAnimalControl"
                                    value="true"
                                    onChange={({ target }) => {
                                        dispatch({ type: 'surrenderedByAnimalControl', payload: target.value })
                                    }}
                                />
                                <label htmlFor="was" style={{ display: "inline" }}>Yes</label>
                            </p>
                            <p>
                                <input type="radio"
                                    id="wasNot"
                                    name="surrenderedByAnimalControl"
                                    value="false"
                                    onChange={({ target }) => {
                                        dispatch({ type: 'surrenderedByAnimalControl', payload: target.value })
                                    }}
                                />
                                <label htmlFor="wasNot" style={{ display: "inline" }}>No</label>
                            </p>
                        </div>
                        <button type="button" onClick={addDog}>Submit</button><span style={{ color: "red", paddingLeft: "10px" }}>{errorMessage}</span>
                    </form>
                    {showModal && (
                        <div style={{ display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", overflow: "hidden" }}>
                            <div style={{ position: "fixed", top: 0, right: 0, bottom: 0, left: 0, zIndex: 10, backgroundColor: "black", opacity: 0.5 }}></div>
                            <div style={{ position: "fixed", left: "50%", bottom: "50%", marginLeft: "-300px", marginBottom: "-150px", display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", zIndex: 100, height: "300px", width: "600px", backgroundColor: "white" }}>
                                <div style={{ marginBottom: "20px" }}>
                                    <span>Dog successfully created.</span>
                                </div>
                                <div style={{ margin: "20px" }}>Would you like to go to</div>
                                <div>
                                    <button onClick={() => dispatch({ type: "dogDetail", payload: id })}>Dog Detail View</button>
                                    <span style={{ margin: "20px" }}> or </span>
                                    <button onClick={() => dispatch({ type: "dogDashboard" })}>Back to Dog Dashboard</button>
                                </div>
                            </div>
                        </div>
                    )}
                </div>
            </main>
        </>
    );
}

export default AddDog;