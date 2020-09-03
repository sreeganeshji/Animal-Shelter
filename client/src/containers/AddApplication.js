import React, { useEffect, useState } from "react";
import { Redirect } from "react-router-dom";

import { getData, sendData } from "../services/api";
import { isLoggedIn } from "../services/auth";
import { NavLinks } from "../components";

const AddApplication = () => {
    const [state, setState] = useState({
        showModal: false,
        dashboardRedirect: false,
        errorMessage: "",
        timeout: -1,
        applicationId: "",
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

    const [applicant, setApplicant] = useState({
        city: "",
        email: "",
        firstName: "",
        lastName: "",
        phoneNumber: "",
        state: "",
        street: "",
        zip: ""
    });

    const [application, setApplication] = useState({
        coApplicantFirstName: "",
        coApplicantLastName: "",
        applicantIdFk: "",
        date: "",
        state: ""
    });

    const updateApplicant = (attribute) => {
        let updatedState = applicant;
        Object.assign(updatedState, attribute);
        setApplicant({ ...updatedState });
    };

    const updateApplication = (attribute) => {
        let updatedState = application;
        Object.assign(updatedState, attribute);
        setApplication({ ...updatedState });
    };

    const postApplicant = () => {
        if (applicant.email === "") {
            setState({ ...state, errorMessage: "Please add applicant email." });
            return;
        }
        if (applicant.firstName === "") {
            setState({ ...state, errorMessage: "Please add applicant first name." });
            return;
        }
        if (applicant.lastName === "") {
            setState({ ...state, errorMessage: "Please add applicant last name." });
            return;
        }
        if (applicant.street === "") {
            setState({ ...state, errorMessage: "Please add applicant street address." });
            return;
        }
        if (applicant.city === "") {
            setState({ ...state, errorMessage: "Please add applicant city." });
            return;
        }
        if (applicant.state === "") {
            setState({ ...state, errorMessage: "Please add applicant state." });
            return;
        }
        if (applicant.state.length > 2) {
            setState({ ...state, errorMessage: "Please use 2-letter state abbreviation." });
            return;
        }
        if (applicant.zip === "") {
            setState({ ...state, errorMessage: "Please add applicant zip code." });
            return;
        }
        if (applicant.phoneNumber === "") {
            setState({ ...state, errorMessage: "Please add applicant phone number." });
            return;
        }
        sendData(applicant, `/applicant`)
            .then(({ target }) => {
                const applicant = JSON.parse(target.response);
                if (target.status === 201 && applicant.email) {
                    // show success modal
                    setApplicant(applicant)
                    updateApplication({ applicantIdFk: applicant.id })
                    postApplication()
                } else {
                    console.log("error adding application (non-201):", target.response);
                    setState({ ...state, "errorMessage": target.response });
                }
            })
            .catch(error => {
                // show error modal
                console.log("error adding application (catch):", error);
                setState({ ...state, "errorMessage": error });
            });
    };

    const postApplication = () => {
        if (applicant.email === "") {
            setState({ ...state, errorMessage: "Please add applicant email." });
            return;
        }
        if (applicant.firstName === "") {
            setState({ ...state, errorMessage: "Please add applicant first name." });
            return;
        }
        if (applicant.lastName === "") {
            setState({ ...state, errorMessage: "Please add applicant last name." });
            return;
        }
        if (applicant.street === "") {
            setState({ ...state, errorMessage: "Please add applicant street address." });
            return;
        }
        if (applicant.city === "") {
            setState({ ...state, errorMessage: "Please add applicant city." });
            return;
        }
        if (applicant.state === "") {
            setState({ ...state, errorMessage: "Please add applicant state." });
            return;
        }
        if (applicant.state.length > 2) {
            setState({ ...state, errorMessage: "Please use 2-letter state abbreviation." });
            return;
        }
        if (applicant.zip === "") {
            setState({ ...state, errorMessage: "Please add applicant zip code." });
            return;
        }
        if (applicant.phoneNumber === "") {
            setState({ ...state, errorMessage: "Please add applicant phone number." });
            return;
        }
        sendData(application, `/application`)
            .then(({ target }) => {
                if (target.status === 201) {
                    // show success modal
                    setState({ ...state, applicationId: JSON.parse(target.response).id, showModal: true });
                } else {
                    console.log("error adding application (non-201):", target.response);
                    setState({ ...state, "errorMessage": target.response });
                }
            })
            .catch(error => {
                // show error modal
                console.log("error adding application (catch):", error);
                setState({ ...state, "errorMessage": error });
            });
    };

    const submitApplication = () => {
        if (application.applicantIdFk) {
            postApplication();
        } else {
            postApplicant();
        }
    };

    const resetPage = () => {
        setApplicant({
            city: "",
            email: "",
            firstName: "",
            lastName: "",
            phoneNumber: "",
            state: "",
            street: "",
            zip: ""
        });

        setApplication({
            coApplicantFirstName: "",
            coApplicantLastName: "",
            applicantIdFk: "",
            date: "",
            state: ""
        });

        setState({
            ...state,
            showModal: false,
            dashboardRedirect: false,
            errorMessage: "",
            timeout: -1,
            applicationId: "",
            loginRedirect: false
        });
    }

    return (
        <>
            {state.dashboardRedirect && <Redirect to={{ pathname: "/dog-dashboard", state: { from: "/add-dog" } }} />}
            <NavLinks isTrustedVolunteer={state.cookie.isTrustedVolunteer} />
            <main className="container">
                <h1>Add Application</h1>
                <form>
                    <div className="form-unit">
                        <label htmlFor="">Applicant Email Address:</label>
                        <input
                            type="text"
                            id="applicant-email"
                            value={applicant.email}
                            onChange={({ target }) => {
                                updateApplicant({ "email": target.value })
                                // debounce for 300 milliseconds
                                clearTimeout(state.timeout);
                                setState({
                                    ...state, timeout: setTimeout(() => {
                                        // make http call for applicant
                                        getData(`/applicant/${applicant.email}`)
                                            .then(({ target }) => {
                                                let applicant = JSON.parse(target.response);
                                                if (target.status === 200 && applicant.email) {
                                                    setApplicant(applicant);
                                                    updateApplication({ applicantIdFk: applicant.id })
                                                }
                                            })
                                            .catch(error => {
                                                console.log(error);
                                            })
                                    }, 400)
                                });
                            }}
                        />
                    </div>
                    <div className="form-unit">
                        <label htmlFor="">Applicant First Name:</label>
                        <input
                            type="text"
                            id="applicant-first-name"
                            value={applicant.firstName}
                            onChange={({ target }) => updateApplicant({ "firstName": target.value })}
                        />
                    </div>
                    <div className="form-unit">
                        <label htmlFor="">Applicant Last Name:</label>
                        <input
                            type="text"
                            id="applicant-last-name"
                            value={applicant.lastName}
                            onChange={({ target }) => updateApplicant({ "lastName": target.value })}
                        />
                    </div>
                    <div className="form-unit">
                        <label htmlFor="">Street Address:</label>
                        <input
                            type="text"
                            id="applicant-stree-address"
                            value={applicant.street}
                            onChange={({ target }) => updateApplicant({ "street": target.value })}
                        />
                    </div>
                    <div className="form-unit">
                        <label htmlFor="">City:</label>
                        <input
                            type="text"
                            id="applicant-city"
                            value={applicant.city}
                            onChange={({ target }) => updateApplicant({ "city": target.value })}
                        />
                    </div>
                    <div className="form-unit">
                        <label htmlFor="">State: (2 Letters)</label>
                        <input
                            type="text"
                            id="applicant-state"
                            value={applicant.state}
                            onChange={({ target }) => updateApplicant({ "state": target.value })}
                        />
                    </div>
                    <div className="form-unit">
                        <label htmlFor="">Zip Code:</label>
                        <input
                            type="text"
                            id="applicant-zip-code"
                            value={applicant.zip}
                            onChange={({ target }) => updateApplicant({ "zip": target.value })}
                        />
                    </div>
                    <div className="form-unit">
                        <label htmlFor="">Phone Number:</label>
                        <input
                            type="text"
                            id="applicant-phone-number"
                            value={applicant.phoneNumber}
                            onChange={({ target }) => updateApplicant({ "phoneNumber": target.value })}
                        />
                    </div>
                    <div className="form-unit">
                        <label htmlFor="">Co-Applicant First Name:</label>
                        <input
                            type="text"
                            id="applicant-co-applicant-first-name"
                            value={applicant.coApplicantFirstName}
                            onChange={({ target }) => updateApplication({ "coApplicantFirstName": target.value })}
                        />
                    </div>
                    <div className="form-unit">
                        <label htmlFor="">Co-Applicant Last Name:</label>
                        <input
                            type="text"
                            id="applicant-co-applicant-last-name"
                            value={applicant.coApplicantLastName}
                            onChange={({ target }) => updateApplication({ "coApplicantLastName": target.value })}
                        />
                    </div>
                    <button type="button" onClick={submitApplication}>Submit Application</button><span style={{ color: "red", paddingLeft: "10px" }}>{state.errorMessage}</span>
                </form>
                {state.showModal && (
                    <div style={{ display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", overflow: "hidden" }}>
                        <div style={{ position: "fixed", top: 0, right: 0, bottom: 0, left: 0, zIndex: 10, backgroundColor: "black", opacity: 0.5 }}></div>
                        <div style={{ position: "fixed", left: "50%", bottom: "50%", marginLeft: "-300px", marginBottom: "-150px", display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", zIndex: 100, height: "300px", width: "600px", backgroundColor: "white" }}>
                            <div style={{ marginBottom: "20px" }}>
                                <span>Application successfully created.</span>
                            </div>
                            <div style={{ marginBottom: "20px" }}>
                                <span>Application #{state.applicationId}</span>
                            </div>
                            <div style={{ margin: "20px" }}>Would you like to go to</div>
                            <div>
                                <button onClick={resetPage}>Add another application</button>
                                <span style={{ margin: "20px" }}> or </span>
                                <button onClick={() => setState({ ...state, dashboardRedirect: true })}>Back to Dog Dashboard</button>
                            </div>
                        </div>
                    </div>
                )}
            </main>
        </>
    )


}

export default AddApplication;