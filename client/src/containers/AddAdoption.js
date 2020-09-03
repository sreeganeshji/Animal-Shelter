import React, { useEffect, useState } from "react";
import { Redirect, useLocation } from "react-router-dom";

import { getData, sendData } from "../services/api";
import { isLoggedIn } from "../services/auth";
import { NavLinks } from "../components";

const AddAdoption = () => {
    const locationState = useLocation().state;

    const [state, setState] = useState({
        approvedApplicants: [],
        timeout: -1,
        latestApplication: {},
        dog: locationState && locationState.dog,
        dateAdopted: "",
        dashboardRedirect: false,
        errorMessage: "",
        cookie: {}
    });

    const [approvedApplicantLastNameFragment, setApprovedApplicantLastNameFragment] = useState("");

    // check auth
    useEffect(() => {
        const cookie = isLoggedIn();
        if (cookie.hasOwnProperty("isTrustedVolunteer") && cookie.isTrustedVolunteer) {
            setState({ ...state, cookie });
        }
        /* eslint-disable react-hooks/exhaustive-deps */
    }, []);

    // Don't pre-populate
    // useEffect(() => {
    //     if (state.cookie.hasOwnProperty("isTrustedVolunteer") && state.cookie.isTrustedVolunteer) {
    //         getData(`/applicant?status=approved`)
    //             .then(({ target }) => {
    //                 setState({ ...state, approvedApplicants: JSON.parse(target.response) });
    //             })
    //             .catch(error => {
    //                 console.log(error);
    //             })
    //     }
    // }, [state.cookie]);

    const getLatestApplication = (params, applicant) => {
        const query = Object.keys(params)
            .filter(key => params[key])
            .map((key, index) => `${index ? "&" : "?"}${key}=${params[key]}`)
            .join("")

        getData(`/application${query}`)
            .then(({ target }) => {
                if (target.status === 200) {
                    const latestApplication = JSON.parse(target.response)[0]
                    latestApplication.applicant = applicant;
                    setState({ ...state, latestApplication });
                } else {
                    console.log(target);
                }

            })
            .catch(error => {
                console.log(error);
            });
    }

    const submitAdoption = () => {
        if (state.dateAdopted === "") {
            setState({ ...state, errorMessage: "Please enter adoption date." });
            return;
        }

        const adoption = {
            dateAdopted: state.dateAdopted,
            applicationNumberFk: +state.latestApplication.id,
            dogIdFk: state.dog.id
        }
        sendData(adoption, "/adoption")
            .then(data => {
                const { target } = data;
                if (target.status === 201) {
                    // show success modal
                    setState({ ...state, dashboardRedirect: true });
                } else {
                    console.log("error adding dog (non-201):", target.response);
                    setState({ ...state, errorMessage: target.response });
                }
            })
            .catch(error => {
                // show error modal
                console.log("error adding dog (catch):", error);
                setState({ ...state, errorMessage: error });
            });
    }

    const totalExpenses = state.dog && state.dog.expenses.reduce((sum, current) => sum + (current.amountInCents / 100), 0);

    const animalControlPays = state.dog && ((state.dog.surrenderedByAnimalControl && totalExpenses) || 0);

    const adoptionFee = (1.15 * totalExpenses) - animalControlPays;

    const adoptionFeeCents = Math.round((adoptionFee % 1) * 100); // round to nearest cent

    const adoptionFeeDollars = Math.floor(adoptionFee); // round down to nearest dollar

    return (
        <>
            {state.dashboardRedirect && <Redirect to={{ pathname: "/dog-dashboard", state: { from: "/add-adoption" } }} />}
            <NavLinks isTrustedVolunteer={state.cookie.isTrustedVolunteer} />
            <main className="container">
                <h1>Add Adoption</h1>
                {!state.latestApplication.id && <form>
                    <div className="form-unit">
                        <label htmlFor="approved-applicant-last-name">Search by approved applicant last name:</label>
                        <input
                            id="approved-applicant-last-name"
                            onChange={({ target }) => {
                                setApprovedApplicantLastNameFragment(target.value);
                                clearTimeout(state.timeout);
                                setState({
                                    ...state, timeout: setTimeout(() => {
                                        // make http call for approved applicants like
                                        getData(`/applicant?like=${target.value}`)
                                            .then(({ target }) => {
                                                let approvedApplicants = JSON.parse(target.response);
                                                if (target.status === 200) {
                                                    setState({ ...state, approvedApplicants });
                                                }
                                                else {
                                                    console.log(target);
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

                </form>}

                {!state.latestApplication.id && state.approvedApplicants.length > 0 && (
                    <>
                        <h3>Select Applicant</h3>
                        <table>
                            <thead>
                                <tr>
                                    <th></th>
                                    <th>First Name</th>
                                    <th>Last Name</th>
                                    <th>Street</th>
                                    <th>City</th>
                                    <th>State</th>
                                    <th>Zip</th>
                                    <th>Email</th>
                                    <th>Phone Number</th>
                                    <th>Co-applicant Name</th>
                                </tr>
                            </thead>
                            <tbody>
                                {state.approvedApplicants.length > 0 && state.approvedApplicants.map((applicant, index) => {
                                    if (applicant.coApplicantLastName) {
                                        console.log(applicant.coApplicantLastName.indexOf(approvedApplicantLastNameFragment));
                                    }

                                    return (
                                        <tr key={index} >
                                            <td>
                                                <button type="button" onClick={() => {
                                                    if (applicant.coApplicantLastName && applicant.coApplicantLastName.indexOf(approvedApplicantLastNameFragment) !== -1) {
                                                        getLatestApplication({ applicantId: applicant.id, coApplicantFirstName: applicant.coApplicantFirstName, coApplicantLastName: applicant.coApplicantLastName }, applicant);
                                                    } else {
                                                        getLatestApplication({ applicantId: applicant.id }, applicant);
                                                    }
                                                }}>Select</button>
                                            </td>
                                            <td>{applicant.firstName}</td>
                                            <td>{applicant.lastName}</td>
                                            <td>{applicant.street}</td>
                                            <td>{applicant.city}</td>
                                            <td>{applicant.state}</td>
                                            <td>{applicant.zip}</td>
                                            <td>{applicant.email}</td>
                                            <td>{applicant.phoneNumber}</td>
                                            <td>{(applicant.coApplicantLastName && applicant.coApplicantLastName.indexOf(approvedApplicantLastNameFragment) !== -1) ? `${applicant.coApplicantFirstName || ""} ${applicant.coApplicantLastName || ""}` : ""}</td>
                                        </tr>
                                    );

                                })}
                            </tbody>
                        </table>
                    </>
                )}
                {/* {!state.latestApplication.id && !state.approvedApplicants.length && (<h3>There are not currently any approved applications.</h3>)} */}
                {state.latestApplication.id && (
                    <>
                        <h3>Selected Application</h3>
                        <div className="form-unit">Application Number: {state.latestApplication.id}</div>
                        <div className="form-unit">Date: {new Date(state.latestApplication.date).toString()}</div>
                        <div className="form-unit">Applicant Name: {state.latestApplication.applicant.firstName} {state.latestApplication.applicant.lastName}</div>
                        <div className="form-unit">Applicant Address: {state.latestApplication.applicant.street}, {state.latestApplication.applicant.city}, {state.latestApplication.applicant.state}  {state.latestApplication.applicant.zipCode}</div>
                        <div className="form-unit">Email: {state.latestApplication.applicant.email}</div>
                        <div className="form-unit">Phone Number: {state.latestApplication.applicant.phoneNumber}</div>
                        {(state.latestApplication.coApplicantFirstName || state.latestApplication.coApplicantLastName) && <div className="form-unit">Co-Applicant: {state.latestApplication.coApplicantFirstName} {state.latestApplication.coApplicantLastName}</div>}
                        <h3>Expenses</h3>
                        {state.dog.expenses.map(expense => (
                            <div className="form-unit" key={expense.vendor + expense.date}>
                                <div>Vendor: {expense.vendor}</div>
                                <div>Amount: ${expense.amountInCents / 100}</div>
                                <div>Description: {expense.description}</div>
                                <div>Date: {expense.date}</div>
                            </div>
                        ))}
                        <h3>Total Expenses: ${totalExpenses}</h3>
                        <h3>Animal Control pays: ${animalControlPays}</h3>
                        <h3>Adoption Fee: ${adoptionFeeDollars}.{adoptionFeeCents}</h3>
                        <div>
                            <label htmlFor="adoption-date">Adoption Date:</label>
                            <input type="date" onChange={({ target }) => setState({ ...state, dateAdopted: target.value })} />
                        </div>
                        <div className="form-unit">
                            <span>I confirm every looks correct: </span><button type="button" onClick={submitAdoption}>Complete Adoption</button><button type="button" onClick={() => setState({ ...state, latestApplication: {}, adoptionDate: "" })}>Cancel</button><span style={{ color: "red", paddingLeft: "10px" }}>{state.errorMessage}</span>
                        </div>
                    </>
                )}
            </main>
        </>
    )


}

export default AddAdoption;