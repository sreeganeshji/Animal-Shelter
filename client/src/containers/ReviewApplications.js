import React, { useEffect, useState } from "react";
import { Redirect } from "react-router-dom"

import { getData, sendData } from "../services/api";
import { isLoggedIn } from "../services/auth";
import { NavLinks } from "../components";

const ReviewApplications = () => {
    const [state, setState] = useState({
        loginRedirect: false,
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

    const [applications, setApplications] = useState([]);

    useEffect(() => {
        if (state.cookie.hasOwnProperty("isTrustedVolunteer") && state.cookie.isTrustedVolunteer) {
            getData("/application?status=pending")
                .then(({ target }) => {
                    setApplications(JSON.parse(target.response));
                })
                .catch(error => {
                    console.log(error);
                })
        }
    }, [state.cookie]);

    const submitDecision = (applicationId, decision) => {
        const updatedApplicationState = decision ? "true" : "false";
        sendData(decision, `/application/${applicationId}?approve=${updatedApplicationState}`, "PUT")
            .then(({ target }) => {
                if (target.status === 200) {
                    const remainingApplications = applications.filter(application => {
                        return application.id !== applicationId;
                    });
                    setApplications(remainingApplications);
                }

            })
            .catch(error => {
                console.log("Error updating application status: ", error);
            })
    }

    return (
        <>
            {state.loginRedirect && <Redirect to={{ pathname: "/login", state: { from: "/review-applications" } }} />}
            <NavLinks isTrustedVolunteer={state.cookie.isTrustedVolunteer} />
            <main class="container">
                <h1>Review Applications</h1>
                {applications.length ? (
                    <table>
                        <thead>
                            <tr>
                                <th>First Name</th>
                                <th>Last Name</th>
                                <th>Email</th>
                                <th>Street</th>
                                <th>City</th>
                                <th>State</th>
                                <th>Zip</th>
                                <th>Phone Number</th>
                                <th>Co-Applicant First Name</th>
                                <th>Co-Applicant Last Name</th>
                                <th>Decision</th>
                            </tr>
                        </thead>
                        <tbody>
                            {applications.map(application => {
                                return (
                                    <tr key={application.id}>
                                        <td>{application.applicant.firstName}</td>
                                        <td>{application.applicant.lastName}</td>
                                        <td>{application.applicant.email}</td>
                                        <td>{application.applicant.street}</td>
                                        <td>{application.applicant.city}</td>
                                        <td>{application.applicant.state}</td>
                                        <td>{application.applicant.zip}</td>
                                        <td>{application.applicant.phoneNumber}</td>
                                        <td>{application.coApplicantFirstName}</td>
                                        <td>{application.coApplicantLastName}</td>
                                        <td>
                                            <button type="button" onClick={() => {
                                                submitDecision(application.id, true);
                                            }} >Approve</button>
                                            <button type="button" onClick={() => {
                                                submitDecision(application.id, false);
                                            }} >Reject</button>
                                        </td>
                                    </tr>
                                );
                            })}
                        </tbody>

                    </table >
                ) : (<h3>There are not any applications currently pending review.</h3>)
                }
            </main>
        </>
    );
}

export default ReviewApplications;